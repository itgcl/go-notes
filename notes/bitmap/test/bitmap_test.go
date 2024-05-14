package test

import (
	"fmt"
	"log"
	"testing"

	"go-notes/notes/bitmap"

	"github.com/stretchr/testify/assert"
)

func TestByteToBitMap(t *testing.T) {
	s := "\u0000\u0000\u0000\u0000\u0000\u0080@\u0000\u0000@\u0000@\u0010\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000"
	var bytes []byte
	for _, v := range s {
		bytes = append(bytes, byte(v))
	}
	fmt.Println(bytes)
	bm := bitmap.New()
	if err := bm.Scan(bytes); err != nil {
		log.Fatal(err)
	}
	fmt.Println(bm.ToSlice64())
}

func TestByte(t *testing.T) {
	var b []byte
	b = append(b, 0)
	b = append(b, 1)
	b = append(b, 100)
	b = append(b, 200)
	b = append(b, 255)
	b = append(b, 'a')
	b = append(b, []byte("a")...)
	b = append(b, []byte("张")...)
	fmt.Println(b)
}

func TestBitmap_Add(t *testing.T) {
	b := bitmap.New()

	b.Add(1)
	assert.Equal(t, 1, b.Used)
	assert.Equal(t, 8, b.Capacity)

	b.Add(2)
	assert.Equal(t, 2, b.Used)
	assert.Equal(t, 8, b.Capacity)

	b.Add(8)
	assert.Equal(t, 3, b.Used)
	assert.Equal(t, 16, b.Capacity)

	b.Add(15)
	assert.Equal(t, 4, b.Used)
	assert.Equal(t, 16, b.Capacity)

	b.Add(16)
	assert.Equal(t, 5, b.Used)
	assert.Equal(t, 24, b.Capacity)

	b.Add(16)
	assert.Equal(t, 5, b.Used)
	assert.Equal(t, 24, b.Capacity)
}

func TestBitmap_Remove(t *testing.T) {
	b := bitmap.New()

	b.Add(1)
	b.Add(2)
	b.Add(8)
	b.Add(15)
	b.Add(16)

	b.Remove(1)
	assert.Equal(t, 4, b.Used)
	assert.Equal(t, 24, b.Capacity)

	b.Remove(2)
	assert.Equal(t, 3, b.Used)
	assert.Equal(t, 24, b.Capacity)

	b.Remove(8)
	assert.Equal(t, 2, b.Used)
	assert.Equal(t, 24, b.Capacity)

	b.Remove(15)
	assert.Equal(t, 1, b.Used)
	assert.Equal(t, 24, b.Capacity)

	b.Remove(16)
	assert.Equal(t, 0, b.Used)
	assert.Equal(t, 24, b.Capacity)

	b.Remove(16)
	assert.Equal(t, 0, b.Used)
	assert.Equal(t, 24, b.Capacity)
}

func TestBitmap_Contains(t *testing.T) {
	b := bitmap.New()

	b.Add(1)
	b.Add(2)
	b.Add(8)
	b.Add(15)
	b.Add(16)

	assert.True(t, b.Contains(1))
	assert.True(t, b.Contains(2))
	assert.True(t, b.Contains(8))
	assert.True(t, b.Contains(15))
	assert.True(t, b.Contains(16))

	assert.False(t, b.Contains(0))
	assert.False(t, b.Contains(3))
	assert.False(t, b.Contains(4))
	assert.False(t, b.Contains(10))
	assert.False(t, b.Contains(17))
	assert.False(t, b.Contains(99))
}

func TestBitmap_String(t *testing.T) {
	b := bitmap.New()

	b.Add(1)
	b.Add(2)

	assert.Equal(t, "011", b.String())

	b.Add(8)
	assert.Equal(t, "011000001", b.String())

	b.Add(5)
	assert.Equal(t, "011001001", b.String())
}

func TestBitmap_Scan(t *testing.T) {
	b := bitmap.New()
	err := b.Scan([]byte{(1 << 1) | (1 << 2) | (1 << 4), (1 << 3) | (1 << 6)})
	assert.NoError(t, err)
	assert.Equal(t, 5, b.Used)
	assert.Equal(t, 16, b.Capacity)
	assert.Equal(t, "011010000001001", b.String())

	err = b.Scan(nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, b.Used)
	assert.Equal(t, 0, b.Capacity)
	assert.Equal(t, "", b.String())

	b.Add(3)

	err = b.Scan([]byte{})
	assert.NoError(t, err)
	assert.Equal(t, 0, b.Used)
	assert.Equal(t, 0, b.Capacity)
	assert.Equal(t, "", b.String())

	err = b.Scan([]uint64{1, 2, 3})
	assert.Error(t, err)
}

func TestBitmap_Reset(t *testing.T) {
	b := bitmap.New()
	b.Add(1)
	b.Add(2)
	b.Add(8)
	b.Add(15)
	b.Add(16)

	b.Reset()
	assert.Equal(t, 0, b.Used)
	assert.Equal(t, 0, b.Capacity)
	assert.Equal(t, "", b.String())
}

func BenchmarkBitmap_Add(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()
	b.StartTimer()

	b.RunParallel(func(pb *testing.PB) {
		bm := bitmap.New()
		for pb.Next() {
			bm.Add(1)
		}
	})
}

func BenchmarkBitmap_Contains(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()
	b.StartTimer()

	b.RunParallel(func(pb *testing.PB) {
		bm := bitmap.New()
		for i := 0; i < 100; i++ {
			bm.Add(1)
		}
		for pb.Next() {
			bm.Contains(1)
		}
	})
}

func BenchmarkBitmap_Remove(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()
	b.StartTimer()

	b.RunParallel(func(pb *testing.PB) {
		bm := bitmap.New()
		for i := 0; i < 100; i++ {
			bm.Add(1)
		}
		for pb.Next() {
			bm.Remove(1)
		}
	})
}
