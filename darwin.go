// Copyright 2016 Tim O'Brien. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin
// +build darwin

package jnigi

/*
#cgo LDFLAGS:-ldl -framework CoreFoundation

#include <dlfcn.h>
#include <jni.h>
#include <CoreFoundation/CoreFoundation.h>

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

// call back for dummy source used to make sure the CFRunLoop doesn't exit right away
// This callback is called when the source has fired.
void jnigiCFRSourceCallBack (  void *info  ) {
}

void jnigiRunCFRLoop(void) {
    CFRunLoopSourceContext sourceContext;
    //Create a a sourceContext to be used by our source that makes
	//sure the CFRunLoop doesn't exit right away
	sourceContext.version = 0;
	sourceContext.info = NULL;
	sourceContext.retain = NULL;
	sourceContext.release = NULL;
	sourceContext.copyDescription = NULL;
	sourceContext.equal = NULL;
	sourceContext.hash = NULL;
	sourceContext.schedule = NULL;
	sourceContext.cancel = NULL;
	sourceContext.perform = &jnigiCFRSourceCallBack;
        // Create the Source from the sourceContext
	CFRunLoopSourceRef sourceRef = CFRunLoopSourceCreate (NULL, 0, &sourceContext);
	// Use the constant kCFRunLoopCommonModes to add the source to the set of objects
	// monitored by all the common modes
	CFRunLoopAddSource (CFRunLoopGetCurrent(),sourceRef,kCFRunLoopCommonModes);
	// Park this thread in the runloop
	CFRunLoopRun();
}

*/
import "C"

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"unsafe"
)

const (
	JLI_LOAD_ENV   = "LIBJLI_LOAD"
	JLI_LOAD_YES   = "yes"
	JLI_LOAD_FORCE = "force"
)

func jni_GetDefaultJavaVMInitArgs(args unsafe.Pointer) jint {
	return jint(C.dyn_JNI_GetDefaultJavaVMInitArgs((unsafe.Pointer)(args)))
}

func jni_CreateJavaVM(pvm unsafe.Pointer, penv unsafe.Pointer, args unsafe.Pointer) jint {
	return jint(C.dyn_JNI_CreateJavaVM((**C.JavaVM)(pvm), (*unsafe.Pointer)(penv), (unsafe.Pointer)(args)))
}

// LoadJVMLib loads libjvm.dyo as specified in jvmLibPath
func LoadJVMLib(jvmLibPath string) error {
	// On MacOS we need to preload libjli.dylib to workaround JDK-7131356
	// "No Java runtime present, requesting install".
	// If envar LIBJLI_LOAD; = "yes": load but just log error if load fails, =
	// "force": load and exit with error if load fails.
	if jliLoadEnv := os.Getenv(JLI_LOAD_ENV); jliLoadEnv == JLI_LOAD_YES || jliLoadEnv == JLI_LOAD_FORCE {
		libjliPath := filepath.Join(filepath.Dir(jvmLibPath), "..", "libjli.dylib")
		clibjliPath := cString(libjliPath)
		defer func() {
			if clibjliPath != nil {
				free(clibjliPath)
			}
		}()

		// Do not close JLI library handle until JVM closes
		handlelibjli := C.dlopen((*C.char)(clibjliPath), C.RTLD_NOW|C.RTLD_GLOBAL)
		if handlelibjli == nil {
			if jliLoadEnv == JLI_LOAD_YES {
				log.Printf("WARNING could not dynamically load %s", libjliPath)
			} else if jliLoadEnv == JLI_LOAD_FORCE {
				log.Fatalf("ERROR could not dynamically load %s", libjliPath)
			}
		}
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

func RunCFRLoop() {
	C.jnigiRunCFRLoop()
}
