package bitmap

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"sync"
)

// implements a bitmap

// BitmapDynamic represents a bitmap
type BitmapDynamic struct {
	// the bitmap
	bitmap []byte
	// the maximum number of bits in the bitmap
	// this equals to the length of bitmap * 8
	Capacity int
	// the number of bits set to 1 in the bitmap
	Used int
	// the mutex protecting the bitmap
	mutex sync.Mutex
}

// New creates a new bitmap
func New() *BitmapDynamic {
	return &BitmapDynamic{
		bitmap:   make([]byte, 0),
		Capacity: 0,
		Used:     0,
		mutex:    sync.Mutex{},
	}
}

// coordinate of bit, return the index of bitmap and the bit index
func (b *BitmapDynamic) coordinate(bit int) (index, bitIndex int) {
	// calculate the index of bitmap
	index = bit / 8
	// calculate the bit index
	bitIndex = bit % 8
	return index, bitIndex
}

// needExpand returns true if bit is greater than the capacity of the bitmap
func (b *BitmapDynamic) needExpand(bit int) bool {
	// if bit is greater than the capacity of the bitmap, the bitmap will be expanded
	return bit >= b.Capacity
}

// expand the bitmap
func (b *BitmapDynamic) expand(index int) {
	// calculate the new bitmap
	newBitmap := make([]byte, index+1)
	// copy the old bitmap to the new bitmap
	copy(newBitmap, b.bitmap)
	// set the new bitmap
	b.bitmap = newBitmap
	// set the new capacity
	b.Capacity = (index + 1) * 8
}

// Add a bit to the bitmap
func (b *BitmapDynamic) Add(val int64) error {
	if b.Contains(val) {
		return nil
	}
	bit := int(val)
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// calculate the index of bitmap and the bit index
	index, bitIndex := b.coordinate(bit)
	// if the bitmap need to be expanded, expand the bitmap
	if b.needExpand(bit) {
		b.expand(index)
	}
	// set the bit to 1
	b.bitmap[index] |= 1 << bitIndex
	// increase the number of bits set to 1 in the bitmap
	b.Used++
	return nil
}

// Remove a bit from the bitmap
func (b *BitmapDynamic) Remove(val int64) {
	if !b.Contains(val) {
		return
	}
	bit := int(val)
	b.mutex.Lock()
	defer b.mutex.Unlock()
	// calculate the index of bitmap
	index, bitIndex := b.coordinate(bit)
	// set the bit to 0
	b.bitmap[index] &= ^(1 << bitIndex)
	// decrease the number of bits set to 1 in the bitmap
	b.Used--
}

// Contains checks if the bitmap contains the bit
func (b *BitmapDynamic) Contains(val int64) bool {
	bit := int(val)
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// calculate the index of bitmap and the bit index
	index, bitIndex := b.coordinate(bit)
	// if index greater than the length of the bitmap, return false
	if index >= len(b.bitmap) {
		return false
	}
	// if the bit is set to 1, return true
	return b.bitmap[index]&(1<<bitIndex) != 0
}

// String returns the bitmap string representation, is a string of '1' and '0'
// break when numbers of '1' equals to used
func (b *BitmapDynamic) String() string {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// create a string buffer
	var buffer bytes.Buffer
	var count int
	// for each bit in the bitmap
	for i := 0; i < len(b.bitmap); i++ {
		// for each bit in the byte
		for j := 0; j < 8; j++ {
			// if the bit is set to 1, write '1' to the string buffer
			if b.bitmap[i]&(1<<j) != 0 {
				buffer.WriteString("1")
				count++
			} else {
				// if the bit is set to 0, write '0' to the string buffer
				buffer.WriteString("0")
			}
			// if the number of '1' equals to used, break
			if count == b.Used {
				break
			}
		}
	}
	// return the string buffer
	return buffer.String()
}

// Bytes returns the bitmap byte array
func (b *BitmapDynamic) Bytes() []byte {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.bitmap
}

// Reset the bitmap
func (b *BitmapDynamic) Reset() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.bitmap = make([]byte, 0)
	b.Capacity = 0
	b.Used = 0
}

// ToSlice returns the bitmap as a slice of int
func (b *BitmapDynamic) ToSlice() []int {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// create a slice of int
	var slice []int
	// for each bit in the bitmap
	for i := 0; i < len(b.bitmap); i++ {
		// for each bit in the byte
		for j := 0; j < 8; j++ {
			// if the bit is set to 1, append the bit to the slice
			if b.bitmap[i]&(1<<j) != 0 {
				slice = append(slice, i*8+j)
			}
		}
	}
	return slice
}

// ToSlice64 returns the bitmap list representation
func (b *BitmapDynamic) ToSlice64() []int64 {
	r := b.ToSlice()
	var slice []int64
	for i := range r {
		slice = append(slice, int64(r[i]))
	}
	return slice
}

var (
	_ sql.Scanner   = (*BitmapDynamic)(nil)
	_ driver.Valuer = (*BitmapDynamic)(nil)
)

// Scan implements the sql.Scanner interface.
/* #nosec */
func (b *BitmapDynamic) Scan(src interface{}) error {
	if src == nil {
		b.Reset()
		return nil
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	// only accept []byte
	srcBytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("bitmap: cannot convert %T to Bitmap", src)
	}

	// set the bitmap
	b.bitmap = srcBytes
	// set the capacity
	b.Capacity = len(srcBytes) * 8
	// set the used
	b.Used = 0
	// for each byte in the bitmap
	for _, i := range b.bitmap {
		// for each bit in the i
		for j := 0; j < 8; j++ {
			// if the bit is set to 1, increase the used
			if i&(1<<j) != 0 {
				b.Used++
			}
		}
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (b *BitmapDynamic) Value() (driver.Value, error) {
	return b.Bytes(), nil
}
