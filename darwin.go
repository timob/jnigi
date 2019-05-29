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
	"os"
	"path"
)

func jni_GetDefaultJavaVMInitArgs(args unsafe.Pointer) jint {
	return jint(C.dyn_JNI_GetDefaultJavaVMInitArgs((unsafe.Pointer)(args)))
}

func jni_CreateJavaVM(pvm unsafe.Pointer, penv unsafe.Pointer, args unsafe.Pointer) jint {
	return jint(C.dyn_JNI_CreateJavaVM((**C.JavaVM)(pvm), (*unsafe.Pointer)(penv), (unsafe.Pointer)(args)))
}

func init() {
	var (
		cs unsafe.Pointer
	)
	defer free(cs)

	// First, check if JAVA_HOME is set as an environment variable.
	// On Darwin, this usually is set to something like:
	// /Library/Java/JavaVirtualMachines/jdkVERSION.jdk/Contents/Home
	// Where VERSION is the Java version (i.e. 1.8.0).
	// Just use JAVA_HOME so we don't load the wrong JVM
	if key, ok := os.LookupEnv("JAVA_HOME"); ok {
		cs = cString(path.Join(key, "/jre/lib/server/libjvm.dylib"))
	} else {
		panic("JAVA_HOME is not set, set it to the JDK path.")
	}

	libHandle := uintptr(C.dlopen((*C.char)(cs), C.RTLD_NOW|C.RTLD_GLOBAL))
	if libHandle == 0 {
		panic("could not dyanmically load libjvm.dylib")
	}

	cs2 := cString("JNI_GetDefaultJavaVMInitArgs")
	defer free(cs2)
	ptr := C.dlsym(unsafe.Pointer(libHandle), (*C.char)(cs2))
	if ptr == nil {
		panic("could not find JNI_GetDefaultJavaVMInitArgs in libjvm.dylib")
	}
	C.var_JNI_GetDefaultJavaVMInitArgs = C.type_JNI_GetDefaultJavaVMInitArgs(ptr)

	cs3 := cString("JNI_CreateJavaVM")
	defer free(cs3)
	ptr = C.dlsym(unsafe.Pointer(libHandle), (*C.char)(cs3))
	if ptr == nil {
		panic("could not find JNI_CreateJavaVM in libjvm.dylib")
	}
	C.var_JNI_CreateJavaVM = C.type_JNI_CreateJavaVM(ptr)
}
