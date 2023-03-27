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

import "math"

func computeDigitsUInt32(num uint32) uint {
	switch {
	case num < 10:
		return 1
	case num < 100:
		return 2
	case num < 1_000:
		return 3
	case num < 10_000:
		return 4
	case num < 100_000:
		return 5
	case num < 1_000_000:
		return 6
	case num < 10_000_000:
		return 7
	case num < 100_000_000:
		return 8
	case num < 1_000_000_000:
		return 9
	default:
		return 10
	}
}

func computeDigits(num uint64) uint {
	switch {
	case num < 10:
		return 1
	case num < 100:
		return 2
	case num < 1_000:
		return 3
	case num < 10_000:
		return 4
	case num < 100_000:
		return 5
	case num < 1_000_000:
		return 6
	case num < 10_000_000:
		return 7
	case num < 100_000_000:
		return 8
	case num < 1_000_000_000:
		return 9
	case num < 10_000_000_000:
		return 10
	case num < 100_000_000_000:
		return 11
	case num < 1_000_000_000_000:
		return 12
	case num < 10_000_000_000_000:
		return 13
	case num < 100_000_000_000_000:
		return 14
	case num < 1_000_000_000_000_000:
		return 15
	case num < 10_000_000_000_000_000:
		return 16
	case num < 100_000_000_000_000_000:
		return 17
	case num < 1_000_000_000_000_000_000:
		return 18
	case num < 10_000_000_000_000_000_000:
		return 19
	default:
		return 20
	}
}

func computeDigitsUInt32_Log10(num uint32) uint {
	return uint(math.Log10(float64(num)) + 1)
}

func computeDigitsUInt32_Branch(num uint32) uint {
	if num >= 10_000 {
		if num >= 1_000_000 {
			if num >= 100_000_000 {
				if num >= 1_000_000_000 {
					return 10 // 4
				}

				return 9 // 4
			}

			if num >= 10_000_000 {
				return 8 // 4
			}

			return 7 // 4
		}

		if num >= 100_000 {
			return 6 // 3
		}

		return 5 // 3
	}

	if num >= 100 {
		if num >= 1000 {
			return 4 // 3
		}

		return 3 // 3
	}

	if num > 10 {
		return 2 // 3
	}

	return 1 // 3
}

func computeDigitsUInt32_Compute(num uint32) uint {
	var n uint = 1
	for num > 9 {
		num /= 10
		n += 1
	}
	return n
}

func computeDigitsUInt32_Switch(num uint32) uint {
	switch {
	case num > 999_999_999:
		return 10
	case num > 99_999_999:
		return 9
	case num > 9_999_999:
		return 8
	case num > 999_999:
		return 7
	case num > 99_999:
		return 6
	case num > 9_999:
		return 5
	case num > 999:
		return 4
	case num > 99:
		return 3
	case num > 9:
		return 2
	default:
		return 1
	}
}
