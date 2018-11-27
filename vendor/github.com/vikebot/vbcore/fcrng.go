package vbcore

import (
	"math"
)

// FCRNG implements a full-cyclic-random-number-generator
type FCRNG struct {
	Seed      uint64
	Size      uint64
	Increment uint64

	gn uint64
}

func NewFCRNG(size uint64) *FCRNG {
	return &FCRNG{
		Seed:      0,
		Size:      size,
		Increment: 7,
	}
}

func (f *FCRNG) Next32() int32 {
	return int32(f.NextU64() % math.MaxInt32)
}

func (f *FCRNG) NextU32() uint32 {
	return uint32(f.NextU64() % math.MaxUint32)
}

func (f *FCRNG) NextU64() uint64 {
	f.gn = (f.gn + f.Increment) % f.Size
	return f.gn
}

func (f *FCRNG) Next64() int64 {
	return int64(f.NextU64() % math.MaxInt64)
}
