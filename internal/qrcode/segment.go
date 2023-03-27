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

import "strings"

const (
	NumericMode SegmentMode = 1 << iota
	AlphaNumericMode
	ByteMode
)

type SegmentMode uint8

type Segment struct {
	mode SegmentMode
	data []byte
}

func NewSegment(mode SegmentMode, data []byte) Segment {
	return Segment{
		mode,
		data,
	}
}

func NewNumericSegment(num uint64) Segment {
	c := computeDigits(num)
	buf := make([]byte, c)

	for c > 0 {
		buf[c-1] = byte(num % 10)
		num /= 10
		c -= 1
	}

	return Segment{
		mode: NumericMode,
		data: buf,
	}
}

func NewAlphaSegment(s string) (Segment, error) {
	n := len(s)
	buf := make([]byte, n)

	for i, ch := range s {
		code, err := GetAlphaCode(ch)
		if err != nil {
			return Segment{}, err
		}
		buf[i] = code
	}

	return Segment{
		mode: AlphaNumericMode,
		data: buf,
	}, nil
}

func NewUTF8Segment(s string) (Segment, error) {
	str := strings.ToValidUTF8(s, "?")

	return Segment{
		mode: ByteMode,
		data: []byte(str),
	}, nil
}

func (seg Segment) DataBits() int {
	l := len(seg.data)

	switch seg.mode {
	case NumericMode:
		c := (l / 3) * 10
		switch l % 3 {
		case 0:
			return c
		case 1:
			return c + 4
		case 2:
			return c + 7
		}
	case AlphaNumericMode:
		c := (l / 2) * 11
		if l%2 == 1 {
			c += 1
		}
		return c
	case ByteMode:
		return l
	default:
		panic("unsupported data mode in segment")
	}

	return 0
}

func (seg Segment) LengthBits(v int) int {
	switch seg.mode {
	case NumericMode:
		switch {
		case v >= 27 && v <= 40:
			return 14
		case v >= 10 && v <= 26:
			return 12
		case v >= 1 && v <= 10:
			return 10
		default:
			panic("invalid version")
		}
	case AlphaNumericMode:
		switch {
		case v >= 27 && v <= 40:
			return 13
		case v >= 10 && v <= 26:
			return 11
		case v >= 1 && v <= 10:
			return 9
		default:
			panic("invalid version")
		}
	case ByteMode:
		switch {
		case v >= 10 && v <= 40:
			return 16
		case v >= 1 && v <= 10:
			return 8
		default:
			panic("invalid version")
		}
	default:
		panic("unsupported data mode in segment")
	}

}

func (seg Segment) Encode(v int, bits *bitStream) {
	bits.Write4(uint8(seg.mode))

	l := len(seg.data)

	switch seg.mode {
	case NumericMode:
		switch {
		case v >= 27 && v <= 40:
			bits.Write8(uint8(l >> 6))
			bits.Write4(uint8(l >> 2))
			bits.Write2(uint8(l))
		case v >= 10 && v <= 26:
			bits.Write8(uint8(l >> 4))
			bits.Write4(uint8(l))
		case v >= 1 && v <= 10:
			bits.Write8(uint8(l >> 2))
			bits.Write2(uint8(l))
		default:
			panic("invalid version")
		}
		encodeNumeric(seg, bits)
	case AlphaNumericMode:
		switch {
		case v >= 27 && v <= 40:
			bits.Write8(uint8(l >> 6))
			bits.Write4(uint8(l >> 1))
			bits.Write1(uint8(l))
		case v >= 10 && v <= 26:
			bits.Write8(uint8(l >> 3))
			bits.Write2(uint8(l >> 1))
			bits.Write1(uint8(l))
		case v >= 1 && v <= 10:
			bits.Write8(uint8(l >> 1))
			bits.Write1(uint8(l))
		default:
			panic("invalid version")
		}
		encodeAlphaNumeric(seg, bits)
	case ByteMode:
		switch {
		case v >= 10 && v <= 40:
			bits.Write8(uint8(l >> 8))
			bits.Write8(uint8(l))
		case v >= 1 && v <= 10:
			bits.Write8(uint8(l >> 1))
		default:
			panic("invalid version")
		}
		encodeBytes(seg, bits)
	default:
		panic("unsupported data mode in segment")
	}
}
