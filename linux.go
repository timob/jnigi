// Copyright 2016 Tim O'Brien. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux
// +build linux

package jnigi

/*
#cgo LDFLAGS:-ldl

#include <dlfcn.h>
#include <jni.h>

typedef jint (*type_JNI_GetDefaultJavaVMInitArgs)(void*);

type_JNI_GetDefaultJavaVMInitArgs var_JNI_GetDefaultJavaVMInitArgs;

jint dyn_JNI_GetDefaultJavaVMInitArgs(void *args) {
    return var_JNI_GetDefaultJavaVMInitArgs(args);
}

typedef jint (*type_JNI_CreateJavaVM)(JavaVM**, void**, void*);

type_JNI_CreateJavaVM var_JNI_CreateJavaVM;

jint dyn_JNI_CreateJavaVM(JavaVM **pvm, void **penv, void *args) {
    return var_JNI_CreateJavaVM(pvm, penv, args);
}

*/
import "C"

import (
	"errors"
	"unsafe"
)

func jni_GetDefaultJavaVMInitArgs(args unsafe.Pointer) jint {
	return jint(C.dyn_JNI_GetDefaultJavaVMInitArgs((unsafe.Pointer)(args)))
}

func jni_CreateJavaVM(pvm unsafe.Pointer, penv unsafe.Pointer, args unsafe.Pointer) jint {
	return jint(C.dyn_JNI_CreateJavaVM((**C.JavaVM)(pvm), (*unsafe.Pointer)(penv), (unsafe.Pointer)(args)))
}

// LoadJVMLib loads libjvm.so as specified in jvmLibPath
func LoadJVMLib(jvmLibPath string) error {
	cs := cString(jvmLibPath)
	defer free(cs)
	libHandle := uintptr(C.dlopen((*C.char)(cs), C.RTLD_NOW|C.RTLD_GLOBAL))
	if libHandle == 0 {
		return errors.New("could not dynamically load libjvm.so")
	}

	cs2 := cString("JNI_GetDefaultJavaVMInitArgs")
	defer free(cs2)
	ptr := C.dlsym(unsafe.Pointer(libHandle), (*C.char)(cs2))
	if ptr == nil {
		return errors.New("could not find JNI_GetDefaultJavaVMInitArgs in libjvm.so")
	}
	C.var_JNI_GetDefaultJavaVMInitArgs = C.type_JNI_GetDefaultJavaVMInitArgs(ptr)

	cs3 := cString("JNI_CreateJavaVM")
	defer free(cs3)
	ptr = C.dlsym(unsafe.Pointer(libHandle), (*C.char)(cs3))
	if ptr == nil {
		return errors.New("could not find JNI_CreateJavaVM in libjvm.so")
	}
	C.var_JNI_CreateJavaVM = C.type_JNI_CreateJavaVM(ptr)
	return nil
}
