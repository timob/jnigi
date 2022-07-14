package jnigi

/*
#include <jni.h>
#include <stdlib.h>

typedef struct ClassLoaderRef {
	jclass clazz;
	jmethodID findClass;
} ClassLoaderRef;

// Gets a global ref to the class loader (and its findClass method) currently active for 'thiz'.
// This ref can then be cached and used on any thread to find additional dependencies that the
// default class loader cannot.
ClassLoaderRef* getClassLoaderRef(JNIEnv *env, jobject thiz) {
	jclass thizClazz = (*env)->GetObjectClass(env, thiz);
	jclass thizClazzClazz = (*env)->GetObjectClass(env, thizClazz);
	jclass classLoaderClazz = (*env)->FindClass(env, "java/lang/ClassLoader");

	jmethodID getClassLoaderMethod = (*env)->GetMethodID(env, thizClazzClazz, "getClassLoader", "()Ljava/lang/ClassLoader;");
	jclass classLoader = (*env)->CallObjectMethod(env, thizClazz, getClassLoaderMethod);

	ClassLoaderRef *ref = (ClassLoaderRef*)malloc(sizeof(ClassLoaderRef));
	ref->clazz = (*env)->NewGlobalRef(env, classLoader);
	ref->findClass = (*env)->GetMethodID(env, classLoaderClazz, "findClass", "(Ljava/lang/String;)Ljava/lang/Class;");

	return ref;
}

// Find a jclass with the given name using both the default and an optional additional loader.
jclass findClass(JNIEnv *env, const char *className, ClassLoaderRef *addtlLoader) {
	// The default FindClass finds system classes (e.g. java/lang/String, android/app/Application),
	// but not custom classes (e.g. your own app classes, or app dependencies).
	// https://stackoverflow.com/questions/13263340/findclass-from-any-thread-in-android-jni#comment58977872_16302771
	jclass clazz;
	clazz = (*env)->FindClass(env, className);
	if (clazz) {
		return clazz;
	}
	if (!addtlLoader) {
		return NULL;
	}

	// Likely "NoClassDefFoundError" is being thrown by env->FindClass. Clear it so it's possible
	// to attempt again.
	(*env)->ExceptionClear(env);

	// An additional class loader can be used to find custom classes.
	// The loader can be derived from JNI_OnLoad or from 'thiz' provided to any JNI function
	// implementation. A loader derived from these sources will find custom classes, but will not
	// find system classes.
	jstring jClassName = (*env)->NewStringUTF(env, className);
	clazz = (*env)->CallObjectMethod(env, addtlLoader->clazz, addtlLoader->findClass, jClassName);
	(*env)->DeleteLocalRef(env, jClassName);
	return clazz;
}
*/
import "C"
import "unsafe"

// findClass attempts to find the class with the given 'name' within env.
// Optionally, addtlLoader can be specified as a second source.
func findClass(env unsafe.Pointer, name unsafe.Pointer, addtlLoader unsafe.Pointer) jclass {
	return jclass(C.findClass((*C.JNIEnv)(env), (*C.char)(name), (*C.ClassLoaderRef)(addtlLoader)))
}

// getClassLoader returns a reference to the class loader used by 'thiz'
func getClassLoader(env unsafe.Pointer, thiz unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.getClassLoaderRef((*C.JNIEnv)(env), (C.jobject)(thiz)))
}
