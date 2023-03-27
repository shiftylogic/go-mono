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

func TestBitStreamWrite1(t *testing.T) {
	bits := newBits(7)
	test.Require(t, bits.Capacity() == 8, "requested capacity rounded up to multiple of 8")
	test.Require(t, len(bits.data) == 1, "1 byte needed for 7 bits of data")
	test.Require(t, bits.index == 0, "index should start at 0")
	test.Require(t, bits.left == 8, "left should start at 8")

	var i uint = 0
	for ; i < 8; i++ {
		test.Require(t, bits.index == 0, "should always be writting to first data byte")
		test.Require(t, bits.Count() == i, fmt.Sprintf("should have %d bits written", i))
		test.Require(t, bits.left == 8-i, fmt.Sprintf("should have %d bits left of allocated 8 bits", 8-i))
		bits.Write1(uint8(i))
	}

	test.Require(t, bits.Count() == 8, "should have filled 8 bits")
	test.Require(t, bits.data[0] == 0b01010101, "written data is incorrect")
}

func TextBitStreamWrite2(t *testing.T) {
	bits := newBits(10)
	test.Require(t, bits.Capacity() == 16, "requested capacity rounded up to multiple of 8")
	test.Require(t, len(bits.data) == 2, "2 bytes needed for 10 bits of data")
	test.Require(t, bits.index == 0, "index should start at 0")
	test.Require(t, bits.left == 8, "left should start at 8")

	bits.Write2(0b11)
	test.Require(t, bits.Count() == 2, "should have 2 bits written")
	test.Require(t, bits.index == 0, "index should not have changed yet")
	test.Require(t, bits.left == 6, "used only 2 of 8 bits")
	test.Require(t, bits.data[0] == 0b11000000, "data incorrectly written")

	bits.Write2(0b01)
	test.Require(t, bits.Count() == 4, "should have 4 bits written")
	test.Require(t, bits.index == 0, "index should not have changed yet")
	test.Require(t, bits.left == 4, "used only 4 of 8 bits")
	test.Require(t, bits.data[0] == 0b11010000, "data incorrectly written")

	bits.Write2(0b10)
	test.Require(t, bits.Count() == 6, "should have 6 bits written")
	test.Require(t, bits.index == 0, "index should not have changed yet")
	test.Require(t, bits.left == 2, "used only 6 of 8 bits")
	test.Require(t, bits.data[0] == 0b11011000, "data incorrectly written")

	bits.Write2(0b00)
	test.Require(t, bits.Count() == 8, "should have 8 bits written")
	test.Require(t, bits.index == 1, "index should have moved to second byte")
	test.Require(t, bits.left == 8, "moved index should refresh left")
	test.Require(t, bits.data[0] == 0b11011000, "data incorrectly written")

	bits.Write2(0b01)
	test.Require(t, bits.Count() == 10, "should have 10 bits written")
	test.Require(t, bits.index == 1, "index should not have changed yet")
	test.Require(t, bits.left == 6, "moved index should refresh left")
	test.Require(t, bits.data[0] == 0b11011000, "data byte 0 incorrect")
	test.Require(t, bits.data[1] == 0b11000000, "data byte 1 incorrect")
}

func TestBitStreamWrite4(t *testing.T) {
	bits := newBits(20)
	test.Require(t, bits.Capacity() == 24, "requested capacity rounded up to multiple of 8")
	test.Require(t, len(bits.data) == 3, "3 bytes needed for 20 bits of data")
	test.Require(t, bits.index == 0, "index should start at 0")
	test.Require(t, bits.left == 8, "left should start at 8")

	bits.Write4(0b1011)
	test.Require(t, bits.Count() == 4, "should have 4 bits written")
	test.Require(t, bits.index == 0, "index should not have changed yet")
	test.Require(t, bits.left == 4, "used only 4 of 8 bits")
	test.Require(t, bits.data[0] == 0b10110000, "data incorrectly written")

	bits.Write4(0b0110)
	test.Require(t, bits.Count() == 8, "should have 8 bits written")
	test.Require(t, bits.index == 1, "index should have moved to second byte")
	test.Require(t, bits.left == 8, "moved index should refresh left")
	test.Require(t, bits.data[0] == 0b10110110, "data incorrectly written")

	bits.Write4(0b0101)
	test.Require(t, bits.Count() == 12, "should have 12 bits written")
	test.Require(t, bits.index == 1, "index should not have changed yet")
	test.Require(t, bits.left == 4, "used only 4 of 8 bits")
	test.Require(t, bits.data[0] == 0b10110110, "data byte 0 incorrect")
	test.Require(t, bits.data[1] == 0b01010000, "data byte 1 incorrect")

	bits.Write4(0b1100)
	test.Require(t, bits.Count() == 16, "should have 16 bits written")
	test.Require(t, bits.index == 2, "index should have moved to third byte")
	test.Require(t, bits.left == 8, "moved index should refresh left")
	test.Require(t, bits.data[0] == 0b10110110, "data byte 0 incorrect")
	test.Require(t, bits.data[1] == 0b01011100, "data byte 1 incorrect")

	bits.Write4(0b1001)
	test.Require(t, bits.Count() == 20, "should have 20 bits written")
	test.Require(t, bits.index == 2, "index should not have changed yet")
	test.Require(t, bits.left == 4, "used only 4 of 8 bits")
	test.Require(t, bits.data[0] == 0b10110110, "data byte 0 incorrect")
	test.Require(t, bits.data[1] == 0b01011100, "data byte 1 incorrect")
	test.Require(t, bits.data[2] == 0b10010000, "data byte 2 incorrect")
}

func TestBitStreamWrite8(t *testing.T) {
	bits := newBits(16)
	test.Require(t, bits.Capacity() == 16, "requested capacity rounded up to multiple of 8")
	test.Require(t, len(bits.data) == 2, "2 bytes needed for 10 bits of data")
	test.Require(t, bits.index == 0, "index should start at 0")
	test.Require(t, bits.left == 8, "left should start at 8")

	bits.Write8(0b11001100)
	test.Require(t, bits.Count() == 8, "should have 8 bits written")
	test.Require(t, bits.index == 1, "index should have moved to second byte")
	test.Require(t, bits.left == 8, "moved index should refresh left")
	test.Require(t, bits.data[0] == 0b11001100, "data incorrectly written")

	bits.Write8(0b10100101)
	test.Require(t, bits.Count() == 16, "should have 16 bits written")
	test.Require(t, bits.index == 2, "index should have moved to third byte")
	test.Require(t, bits.left == 8, "moved index should refresh left")
	test.Require(t, bits.data[0] == 0b11001100, "data byte 0 incorrect")
	test.Require(t, bits.data[1] == 0b10100101, "data byte 1 incorrect")
}

func TestBitStreamUnaligned2(t *testing.T) {
	bits := newBits(10)
	bits.Write1(0b1)
	bits.Write2(0b11)
	bits.Write2(0b00)
	bits.Write2(0b10)
	bits.Write2(0b01)
	test.Require(t, bits.Count() == 9, "should have 9 bits written")
	test.Require(t, bits.index == 1, "should be on second byte")
	test.Require(t, bits.left == 7, "should have only used 1 bit in second byte")
	test.Require(t, bits.data[0] == 0b11100100, "data byte 0 incorrect")
	test.Require(t, bits.data[1] == 0b10000000, "data byte 1 incorrect")
}

func TestBitStreamUnaligned4(t *testing.T) {
	// Misaligned by 1 bit
	{
		bits := newBits(16)
		bits.Write1(0b1)
		bits.Write4(0b1010)
		bits.Write4(0b1101)
		test.Require(t, bits.Count() == 9, "should have 9 bits written")
		test.Require(t, bits.index == 1, "should be on second byte")
		test.Require(t, bits.left == 7, "should have only used 1 bit in second byte")
		test.Require(t, bits.data[0] == 0b11010110, "data byte 0 incorrect")
		test.Require(t, bits.data[1] == 0b10000000, "data byte 1 incorrect")
	}

	// Misaligned by 2 bits
	{
		bits := newBits(16)
		bits.Write2(0b01)
		bits.Write4(0b1010)
		bits.Write4(0b1101)
		test.Require(t, bits.Count() == 10, "should have 10 bits written")
		test.Require(t, bits.index == 1, "should be on second byte")
		test.Require(t, bits.left == 6, "should have only used 2 bit in second byte")
		test.Require(t, bits.data[0] == 0b01101011, "data byte 0 incorrect")
		test.Require(t, bits.data[1] == 0b01000000, "data byte 1 incorrect")
	}

	// Misaligned by 3 bits
	{
		bits := newBits(16)
		bits.Write1(0b1)
		bits.Write2(0b01)
		bits.Write4(0b1010)
		bits.Write4(0b1101)
		test.Require(t, bits.Count() == 11, "should have 11 bits written")
		test.Require(t, bits.index == 1, "should be on second byte")
		test.Require(t, bits.left == 5, "should have only used 3 bit in second byte")
		test.Require(t, bits.data[0] == 0b10110101, "data byte 0 incorrect")
		test.Require(t, bits.data[1] == 0b10100000, "data byte 1 incorrect")
	}
}

func TestBitStreamUnaligned8(t *testing.T) {
	// Misaligned by 1 bit
	{
		bits := newBits(16)
		bits.Write1(0b1)
		bits.Write8(0b10101011)
		test.Require(t, bits.Count() == 9, "should have 9 bits written")
		test.Require(t, bits.index == 1, "should be on second byte")
		test.Require(t, bits.left == 7, "should have only used 1 bit in second byte")
		test.Require(t, bits.data[0] == 0b11010101, "data byte 0 incorrect")
		test.Require(t, bits.data[1] == 0b10000000, "data byte 1 incorrect")
	}

	// Misaligned by 2 bits
	{
		bits := newBits(16)
		bits.Write2(0b01)
		bits.Write8(0b11001100)
		bits.Write1(0b1)
		test.Require(t, bits.Count() == 11, "should have 11 bits written")
		test.Require(t, bits.index == 1, "should be on second byte")
		test.Require(t, bits.left == 5, "should have only used 3 bit in second byte")
		test.Require(t, bits.data[0] == 0b01110011, "data byte 0 incorrect")
		test.Require(t, bits.data[1] == 0b00100000, "data byte 1 incorrect")
	}

	// Misaligned by 3 bits
	{
		bits := newBits(16)
		bits.Write1(0b1)
		bits.Write2(0b01)
		bits.Write4(0b1010)
		bits.Write4(0b1101)
		test.Require(t, bits.Count() == 11, "should have 11 bits written")
		test.Require(t, bits.index == 1, "should be on second byte")
		test.Require(t, bits.left == 5, "should have only used 3 bit in second byte")
		test.Require(t, bits.data[0] == 0b10110101, "data byte 0 incorrect")
		test.Require(t, bits.data[1] == 0b10100000, "data byte 1 incorrect")
	}

	// Misaligned by 4 bits
	{
		bits := newBits(16)
		bits.Write4(0b1001)
		bits.Write8(0b00110011)
		test.Require(t, bits.Count() == 12, "should have 12 bits written")
		test.Require(t, bits.index == 1, "should be on second byte")
		test.Require(t, bits.left == 4, "should have only used 4 bit in second byte")
		test.Require(t, bits.data[0] == 0b10010011, "data byte 0 incorrect")
		test.Require(t, bits.data[1] == 0b00110000, "data byte 1 incorrect")
	}

	// Misaligned by 5 bits
	{
		bits := newBits(16)
		bits.Write1(0b0)
		bits.Write4(0b1001)
		bits.Write8(0b00111001)
		test.Require(t, bits.Count() == 13, "should have 13 bits written")
		test.Require(t, bits.index == 1, "should be on second byte")
		test.Require(t, bits.left == 3, "should have only used 5 bit in second byte")
		test.Require(t, bits.data[0] == 0b01001001, "data byte 0 incorrect")
		test.Require(t, bits.data[1] == 0b11001000, "data byte 1 incorrect")
	}

	// Misaligned by 6 bits
	{
		bits := newBits(16)
		bits.Write2(0b11)
		bits.Write4(0b1001)
		bits.Write8(0b11000101)
		test.Require(t, bits.Count() == 14, "should have 14 bits written")
		test.Require(t, bits.index == 1, "should be on second byte")
		test.Require(t, bits.left == 2, "should have only used 6 bit in second byte")
		test.Require(t, bits.data[0] == 0b11100111, "data byte 0 incorrect")
		test.Require(t, bits.data[1] == 0b00010100, "data byte 1 incorrect")
	}

	// Misaligned by 7 bits
	{
		bits := newBits(16)
		bits.Write1(0b1)
		bits.Write2(0b01)
		bits.Write4(0b1010)
		bits.Write8(0b11111111)
		test.Require(t, bits.Count() == 15, "should have 15 bits written")
		test.Require(t, bits.index == 1, "should be on second byte")
		test.Require(t, bits.left == 1, "should have only used 7 bit in second byte")
		test.Require(t, bits.data[0] == 0b10110101, "data byte 0 incorrect")
		test.Require(t, bits.data[1] == 0b11111110, "data byte 1 incorrect")
	}
}
