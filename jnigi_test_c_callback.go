//go:build cgo_testing

package jnigi

import (
	"unsafe"
)

/*
#include<stdint.h>

extern uintptr_t go_callback_Greet(void *env, uintptr_t obj, uintptr_t arg_0 );

*/
import "C"

//export go_callback_Greet
func go_callback_Greet(jenv unsafe.Pointer, jobj uintptr, arg_0 uintptr) uintptr {
	env := WrapEnv(jenv)
	defer env.DeleteGlobalRefCache()
	env.ExceptionHandler = ThrowableToStringExceptionHandler

	strArgRef := WrapJObject(arg_0, "java/lang/String", false)
	var goBytes []byte
	if err := strArgRef.CallMethod(env, "getBytes", &goBytes); err != nil {
		panic(err)
	}
	retRef, err := env.NewObject("java/lang/String", []byte("Hello "+string(goBytes)+"!"))
	if err != nil {
		panic(err)
	}
	return uintptr(retRef.JObject())
}

var c_go_callback_Greet = C.go_callback_Greet
