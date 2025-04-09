package main

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type LilEnUInt interface {
	uint16 | uint32 | uint64
}

func ToLittleEndianG[N LilEnUInt](number N) N {
	b := (*[8]byte)(unsafe.Pointer(&number))[:]
	var r N
	j := 0
	for i := int(unsafe.Sizeof(number)) - 1; i >= 0; i-- {
		r = r | N(b[j])<<uint(i*8)
		j++
	}
	return r
}

func TestСonversion(t *testing.T) {
	tests16 := map[string]struct {
		number uint16
		result uint16
	}{
		"test case #6 uint16": {
			number: 0x0102,
			result: 0x0201,
		},
	}

	for name, test := range tests16 {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndianG(test.number)
			assert.Equal(t, test.result, result)
		})
	}

	tests32 := map[string]struct {
		number uint32
		result uint32
	}{
		"test case #1": {
			number: 0x00000000,
			result: 0x00000000,
		},
		"test case #2": {
			number: 0xFFFFFFFF,
			result: 0xFFFFFFFF,
		},
		"test case #3": {
			number: 0x00FF00FF,
			result: 0xFF00FF00,
		},
		"test case #4": {
			number: 0x0000FFFF,
			result: 0xFFFF0000,
		},
		"test case #5": {
			number: 0x01020304,
			result: 0x04030201,
		},
	}

	for name, test := range tests32 {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndianG(test.number)
			assert.Equal(t, test.result, result)
		})
	}

	tests64 := map[string]struct {
		number uint64
		result uint64
	}{
		"test case #7 uint64": {
			number: 0x010203040000FFFF,
			result: 0xFFFF000004030201,
		},
	}

	for name, test := range tests64 {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndianG(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}
