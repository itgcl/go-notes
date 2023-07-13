package bitmap

type Bitmap interface {
	coordinate(bit int) (index, bitIndex int)
	needExpand(bit int) bool
	Add(val int64) error
	Remove(val int64)
	Contains(val int64) bool
	String() string
	Bytes() []byte
	ToSlice() []int
	ToSlice64() []int64
	Reset()
}
