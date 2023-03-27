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
	"math/rand"
	"testing"

	"shiftylogic.dev/site-plat/internal/test"
)

func TestAlphaCodes(t *testing.T) {
	inputs := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:"

	for expected, ch := range inputs {
		t.Run(fmt.Sprintf("success_%c", ch), func(t *testing.T) {
			// Table version
			actual, err := GetAlphaCode(ch)
			test.NoError(t, err, "should be no error")
			test.Expect(t, uint8(expected), actual, "unexpected code")

			// Switch version
			actual, err = GetAlphaCode_Switch(ch)
			test.NoError(t, err, "(switch) should be no error")
			test.Expect(t, uint8(expected), actual, "(switch )unexpected code")
		})
	}
}

func TestAlphaCodeErrors(t *testing.T) {
	seed := rand.Int63()
	t.Logf("Random seed: %d", seed)

	r := rand.New(rand.NewSource(seed))

	for _, ch := range "abcdefghijklmnopqrstuvwxyz" {
		// Table version
		_, err := GetAlphaCode(ch)
		test.AnyError(t, err, fmt.Sprintf("expected error (%c)", ch))

		// Switch version
		_, err = GetAlphaCode_Switch(ch)
		test.AnyError(t, err, fmt.Sprintf("(switch) expected error (%c)", ch))
	}

	for i := 0; i < 1_000_000; i++ {
		v := r.Int31n(test.MaxInt32-90) + 90
		ch := rune(v)

		// Table version
		_, err := GetAlphaCode(ch)
		test.AnyError(t, err, fmt.Sprintf("expected error ('%c' | V: %d)", ch, v))

		// Switch version
		_, err = GetAlphaCode_Switch(ch)
		test.AnyError(t, err, fmt.Sprintf("(switch) expected error ('%c' | V: %d)", ch, v))
	}
}

func BenchmarkAlphaCodes(b *testing.B) {
	inputs := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:"
	method := []struct {
		tag string
		fn  func(rune) (uint8, error)
	}{
		{tag: "noop", fn: func(rune) (uint8, error) { return 0, nil }},
		{tag: "table", fn: GetAlphaCode},
		{tag: "switch", fn: GetAlphaCode_Switch},
	}

	for _, m := range method {
		b.Run(fmt.Sprintf("%s", m.tag), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, ch := range inputs {
					m.fn(ch)
				}
			}
		})
	}
}
