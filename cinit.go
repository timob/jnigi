// Copyright 2016 Tim O'Brien. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jnigi

/*
#include<jni.h>
*/
import "C"

import (
	"unsafe"
)

const jni_commit = C.JNI_COMMIT
const jni_abort = C.JNI_ABORT

// These functions deal with initializing C structs

// NewJVMInitArgs builds JNI JavaVMInitArgs using GetDefaultJavaVMInitArgs and parameters.
func NewJVMInitArgs(ignoreUnrecognizedArgs bool, includeDefaultArgs bool, version int, args []string) *JVMInitArgs {
	jvmargs := (*C.JavaVMInitArgs)(calloc(unsafe.Sizeof(C.JavaVMInitArgs{}), 1))
	jvmargs.version = C.jint(version)
	if includeDefaultArgs {
		if jni_GetDefaultJavaVMInitArgs(unsafe.Pointer(jvmargs)) < 0 {
			panic("JNI_GetDefaultJavaVMInitArgs failed")
		}
	}
	if ignoreUnrecognizedArgs {
		jvmargs.ignoreUnrecognized = (C.jboolean)(jboolean(1))
	}
	for _, arg := range args {
		cStr := cString(arg)
		jvmargs.nOptions++
		jvmargs.options = (*C.JavaVMOption)(realloc(unsafe.Pointer(jvmargs.options), uintptr(jvmargs.nOptions)*unsafe.Sizeof(C.JavaVMOption{})))
		if unsafe.Pointer(jvmargs.options) == nil {
			panic("NewJVMInitArgs realloc call failed")
		}
		option := (*C.JavaVMOption)(unsafe.Pointer(uintptr(unsafe.Pointer(jvmargs.options)) + unsafe.Sizeof(C.JavaVMOption{})*uintptr(jvmargs.nOptions-1)))
		option.optionString = (*C.char)(cStr)
		option.extraInfo = unsafe.Pointer(nil)
	}
	return &JVMInitArgs{unsafe.Pointer(jvmargs)}
}

func registerNative(env unsafe.Pointer, class jclass, mnCstr, sigCstr, fptr unsafe.Pointer) int {
	jniNM := (*C.JNINativeMethod)(calloc(unsafe.Sizeof(C.JNINativeMethod{}), 1))
	jniNM.name = (*C.char)(mnCstr)
	jniNM.signature = (*C.char)(sigCstr)
	jniNM.fnPtr = fptr

	r := registerNatives(env, class, unsafe.Pointer(jniNM), 1)
	// free ?
	return int(r)
}
