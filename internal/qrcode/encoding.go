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

func encodeNumeric(segment Segment, bits *bitStream) {
	n := len(segment.data)
	i := 0

	for i+2 < n {
		v := 100 * uint(segment.data[i])
		v += 10 * uint(segment.data[i+1])
		v += uint(segment.data[i+2])

		// 3 digit sets take 10 bits
		bits.Write8(uint8(v >> 2))
		bits.Write2(uint8(v))

		i += 3
	}

	switch n - i {
	case 1:
		// A single digit takes 4 bits
		bits.Write4(segment.data[i])
	case 2:
		v := 10*uint(segment.data[i]) + uint(segment.data[i+1])
		// 2 digit set takes 7 bits
		bits.Write4(uint8(v >> 3))
		bits.Write2(uint8(v >> 1))
		bits.Write1(uint8(v))
	}
}

func encodeAlphaNumeric(segment Segment, bits *bitStream) {
	n := len(segment.data)
	i := 0

	for i+1 < n {
		v := uint(segment.data[i])*45 + uint(segment.data[i+1])

		// 3 digit sets take 10 bits
		bits.Write8(uint8(v >> 3))
		bits.Write2(uint8(v >> 1))
		bits.Write1(uint8(v))

		i += 2
	}

	if i < n {
		v := uint(segment.data[i])
		bits.Write4(uint8(v >> 2))
		bits.Write2(uint8(v))
	}
}

func encodeBytes(segment Segment, bits *bitStream) {
	for _, d := range segment.data {
		bits.Write8(d)
	}
}
