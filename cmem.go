// Copyright 2016 Tim O'Brien. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jnigi

/*
#include<stdlib.h>
*/
import "C"

import (
	"log"
	"unsafe"
)

func malloc(size uintptr) unsafe.Pointer {
	p := C.malloc(C.size_t(size))
	if p == nil {
		log.Panicf("C malloc failed (size = %d)", size)
	}
	return p
}

func calloc(count uintptr, size uintptr) unsafe.Pointer {
	p := C.calloc(C.size_t(count), C.size_t(size))
	if p == nil {
		log.Panicf("C calloc failed (count = %d, size = %d)", count, size)
	}
	return p
}

func realloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer {
	p := C.realloc(ptr, C.size_t(size))
	if p == nil {
		log.Panicf("C realloc failed (size = %d)", size)
	}
	return p
}

func free(ptr unsafe.Pointer) {
	C.free(ptr)
}

func cString(in string) unsafe.Pointer {
	return unsafe.Pointer(C.CString(in))
}
