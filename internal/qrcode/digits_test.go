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

func TestComputeDigitsUInt32(t *testing.T) {
	inputs := []uint32{
		7, 32, 673, 7893, 47893, 190374, 8942303, 47389843, 989932896, 1893673902,
	}

	for i, v := range inputs {
		test.Expect(t, uint(i+1), computeDigitsUInt32(v), "compute digits uint32 (selected)")
		test.Expect(t, uint(i+1), computeDigitsUInt32_Branch(v), "compute digits uint32 (branch)")
		test.Expect(t, uint(i+1), computeDigitsUInt32_Compute(v), "compute digits uint32 (compute)")
		test.Expect(t, uint(i+1), computeDigitsUInt32_Log10(v), "compute digits uint32 (log10)")
		test.Expect(t, uint(i+1), computeDigitsUInt32_Switch(v), "compute digits uint32 (switch)")
	}
}

func BenchmarkComputeDigits(b *testing.B) {
	inputs := []uint32{
		7, 32, 673, 7893, 47893, 190374, 8942303, 47389843, 989932896, 1893673902,
	}
	method := []struct {
		tag string
		fn  func(uint32) uint
	}{
		{tag: "noop", fn: func(uint32) uint { return 0 }},
		{tag: "selected", fn: computeDigitsUInt32},
		{tag: "branch", fn: computeDigitsUInt32_Branch},
		{tag: "compute", fn: computeDigitsUInt32_Compute},
		{tag: "log10", fn: computeDigitsUInt32_Log10},
		{tag: "switch", fn: computeDigitsUInt32_Switch},
	}

	for _, m := range method {
		for _, v := range inputs {
			b.Run(fmt.Sprintf("%s_%d", m.tag, v), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					m.fn(v)
				}
			})
		}
	}

	for _, v := range inputs {
		b.Run(fmt.Sprintf("64bit_%d", v), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				computeDigits(uint64(v))
			}
		})
	}
}
