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

func TestEncoderNumeric(t *testing.T) {
	cases := []struct {
		num      uint64
		bits     uint
		expected []uint8
	}{
		{num: 0, bits: 4, expected: []uint8{0x00}},
		{num: 9, bits: 4, expected: []uint8{0x90}},
		{num: 11, bits: 7, expected: []uint8{0x16}},
		{num: 81, bits: 7, expected: []uint8{0xA2}},
		{num: 123, bits: 10, expected: []uint8{0x1E, 0xC0}},
		{num: 673, bits: 10, expected: []uint8{0xA8, 0x40}},
		{num: 7893, bits: 14, expected: []uint8{0xC5, 0x4C}},
		{num: 47893, bits: 17, expected: []uint8{0x77, 0xAE, 0x80}},
		{num: 190374, bits: 20, expected: []uint8{0x2F, 0x97, 0x60}},
		{num: 8942303, bits: 24, expected: []uint8{0xDF, 0x8E, 0x63}},
		{num: 47389843, bits: 27, expected: []uint8{0x76, 0x78, 0x25, 0x60}},
		{num: 989932897, bits: 30, expected: []uint8{0xF7, 0x7A, 0x4E, 0x04}},
		{num: 3141592653, bits: 34, expected: []uint8{0x4E, 0x89, 0xF4, 0x24, 0xC0}},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("case_%d", c.num), func(t *testing.T) {
			bits := newBits(64)
			segment := NewNumericSegment(c.num)

			encodeNumeric(segment, &bits)

			test.Expect(t, c.bits, bits.Count(), "bit count")
			test.ExpectBits(t, c.expected, bits.data, c.bits, "bits")
		})
	}
}

func TestEncoderAlphaNumeric(t *testing.T) {
	cases := []struct {
		s        string
		bits     uint
		expected []uint8
	}{
		{s: "", bits: 0, expected: []uint8{}},
		{s: "A", bits: 6, expected: []uint8{0x28}},
		{s: "HELLO WORLD", bits: 61, expected: []uint8{0x61, 0x6F, 0x1A, 0x2E, 0x5B, 0x89, 0xA8, 0x68}},
		{s: "3141592653", bits: 55, expected: []uint8{0x11, 0x02, 0xD4, 0x75, 0x06, 0x01, 0xC8}},
		{s: "PROJECT NAYUKI", bits: 77, expected: []uint8{0x90, 0x11, 0x2D, 0x41, 0x53, 0xD8, 0x2B, 0x86, 0x1C, 0xB0}},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			bits := newBits(128)
			segment, err := NewAlphaSegment(c.s)
			test.NoError(t, err, "unexpected error")

			encodeAlphaNumeric(segment, &bits)

			test.Expect(t, c.bits, bits.Count(), "bit count")
			test.ExpectBits(t, c.expected, bits.data, c.bits, "bits")
		})
	}
}

func TestEncoderUTF8(t *testing.T) {
	cases := []struct {
		s        string
		bits     uint
		expected []uint8
	}{
		{s: "", bits: 0, expected: []uint8{}},
		{s: "a", bits: 8, expected: []uint8{0x61}},
		{s: "HELLO world", bits: 88, expected: []uint8{0x48, 0x45, 0x4C, 0x4C, 0x4F, 0x20, 0x77, 0x6F, 0x72, 0x6C, 0x64}},
		{s: "3141592653", bits: 80, expected: []uint8{0x33, 0x31, 0x34, 0x31, 0x35, 0x39, 0x32, 0x36, 0x35, 0x33}},
		{s: "PROJECT NAYUKI", bits: 112, expected: []uint8{
			0x50, 0x52, 0x4F, 0x4A, 0x45, 0x43, 0x54, 0x20, 0x4E, 0x41, 0x59, 0x55, 0x4B, 0x49, 0x61,
		}},
		{s: "aЉ윇😱", bits: 80, expected: []uint8{0x61, 0xD0, 0x89, 0xEC, 0x9C, 0x87, 0xF0, 0x9F, 0x98, 0xB1}},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			bits := newBits(128)
			segment, err := NewUTF8Segment(c.s)
			test.NoError(t, err, "unexpected error")

			encodeBytes(segment, &bits)

			test.Expect(t, c.bits, bits.Count(), "bit count")
			test.ExpectBits(t, c.expected, bits.data, c.bits, "bits")
		})
	}
}

func BenchmarkEncodeNumeric(b *testing.B) {
	inputs := []uint64{
		9, 11, 123, 7893, 47893, 190374, 8942303, 47389843, 989932896, 3141592653,
		1_000_000_000_000_000, 9_999_999_999_999_999_999, 18_000_111_200_000_999_000,
	}
	method := []struct {
		tag string
		fn  func(Segment, *bitStream)
	}{
		{tag: "noop", fn: func(Segment, *bitStream) { return }},
		{tag: "selected", fn: encodeNumeric},
	}

	buf := make([]uint8, 10)

	for _, m := range method {
		for _, v := range inputs {
			segment := NewNumericSegment(v)
			b.Run(fmt.Sprintf("%s_%d", m.tag, v), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					bits := wrapBits(buf)
					m.fn(segment, &bits)
				}
			})
		}
	}
}
