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

import (
	"fmt"
	"testing"

	"shiftylogic.dev/site-plat/internal/test"
)

func TestNumericSegment(t *testing.T) {
	cases := []struct {
		v        uint64
		bits     int
		expected []uint8
	}{
		{v: 0, bits: 4, expected: []uint8{0}},
		{v: 1, bits: 4, expected: []uint8{1}},
		{v: 9, bits: 4, expected: []uint8{9}},
		{v: 10, bits: 7, expected: []uint8{1, 0}},
		{v: 11, bits: 7, expected: []uint8{1, 1}},
		{v: 23, bits: 7, expected: []uint8{2, 3}},
		{v: 50, bits: 7, expected: []uint8{5, 0}},
		{v: 99, bits: 7, expected: []uint8{9, 9}},
		{v: 100, bits: 10, expected: []uint8{1, 0, 0}},
		{v: 101, bits: 10, expected: []uint8{1, 0, 1}},
		{v: 400, bits: 10, expected: []uint8{4, 0, 0}},
		{v: 999, bits: 10, expected: []uint8{9, 9, 9}},
		{v: 1000, bits: 14, expected: []uint8{1, 0, 0, 0}},
		{v: 1001, bits: 14, expected: []uint8{1, 0, 0, 1}},
		{v: 9999, bits: 14, expected: []uint8{9, 9, 9, 9}},
		{v: 10000, bits: 17, expected: []uint8{1, 0, 0, 0, 0}},
		{v: 10001, bits: 17, expected: []uint8{1, 0, 0, 0, 1}},
		{v: 1_000_000_000, bits: 34, expected: []uint8{1, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{v: 1_000_000_001, bits: 34, expected: []uint8{1, 0, 0, 0, 0, 0, 0, 0, 0, 1}},
		{v: 2_999_000_000, bits: 34, expected: []uint8{2, 9, 9, 9, 0, 0, 0, 0, 0, 0}},
		{v: 3_999_999_999, bits: 34, expected: []uint8{3, 9, 9, 9, 9, 9, 9, 9, 9, 9}},
		{v: 4_000_000_000, bits: 34, expected: []uint8{4, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{v: 4_000_000_001, bits: 34, expected: []uint8{4, 0, 0, 0, 0, 0, 0, 0, 0, 1}},
		{v: 4_294_967_295, bits: 34, expected: []uint8{4, 2, 9, 4, 9, 6, 7, 2, 9, 5}},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("case_%d", c.v), func(t *testing.T) {
			segment := NewNumericSegment(c.v)
			test.Expect(t, NumericMode, segment.mode, "mode")
			test.Expect(t, c.expected, segment.data, "bits")
			test.Expect(t, c.bits, segment.DataBits(), "bit count")
		})
	}
}
