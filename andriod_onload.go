//go:build android
// +build android

package jnigi

/*
#cgo LDFLAGS: -llog -landroid
#include <jni.h>
#include <android/log.h>

void setAndroidJVM(void*, void*);

JNIEXPORT jint JNICALL
JNI_OnLoad(JavaVM *vm, void *reserved) {
    // logging failures to logcat since there's no easy way to return errors within loadLibrary()
    const char *tag = "jnigi.JNI_Onload";
    __android_log_write(ANDROID_LOG_VERBOSE, tag, "Invoked");

    // https://developer.android.com/training/articles/perf-jni
	JNIEnv* env;
	jint res = (*vm)->GetEnv(vm, (void**)(&env), JNI_VERSION_1_6);
    __android_log_write(ANDROID_LOG_VERBOSE, tag, "vm->GetEnv returned");
    if (res != JNI_OK) {
        __android_log_print(ANDROID_LOG_ERROR, tag, "vm->GetEnv failed with code %d", res);
        return JNI_ERR;
    }

	// It's possible for setAndroidJVM to deadlock, so extra logging is useful
    __android_log_write(ANDROID_LOG_VERBOSE, tag, "Invoking setAndroidJVM (Go)");
	setAndroidJVM(vm, env);
    __android_log_write(ANDROID_LOG_VERBOSE, tag, "Done");

	return JNI_VERSION_1_6;
}
*/
import "C"
