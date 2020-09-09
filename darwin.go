// Copyright 2016 Tim O'Brien. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin

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
	"unsafe"
	"errors"
)

func jni_GetDefaultJavaVMInitArgs(args unsafe.Pointer) jint {
	return jint(C.dyn_JNI_GetDefaultJavaVMInitArgs((unsafe.Pointer)(args)))
}

func jni_CreateJavaVM(pvm unsafe.Pointer, penv unsafe.Pointer, args unsafe.Pointer) jint {
	return jint(C.dyn_JNI_CreateJavaVM((**C.JavaVM)(pvm), (*unsafe.Pointer)(penv), (unsafe.Pointer)(args)))
}

func LoadJVMLib(jvmLibPath string) error {
	// On MacOS we need to preload libjli.dylib to workaround JDK-7131356 "No Java runtime present, requesting install"
	libjliPath, err := filepath.Abs(filepath.Join(filepath.Dir(jvmLibPath), "..", "jli", "libjli.dylib"))
	if err != nil {
		return errors.New("error resolving absolute path to 'libjli.dylib'." + err.Error())
	}

	clibjliPath := cString(libjliPath)
	defer func() {
		if clibjliPath != nil {
			free(clibjliPath)
		}
	}()

	// Do not close JLI library handle until JVM closes
	handlelibjli := C.dlopen((*C.char)(clibjliPath), C.RTLD_NOW|C.RTLD_GLOBAL)
	if handlelibjli == nil {
		return errors.New("could not dynamically load 'libjli.dylib'")
	}

	cs := cString(jvmLibPath)
	defer func() {
		if cs != nil {
			free(cs)
		}
	}()

	libHandle := uintptr(C.dlopen((*C.char)(cs), C.RTLD_NOW|C.RTLD_GLOBAL))
	if libHandle == 0 {
		return errors.New("could not dynamically load libjvm.dylib")
	}

	cs2 := cString("JNI_GetDefaultJavaVMInitArgs")
	defer free(cs2)
	ptr := C.dlsym(unsafe.Pointer(libHandle), (*C.char)(cs2))
	if ptr == nil {
		return errors.New("could not find JNI_GetDefaultJavaVMInitArgs in libjvm.dylib")
	}
	C.var_JNI_GetDefaultJavaVMInitArgs = C.type_JNI_GetDefaultJavaVMInitArgs(ptr)

	cs3 := cString("JNI_CreateJavaVM")
	defer free(cs3)
	ptr = C.dlsym(unsafe.Pointer(libHandle), (*C.char)(cs3))
	if ptr == nil {
		return errors.New("could not find JNI_CreateJavaVM in libjvm.dylib")
	}
	C.var_JNI_CreateJavaVM = C.type_JNI_CreateJavaVM(ptr)
	return nil
}
