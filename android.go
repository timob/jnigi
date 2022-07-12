//go:build android
// +build android

package jnigi

import "C"
import (
	"fmt"
	"unsafe"
)

// When running inside an Android Application (i.e. as a shared library), only the App's JVM
// can be used, which is passed to JNI_OnLoad when the library is loaded.
//
// Note: Currently there is *only* support for running as a shared library. Support for standalone
// commandline (non APK) apps could be added using libart.so, however it wouldn't be possible to
// use Android APIs which require an Application context (i.e. all the ones you probably want).

var bootstrapped = make(chan struct{})
var sharedJVM *JVM
var sharedEnv *Env

// Called from JNI_OnLoad in C.
//export setAndroidJVM
func setAndroidJVM(vm unsafe.Pointer, env unsafe.Pointer) {
	sharedJVM = &JVM{vm}
	sharedEnv = &Env{jniEnv: env, classCache: make(map[string]jclass)}
	close(bootstrapped)
}

// AndroidJVM returns references to the Android JVM and the initial environment.
//
// The references are immediately set by JNI when the cgo library is loaded with
// System.loadLibrary() within the JVM. To ensure thread safety, the function blocks until
// the references are set (i.e. loadLibrary() has completed).
//
// Note: Do not call from a package init() function. Calling from init() function results in
// deadlock, since the C code is unable to call back into Go until init() has finished (deadlock).
//
// Must call runtime.LockOSThread() first.
func AndroidJVM() (*JVM, *Env, error) {
	<-bootstrapped
	if sharedJVM == nil || sharedEnv == nil {
		return nil, nil, fmt.Errorf("JNI_OnLoad not invoked")
	}
	return sharedJVM, sharedEnv, nil
}
