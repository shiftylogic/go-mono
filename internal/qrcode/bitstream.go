// MIT License
//
// Copyright (c) 2023 Robert Anderson
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package qrcode

type bitStream struct {
	data  []byte
	index uint
	left  uint
}

func newBits(capacity uint32) bitStream {
	needed := (capacity + 7) >> 3
	buf := make([]uint8, needed)

	return bitStream{
		data:  buf,
		index: 0,
		left:  8,
	}
}

func wrapBits(buf []uint8) bitStream {
	return bitStream{
		data:  buf,
		index: 0,
		left:  8,
	}
}

func (bs *bitStream) ensureState(needed uint) {
	l := uint(len(bs.data))

	if bs.index >= l {
		panic("index is invalid")
	}

	avail := (l-bs.index)*8 + bs.left
	if avail < needed {
		panic("not enough space")
	}
}

func (bs *bitStream) Capacity() uint {
	return uint(len(bs.data) * 8)
}

func (bs *bitStream) Count() uint {
	return bs.index*8 + (8 - bs.left)
}

func (bs *bitStream) Write1(n uint8) {
	bs.ensureState(1)

	// Only the bottom bit matters
	n &= 1

	switch bs.left {
	case 0:
		// This should never happen
		panic("invalid bitstream state")
	case 1:
		bs.data[bs.index] |= n
		bs.index += 1
		bs.left = 8
	default:
		bs.data[bs.index] |= n << (bs.left - 1)
		bs.left -= 1
	}
}

func (bs *bitStream) Write2(n byte) {
	bs.ensureState(2)

	// Only the bottom 2 bits matter
	n &= 3

	switch bs.left {
	case 0:
		// This should never happen
		panic("invalid bitstream state")
	case 1:
		bs.data[bs.index] |= n >> 1
		bs.index += 1
		bs.data[bs.index] |= n << 7
		bs.left = 7
	case 2:
		bs.data[bs.index] |= n
		bs.index += 1
		bs.left = 8
	default:
		bs.data[bs.index] |= n << (bs.left - 2)
		bs.left -= 2
	}
}

func (bs *bitStream) Write4(n byte) {
	bs.ensureState(4)

	// Only the bottom bits matter
	n &= 0xF

	switch {
	case bs.left > 4:
		bs.data[bs.index] |= n << (bs.left - 4)
		bs.left -= 4
	case bs.left == 4:
		bs.data[bs.index] |= n
		bs.index += 1
		bs.left = 8
	case bs.left < 4:
		bs.data[bs.index] |= n >> (4 - bs.left)
		bs.index += 1
		bs.left = 4 + bs.left
		bs.data[bs.index] |= n << bs.left
	}
}

func (bs *bitStream) Write8(n byte) {
	bs.ensureState(8)

	switch bs.left {
	case 0:
		// This should never happen
		panic("invalid bitstream state")
	case 8:
		bs.data[bs.index] |= n
		bs.index += 1
	default:
		bs.data[bs.index] |= n >> (8 - bs.left)
		bs.index += 1
		bs.data[bs.index] |= n << bs.left
	}
}
