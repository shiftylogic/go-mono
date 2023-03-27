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
	"errors"
)

const (
	InvalidAlphaCode = 255
)

var (
	kSymbolTable = []uint8{
		36, // <space>
		255, 255, 255,
		37, 38, // <dollar> <percent>
		255, 255, 255, 255,
		39, 40, // <asterisk> <plus>
		255,
		41, 42, 43, // <minus> <period> <slash>
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, // digits
		44, // <colon>
		255, 255, 255, 255, 255, 255,
		10, 11, 12, 13, 14, 15, 16, 17, 18, 19, // A B C D E F G H I J
		20, 21, 22, 23, 24, 25, 26, 27, 28, 29, // K L M N O P Q R S T
		30, 31, 32, 33, 34, 35, // U V W X Y Z
	}
)

func GetAlphaCode(ch rune) (uint8, error) {
	if ch > 'Z' {
		return InvalidAlphaCode, errors.New("character code too high")
	}

	i := int(ch - ' ') // 0 <= i <= 58
	if i >= len(kSymbolTable) {
		return InvalidAlphaCode, errors.New("table index out of bounds")
	}

	if code := kSymbolTable[i]; code <= 44 {
		return code, nil
	}

	return InvalidAlphaCode, errors.New("invalid character code")
}

func GetAlphaCode_Switch(ch rune) (uint8, error) {
	switch {
	case ch >= 'A' && ch <= 'Z':
		return uint8(ch - 'A' + 10), nil
	case ch >= '0' && ch <= '9':
		return uint8(ch - '0'), nil
	case ch == ' ':
		return 36, nil
	case ch == '$':
		return 37, nil
	case ch == '%':
		return 38, nil
	case ch == '*':
		return 39, nil
	case ch == '+':
		return 40, nil
	case ch == '-':
		return 41, nil
	case ch == '.':
		return 42, nil
	case ch == '/':
		return 43, nil
	case ch == ':':
		return 44, nil
	default:
		return InvalidAlphaCode, errors.New("invalid character")
	}
}
