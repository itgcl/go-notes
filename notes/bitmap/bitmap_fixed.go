package bitmap

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"sync"
)

var ErrorExceedingFixedLength = errors.New("exceeding fixed length")

// implements a bitmap

// BitmapFixed represents a bitmap
type BitmapFixed struct {
	// the bitmap
	bitmap []byte
	// the maximum number of bits in the bitmap
	// this equals to the length of bitmap * 8
	capacity int
	// the number of bits set to 1 in the bitmap
	used int
	// the mutex protecting the bitmap
	mutex           sync.Mutex
	byteFixedLength int
}

// New creates a new bitmap
func NewFixed(byteFixedLength int) *BitmapFixed {
	return &BitmapFixed{
		bitmap:          make([]byte, byteFixedLength),
		capacity:        byteFixedLength * 8,
		used:            0,
		mutex:           sync.Mutex{},
		byteFixedLength: byteFixedLength,
	}
}

// coordinate of bit, return the index of bitmap and the bit index
func (bf *BitmapFixed) coordinate(bit int) (index, bitIndex int) {
	// calculate the index of bitmap
	index = bit / 8
	// calculate the bit index
	bitIndex = bit % 8
	return index, bitIndex
}

// needExpand returns true if bit is greater than the capacity of the bitmap
func (bf *BitmapFixed) needExpand(bit int) bool {
	// if bit is greater than the capacity of the bitmap, the bitmap will be expanded
	return bit >= bf.capacity
}

// Add a bit to the bitmap
func (bf *BitmapFixed) Add(val int64) error {
	if bf.Contains(val) {
		return nil
	}
	bit := int(val)
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	// calculate the index of bitmap and the bit index
	index, bitIndex := bf.coordinate(bit)
	// if the bitmap need to be expanded, expand the bitmap
	if bf.needExpand(bit) {
		return ErrorExceedingFixedLength
	}
	// set the bit to 1
	bf.bitmap[index] |= 1 << bitIndex
	// increase the number of bits set to 1 in the bitmap
	bf.used++
	return nil
}

// Remove a bit from the bitmap
func (bf *BitmapFixed) Remove(val int64) {
	if !bf.Contains(val) {
		return
	}
	bit := int(val)
	bf.mutex.Lock()
	defer bf.mutex.Unlock()
	// calculate the index of bitmap
	index, bitIndex := bf.coordinate(bit)
	// set the bit to 0
	bf.bitmap[index] &= ^(1 << bitIndex)
	// decrease the number of bits set to 1 in the bitmap
	bf.used--
}

// Contains checks if the bitmap contains the bit
func (bf *BitmapFixed) Contains(val int64) bool {
	bit := int(val)
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	// calculate the index of bitmap and the bit index
	index, bitIndex := bf.coordinate(bit)
	// if index greater than the length of the bitmap, return false
	if index >= len(bf.bitmap) {
		return false
	}
	// if the bit is set to 1, return true
	return bf.bitmap[index]&(1<<bitIndex) != 0
}

// String returns the bitmap string representation, is a string of '1' and '0'
// break when numbers of '1' equals to used
func (bf *BitmapFixed) String() string {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()
	// create a string buffer
	var count int
	var buffer bytes.Buffer
	// for each bit in the bitmap
	for i := 0; i < len(bf.bitmap); i++ {
		// for each bit in the byte
		for j := 0; j < 8; j++ {
			// if the bit is set to 1, write '1' to the string buffer
			if bf.bitmap[i]&(1<<j) != 0 {
				buffer.WriteString("1")
				count++
			} else {
				// if the bit is set to 0, write '0' to the string buffer
				buffer.WriteString("0")
			}
			// if the number of '1' equals to used, break
			if count == bf.used {
				break
			}
		}
	}
	// return the string buffer
	return buffer.String()
}

// Bytes returns the bitmap byte array
func (bf *BitmapFixed) Bytes() []byte {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	return bf.bitmap
}

// Reset the bitmap
func (bf *BitmapFixed) Reset() {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	bf.bitmap = make([]byte, 0)
	bf.capacity = 0
	bf.used = 0
}

// ToSlice returns the bitmap as a slice of int
func (bf *BitmapFixed) ToSlice() []int {
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	// create a slice of int
	var slice []int
	// for each bit in the bitmap
	for i := 0; i < len(bf.bitmap); i++ {
		// for each bit in the byte
		for j := 0; j < 8; j++ {
			// if the bit is set to 1, append the bit to the slice
			if bf.bitmap[i]&(1<<j) != 0 {
				slice = append(slice, i*8+j)
			}
		}
	}
	return slice
}

// ToSlice64 returns the bitmap list representation
func (bf *BitmapFixed) ToSlice64() []int64 {
	r := bf.ToSlice()
	var slice []int64
	for i := range r {
		slice = append(slice, int64(r[i]))
	}
	return slice
}

var (
	_ sql.Scanner   = (*BitmapFixed)(nil)
	_ driver.Valuer = (*BitmapFixed)(nil)
)

// Scan implements the sql.Scanner interface. /* #nosec */
func (bf *BitmapFixed) Scan(src interface{}) error {
	if src == nil {
		bf.Reset()
		return nil
	}
	bf.mutex.Lock()
	defer bf.mutex.Unlock()

	// only accept []byte
	srcBytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("bitmap: cannot convert %T to Bitmap", src)
	}
	// set the used
	bf.used = 0
	// set the capacity
	bf.capacity = len(srcBytes) * 8
	// set the bitmap
	bf.bitmap = srcBytes

	// for each byte in the bitmap
	for _, i := range bf.bitmap {
		// for each bit in the i
		for j := 0; j < 8; j++ {
			// if the bit is set to 1, increase the used
			if i&(1<<j) != 0 {
				bf.used++
			}
		}
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (bf *BitmapFixed) Value() (driver.Value, error) {
	return bf.Bytes(), nil
}
