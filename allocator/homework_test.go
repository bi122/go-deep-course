package main

import (
	"reflect"
	"slices"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func Defragment(memory []byte, pointers []unsafe.Pointer) {
	ptrsCnt := len(pointers)
	provenPtrsCnt := 0
	for i := 0; i < len(memory); i++ {
		if provenPtrsCnt == ptrsCnt {
			break
		}
		ptr := unsafe.Pointer(&memory[i])
		if slices.Contains(pointers, ptr) {
			provenPtrsCnt++
		} else {
			memory[i], *(*byte)(pointers[provenPtrsCnt]) = *(*byte)(pointers[provenPtrsCnt]), memory[i]
			pointers[provenPtrsCnt], ptr = ptr, pointers[provenPtrsCnt]
			provenPtrsCnt++
		}
	}
}

func TestDefragmentation(t *testing.T) {
	var fragmentedMemory = []byte{
		0xFF, 0x00, 0x00, 0x00,
		0x00, 0xFF, 0x00, 0x00,
		0x00, 0x00, 0xFF, 0x00,
		0x00, 0x00, 0x00, 0xFF,
	}

	var fragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[5]),
		unsafe.Pointer(&fragmentedMemory[10]),
		unsafe.Pointer(&fragmentedMemory[15]),
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[1]),
		unsafe.Pointer(&fragmentedMemory[2]),
		unsafe.Pointer(&fragmentedMemory[3]),
	}

	var defragmentedMemory = []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	Defragment(fragmentedMemory, fragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}
