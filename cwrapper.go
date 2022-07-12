// Copyright 2016 Tim O'Brien. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jnigi

/*
#include<jni.h>
#include<stdlib.h>

// Android is missing the JNI_VERSION_1_8 constant
#ifndef JNI_VERSION_1_8
#define JNI_VERSION_1_8 0x00010008
#endif

jclass FindClass(JNIEnv* env, char* name) {
	return (*env)->FindClass (env, name);
}

jint Throw(JNIEnv* env, jthrowable obj) {
	return (*env)->Throw (env, obj);
}

jint ThrowNew(JNIEnv* env, jclass clazz, char* msg) {
	return (*env)->ThrowNew (env, clazz, msg);
}

jthrowable ExceptionOccurred(JNIEnv* env) {
	return (*env)->ExceptionOccurred (env);
}

void ExceptionDescribe(JNIEnv* env) {
	(*env)->ExceptionDescribe (env);
}

void ExceptionClear(JNIEnv* env) {
	(*env)->ExceptionClear (env);
}

void FatalError(JNIEnv* env, char* msg) {
	(*env)->FatalError (env, msg);
}

jboolean ExceptionCheck(JNIEnv* env) {
	return (*env)->ExceptionCheck (env);
}

jint PushLocalFrame(JNIEnv* env, jint capacity) {
	return (*env)->PushLocalFrame (env, capacity);
}

jobject PopLocalFrame(JNIEnv* env, jobject result) {
	return (*env)->PopLocalFrame (env, result);
}

jobject NewGlobalRef(JNIEnv* env, jobject lobj) {
	return (*env)->NewGlobalRef (env, lobj);
}

void DeleteGlobalRef(JNIEnv* env, jobject gref) {
	(*env)->DeleteGlobalRef (env, gref);
}

void DeleteLocalRef(JNIEnv* env, jobject obj) {
	(*env)->DeleteLocalRef (env, obj);
}

jboolean IsSameObject(JNIEnv* env, jobject obj1, jobject obj2) {
	return (*env)->IsSameObject (env, obj1, obj2);
}

jobject NewLocalRef(JNIEnv* env, jobject ref) {
	return (*env)->NewLocalRef (env, ref);
}

jint EnsureLocalCapacity(JNIEnv* env, jint capacity) {
	return (*env)->EnsureLocalCapacity (env, capacity);
}

jobject AllocObject(JNIEnv* env, jclass clazz) {
	return (*env)->AllocObject (env, clazz);
}

jobject NewObjectA(JNIEnv* env, jclass clazz, jmethodID methodID, jvalue* args) {
	return (*env)->NewObjectA (env, clazz, methodID, args);
}

jclass GetObjectClass(JNIEnv* env, jobject obj) {
	return (*env)->GetObjectClass (env, obj);
}

jboolean IsInstanceOf(JNIEnv* env, jobject obj, jclass clazz) {
	return (*env)->IsInstanceOf (env, obj, clazz);
}

jmethodID GetMethodID(JNIEnv* env, jclass clazz, char* name, char* sig) {
	return (*env)->GetMethodID (env, clazz, name, sig);
}

jobject CallObjectMethodA(JNIEnv* env, jobject obj, jmethodID methodID, jvalue* args) {
	return (*env)->CallObjectMethodA (env, obj, methodID, args);
}

jboolean CallBooleanMethodA(JNIEnv* env, jobject obj, jmethodID methodID, jvalue* args) {
	return (*env)->CallBooleanMethodA (env, obj, methodID, args);
}

jbyte CallByteMethodA(JNIEnv* env, jobject obj, jmethodID methodID, jvalue* args) {
	return (*env)->CallByteMethodA (env, obj, methodID, args);
}

jchar CallCharMethodA(JNIEnv* env, jobject obj, jmethodID methodID, jvalue* args) {
	return (*env)->CallCharMethodA (env, obj, methodID, args);
}

jshort CallShortMethodA(JNIEnv* env, jobject obj, jmethodID methodID, jvalue* args) {
	return (*env)->CallShortMethodA (env, obj, methodID, args);
}

jint CallIntMethodA(JNIEnv* env, jobject obj, jmethodID methodID, jvalue* args) {
	return (*env)->CallIntMethodA (env, obj, methodID, args);
}

jlong CallLongMethodA(JNIEnv* env, jobject obj, jmethodID methodID, jvalue* args) {
	return (*env)->CallLongMethodA (env, obj, methodID, args);
}

jfloat CallFloatMethodA(JNIEnv* env, jobject obj, jmethodID methodID, jvalue* args) {
	return (*env)->CallFloatMethodA (env, obj, methodID, args);
}

jdouble CallDoubleMethodA(JNIEnv* env, jobject obj, jmethodID methodID, jvalue* args) {
	return (*env)->CallDoubleMethodA (env, obj, methodID, args);
}

void CallVoidMethodA(JNIEnv* env, jobject obj, jmethodID methodID, jvalue* args) {
	(*env)->CallVoidMethodA (env, obj, methodID, args);
}

jfieldID GetFieldID(JNIEnv* env, jclass clazz, char* name, char* sig) {
	return (*env)->GetFieldID (env, clazz, name, sig);
}

jobject GetObjectField(JNIEnv* env, jobject obj, jfieldID fieldID) {
	return (*env)->GetObjectField (env, obj, fieldID);
}

jboolean GetBooleanField(JNIEnv* env, jobject obj, jfieldID fieldID) {
	return (*env)->GetBooleanField (env, obj, fieldID);
}

jbyte GetByteField(JNIEnv* env, jobject obj, jfieldID fieldID) {
	return (*env)->GetByteField (env, obj, fieldID);
}

jchar GetCharField(JNIEnv* env, jobject obj, jfieldID fieldID) {
	return (*env)->GetCharField (env, obj, fieldID);
}

jshort GetShortField(JNIEnv* env, jobject obj, jfieldID fieldID) {
	return (*env)->GetShortField (env, obj, fieldID);
}

jint GetIntField(JNIEnv* env, jobject obj, jfieldID fieldID) {
	return (*env)->GetIntField (env, obj, fieldID);
}

jlong GetLongField(JNIEnv* env, jobject obj, jfieldID fieldID) {
	return (*env)->GetLongField (env, obj, fieldID);
}

jfloat GetFloatField(JNIEnv* env, jobject obj, jfieldID fieldID) {
	return (*env)->GetFloatField (env, obj, fieldID);
}

jdouble GetDoubleField(JNIEnv* env, jobject obj, jfieldID fieldID) {
	return (*env)->GetDoubleField (env, obj, fieldID);
}

void SetObjectField(JNIEnv* env, jobject obj, jfieldID fieldID, jobject val) {
	(*env)->SetObjectField (env, obj, fieldID, val);
}

void SetBooleanField(JNIEnv* env, jobject obj, jfieldID fieldID, jboolean val) {
	(*env)->SetBooleanField (env, obj, fieldID, val);
}

void SetByteField(JNIEnv* env, jobject obj, jfieldID fieldID, jbyte val) {
	(*env)->SetByteField (env, obj, fieldID, val);
}

void SetCharField(JNIEnv* env, jobject obj, jfieldID fieldID, jchar val) {
	(*env)->SetCharField (env, obj, fieldID, val);
}

void SetShortField(JNIEnv* env, jobject obj, jfieldID fieldID, jshort val) {
	(*env)->SetShortField (env, obj, fieldID, val);
}

void SetIntField(JNIEnv* env, jobject obj, jfieldID fieldID, jint val) {
	(*env)->SetIntField (env, obj, fieldID, val);
}

void SetLongField(JNIEnv* env, jobject obj, jfieldID fieldID, jlong val) {
	(*env)->SetLongField (env, obj, fieldID, val);
}

void SetFloatField(JNIEnv* env, jobject obj, jfieldID fieldID, jfloat val) {
	(*env)->SetFloatField (env, obj, fieldID, val);
}

void SetDoubleField(JNIEnv* env, jobject obj, jfieldID fieldID, jdouble val) {
	(*env)->SetDoubleField (env, obj, fieldID, val);
}

jmethodID GetStaticMethodID(JNIEnv* env, jclass clazz, char* name, char* sig) {
	return (*env)->GetStaticMethodID (env, clazz, name, sig);
}

jobject CallStaticObjectMethodA(JNIEnv* env, jclass clazz, jmethodID methodID, jvalue* args) {
	return (*env)->CallStaticObjectMethodA (env, clazz, methodID, args);
}

jboolean CallStaticBooleanMethodA(JNIEnv* env, jclass clazz, jmethodID methodID, jvalue* args) {
	return (*env)->CallStaticBooleanMethodA (env, clazz, methodID, args);
}

jbyte CallStaticByteMethodA(JNIEnv* env, jclass clazz, jmethodID methodID, jvalue* args) {
	return (*env)->CallStaticByteMethodA (env, clazz, methodID, args);
}

jchar CallStaticCharMethodA(JNIEnv* env, jclass clazz, jmethodID methodID, jvalue* args) {
	return (*env)->CallStaticCharMethodA (env, clazz, methodID, args);
}

jshort CallStaticShortMethodA(JNIEnv* env, jclass clazz, jmethodID methodID, jvalue* args) {
	return (*env)->CallStaticShortMethodA (env, clazz, methodID, args);
}

jint CallStaticIntMethodA(JNIEnv* env, jclass clazz, jmethodID methodID, jvalue* args) {
	return (*env)->CallStaticIntMethodA (env, clazz, methodID, args);
}

jlong CallStaticLongMethodA(JNIEnv* env, jclass clazz, jmethodID methodID, jvalue* args) {
	return (*env)->CallStaticLongMethodA (env, clazz, methodID, args);
}

jfloat CallStaticFloatMethodA(JNIEnv* env, jclass clazz, jmethodID methodID, jvalue* args) {
	return (*env)->CallStaticFloatMethodA (env, clazz, methodID, args);
}

jdouble CallStaticDoubleMethodA(JNIEnv* env, jclass clazz, jmethodID methodID, jvalue* args) {
	return (*env)->CallStaticDoubleMethodA (env, clazz, methodID, args);
}

void CallStaticVoidMethodA(JNIEnv* env, jclass cls, jmethodID methodID, jvalue* args) {
	(*env)->CallStaticVoidMethodA (env, cls, methodID, args);
}

jfieldID GetStaticFieldID(JNIEnv* env, jclass clazz, char* name, char* sig) {
	return (*env)->GetStaticFieldID (env, clazz, name, sig);
}

jobject GetStaticObjectField(JNIEnv* env, jclass clazz, jfieldID fieldID) {
	return (*env)->GetStaticObjectField (env, clazz, fieldID);
}

jboolean GetStaticBooleanField(JNIEnv* env, jclass clazz, jfieldID fieldID) {
	return (*env)->GetStaticBooleanField (env, clazz, fieldID);
}

jbyte GetStaticByteField(JNIEnv* env, jclass clazz, jfieldID fieldID) {
	return (*env)->GetStaticByteField (env, clazz, fieldID);
}

jchar GetStaticCharField(JNIEnv* env, jclass clazz, jfieldID fieldID) {
	return (*env)->GetStaticCharField (env, clazz, fieldID);
}

jshort GetStaticShortField(JNIEnv* env, jclass clazz, jfieldID fieldID) {
	return (*env)->GetStaticShortField (env, clazz, fieldID);
}

jint GetStaticIntField(JNIEnv* env, jclass clazz, jfieldID fieldID) {
	return (*env)->GetStaticIntField (env, clazz, fieldID);
}

jlong GetStaticLongField(JNIEnv* env, jclass clazz, jfieldID fieldID) {
	return (*env)->GetStaticLongField (env, clazz, fieldID);
}

jfloat GetStaticFloatField(JNIEnv* env, jclass clazz, jfieldID fieldID) {
	return (*env)->GetStaticFloatField (env, clazz, fieldID);
}

jdouble GetStaticDoubleField(JNIEnv* env, jclass clazz, jfieldID fieldID) {
	return (*env)->GetStaticDoubleField (env, clazz, fieldID);
}

void SetStaticObjectField(JNIEnv* env, jclass clazz, jfieldID fieldID, jobject value) {
	(*env)->SetStaticObjectField (env, clazz, fieldID, value);
}

void SetStaticBooleanField(JNIEnv* env, jclass clazz, jfieldID fieldID, jboolean value) {
	(*env)->SetStaticBooleanField (env, clazz, fieldID, value);
}

void SetStaticByteField(JNIEnv* env, jclass clazz, jfieldID fieldID, jbyte value) {
	(*env)->SetStaticByteField (env, clazz, fieldID, value);
}

void SetStaticCharField(JNIEnv* env, jclass clazz, jfieldID fieldID, jchar value) {
	(*env)->SetStaticCharField (env, clazz, fieldID, value);
}

void SetStaticShortField(JNIEnv* env, jclass clazz, jfieldID fieldID, jshort value) {
	(*env)->SetStaticShortField (env, clazz, fieldID, value);
}

void SetStaticIntField(JNIEnv* env, jclass clazz, jfieldID fieldID, jint value) {
	(*env)->SetStaticIntField (env, clazz, fieldID, value);
}

void SetStaticLongField(JNIEnv* env, jclass clazz, jfieldID fieldID, jlong value) {
	(*env)->SetStaticLongField (env, clazz, fieldID, value);
}

void SetStaticFloatField(JNIEnv* env, jclass clazz, jfieldID fieldID, jfloat value) {
	(*env)->SetStaticFloatField (env, clazz, fieldID, value);
}

void SetStaticDoubleField(JNIEnv* env, jclass clazz, jfieldID fieldID, jdouble value) {
	(*env)->SetStaticDoubleField (env, clazz, fieldID, value);
}

jstring NewString(JNIEnv* env, jchar* unicode, jsize len) {
	return (*env)->NewString (env, unicode, len);
}

jsize GetStringLength(JNIEnv* env, jstring str) {
	return (*env)->GetStringLength (env, str);
}

const jchar* GetStringChars(JNIEnv* env, jstring str, jboolean* isCopy) {
	return (*env)->GetStringChars (env, str, isCopy);
}

void ReleaseStringChars(JNIEnv* env, jstring str, jchar* chars) {
	(*env)->ReleaseStringChars (env, str, chars);
}

jstring NewStringUTF(JNIEnv* env, char* utf) {
	return (*env)->NewStringUTF (env, utf);
}

jsize GetStringUTFLength(JNIEnv* env, jstring str) {
	return (*env)->GetStringUTFLength (env, str);
}

const char* GetStringUTFChars(JNIEnv* env, jstring str, jboolean* isCopy) {
	return (*env)->GetStringUTFChars (env, str, isCopy);
}

void ReleaseStringUTFChars(JNIEnv* env, jstring str, char* chars) {
	(*env)->ReleaseStringUTFChars (env, str, chars);
}

jsize GetArrayLength(JNIEnv* env, jarray array) {
	return (*env)->GetArrayLength (env, array);
}

jobjectArray NewObjectArray(JNIEnv* env, jsize len, jclass clazz, jobject init) {
	return (*env)->NewObjectArray (env, len, clazz, init);
}

jobject GetObjectArrayElement(JNIEnv* env, jobjectArray array, jsize index) {
	return (*env)->GetObjectArrayElement (env, array, index);
}

void SetObjectArrayElement(JNIEnv* env, jobjectArray array, jsize index, jobject val) {
	(*env)->SetObjectArrayElement (env, array, index, val);
}

jbooleanArray NewBooleanArray(JNIEnv* env, jsize len) {
	return (*env)->NewBooleanArray (env, len);
}

jbyteArray NewByteArray(JNIEnv* env, jsize len) {
	return (*env)->NewByteArray (env, len);
}

jcharArray NewCharArray(JNIEnv* env, jsize len) {
	return (*env)->NewCharArray (env, len);
}

jshortArray NewShortArray(JNIEnv* env, jsize len) {
	return (*env)->NewShortArray (env, len);
}

jintArray NewIntArray(JNIEnv* env, jsize len) {
	return (*env)->NewIntArray (env, len);
}

jlongArray NewLongArray(JNIEnv* env, jsize len) {
	return (*env)->NewLongArray (env, len);
}

jfloatArray NewFloatArray(JNIEnv* env, jsize len) {
	return (*env)->NewFloatArray (env, len);
}

jdoubleArray NewDoubleArray(JNIEnv* env, jsize len) {
	return (*env)->NewDoubleArray (env, len);
}

jboolean* GetBooleanArrayElements(JNIEnv* env, jbooleanArray array, jboolean* isCopy) {
	return (*env)->GetBooleanArrayElements (env, array, isCopy);
}

jbyte* GetByteArrayElements(JNIEnv* env, jbyteArray array, jboolean* isCopy) {
	return (*env)->GetByteArrayElements (env, array, isCopy);
}

jchar* GetCharArrayElements(JNIEnv* env, jcharArray array, jboolean* isCopy) {
	return (*env)->GetCharArrayElements (env, array, isCopy);
}

jshort* GetShortArrayElements(JNIEnv* env, jshortArray array, jboolean* isCopy) {
	return (*env)->GetShortArrayElements (env, array, isCopy);
}

jint* GetIntArrayElements(JNIEnv* env, jintArray array, jboolean* isCopy) {
	return (*env)->GetIntArrayElements (env, array, isCopy);
}

jlong* GetLongArrayElements(JNIEnv* env, jlongArray array, jboolean* isCopy) {
	return (*env)->GetLongArrayElements (env, array, isCopy);
}

jfloat* GetFloatArrayElements(JNIEnv* env, jfloatArray array, jboolean* isCopy) {
	return (*env)->GetFloatArrayElements (env, array, isCopy);
}

jdouble* GetDoubleArrayElements(JNIEnv* env, jdoubleArray array, jboolean* isCopy) {
	return (*env)->GetDoubleArrayElements (env, array, isCopy);
}

void ReleaseBooleanArrayElements(JNIEnv* env, jbooleanArray array, jboolean* elems, jint mode) {
	(*env)->ReleaseBooleanArrayElements (env, array, elems, mode);
}

void ReleaseByteArrayElements(JNIEnv* env, jbyteArray array, jbyte* elems, jint mode) {
	(*env)->ReleaseByteArrayElements (env, array, elems, mode);
}

void ReleaseCharArrayElements(JNIEnv* env, jcharArray array, jchar* elems, jint mode) {
	(*env)->ReleaseCharArrayElements (env, array, elems, mode);
}

void ReleaseShortArrayElements(JNIEnv* env, jshortArray array, jshort* elems, jint mode) {
	(*env)->ReleaseShortArrayElements (env, array, elems, mode);
}

void ReleaseIntArrayElements(JNIEnv* env, jintArray array, jint* elems, jint mode) {
	(*env)->ReleaseIntArrayElements (env, array, elems, mode);
}

void ReleaseLongArrayElements(JNIEnv* env, jlongArray array, jlong* elems, jint mode) {
	(*env)->ReleaseLongArrayElements (env, array, elems, mode);
}

void ReleaseFloatArrayElements(JNIEnv* env, jfloatArray array, jfloat* elems, jint mode) {
	(*env)->ReleaseFloatArrayElements (env, array, elems, mode);
}

void ReleaseDoubleArrayElements(JNIEnv* env, jdoubleArray array, jdouble* elems, jint mode) {
	(*env)->ReleaseDoubleArrayElements (env, array, elems, mode);
}

void GetBooleanArrayRegion(JNIEnv* env, jbooleanArray array, jsize start, jsize l, jboolean* buf) {
	(*env)->GetBooleanArrayRegion (env, array, start, l, buf);
}

void GetByteArrayRegion(JNIEnv* env, jbyteArray array, jsize start, jsize len, jbyte* buf) {
	(*env)->GetByteArrayRegion (env, array, start, len, buf);
}

void GetCharArrayRegion(JNIEnv* env, jcharArray array, jsize start, jsize len, jchar* buf) {
	(*env)->GetCharArrayRegion (env, array, start, len, buf);
}

void GetShortArrayRegion(JNIEnv* env, jshortArray array, jsize start, jsize len, jshort* buf) {
	(*env)->GetShortArrayRegion (env, array, start, len, buf);
}

void GetIntArrayRegion(JNIEnv* env, jintArray array, jsize start, jsize len, jint* buf) {
	(*env)->GetIntArrayRegion (env, array, start, len, buf);
}

void GetLongArrayRegion(JNIEnv* env, jlongArray array, jsize start, jsize len, jlong* buf) {
	(*env)->GetLongArrayRegion (env, array, start, len, buf);
}

void GetFloatArrayRegion(JNIEnv* env, jfloatArray array, jsize start, jsize len, jfloat* buf) {
	(*env)->GetFloatArrayRegion (env, array, start, len, buf);
}

void GetDoubleArrayRegion(JNIEnv* env, jdoubleArray array, jsize start, jsize len, jdouble* buf) {
	(*env)->GetDoubleArrayRegion (env, array, start, len, buf);
}

void SetBooleanArrayRegion(JNIEnv* env, jbooleanArray array, jsize start, jsize l, jboolean* buf) {
	(*env)->SetBooleanArrayRegion (env, array, start, l, buf);
}

void SetByteArrayRegion(JNIEnv* env, jbyteArray array, jsize start, jsize len, jbyte* buf) {
	(*env)->SetByteArrayRegion (env, array, start, len, buf);
}

void SetCharArrayRegion(JNIEnv* env, jcharArray array, jsize start, jsize len, jchar* buf) {
	(*env)->SetCharArrayRegion (env, array, start, len, buf);
}

void SetShortArrayRegion(JNIEnv* env, jshortArray array, jsize start, jsize len, jshort* buf) {
	(*env)->SetShortArrayRegion (env, array, start, len, buf);
}

void SetIntArrayRegion(JNIEnv* env, jintArray array, jsize start, jsize len, jint* buf) {
	(*env)->SetIntArrayRegion (env, array, start, len, buf);
}

void SetLongArrayRegion(JNIEnv* env, jlongArray array, jsize start, jsize len, jlong* buf) {
	(*env)->SetLongArrayRegion (env, array, start, len, buf);
}

void SetFloatArrayRegion(JNIEnv* env, jfloatArray array, jsize start, jsize len, jfloat* buf) {
	(*env)->SetFloatArrayRegion (env, array, start, len, buf);
}

void SetDoubleArrayRegion(JNIEnv* env, jdoubleArray array, jsize start, jsize len, jdouble* buf) {
	(*env)->SetDoubleArrayRegion (env, array, start, len, buf);
}

jint RegisterNatives(JNIEnv* env, jclass clazz, JNINativeMethod* methods, jint nMethods) {
	return (*env)->RegisterNatives (env, clazz, methods, nMethods);
}

void GetStringRegion(JNIEnv* env, jstring str, jsize start, jsize len, jchar* buf) {
	(*env)->GetStringRegion (env, str, start, len, buf);
}

void GetStringUTFRegion(JNIEnv* env, jstring str, jsize start, jsize len, char* buf) {
	(*env)->GetStringUTFRegion (env, str, start, len, buf);
}

void* GetPrimitiveArrayCritical(JNIEnv* env, jarray array, jboolean* isCopy) {
	return (*env)->GetPrimitiveArrayCritical (env, array, isCopy);
}

void ReleasePrimitiveArrayCritical(JNIEnv* env, jarray array, void* carray, jint mode) {
	(*env)->ReleasePrimitiveArrayCritical (env, array, carray, mode);
}

const jchar* GetStringCritical(JNIEnv* env, jstring string, jboolean* isCopy) {
	return (*env)->GetStringCritical (env, string, isCopy);
}

void ReleaseStringCritical(JNIEnv* env, jstring string, jchar* cstring) {
	(*env)->ReleaseStringCritical (env, string, cstring);
}

jint AttachCurrentThread(JavaVM* vm, void** penv, void* args) {
	return (*vm)->AttachCurrentThread (vm, penv, args);
}

jint DetachCurrentThread(JavaVM* vm) {
	return (*vm)->DetachCurrentThread (vm);
}

jint DestroyJavaVM(JavaVM* vm) {
 	return (*vm)->DestroyJavaVM (vm);
}

jint GetJavaVM(JNIEnv* env, JavaVM** vm) {
	return (*env)->GetJavaVM (env, vm);
}

jobject CallNonvirtualObjectMethodA(JNIEnv* env, jobject obj, jclass clazz, jmethodID methodID, jvalue* args) {
	return (*env)->CallNonvirtualObjectMethodA (env, obj, clazz, methodID, args);
}

jboolean CallNonvirtualBooleanMethodA(JNIEnv* env, jobject obj, jclass clazz, jmethodID methodID, jvalue* args) {
	return (*env)->CallNonvirtualBooleanMethodA (env, obj, clazz, methodID, args);
}

jbyte CallNonvirtualByteMethodA(JNIEnv* env, jobject obj, jclass clazz, jmethodID methodID, jvalue* args) {
	return (*env)->CallNonvirtualByteMethodA (env, obj, clazz, methodID, args);
}

jchar CallNonvirtualCharMethodA(JNIEnv* env, jobject obj, jclass clazz, jmethodID methodID, jvalue* args) {
	return (*env)->CallNonvirtualCharMethodA (env, obj, clazz, methodID, args);
}

jshort CallNonvirtualShortMethodA(JNIEnv* env, jobject obj, jclass clazz, jmethodID methodID, jvalue* args) {
	return (*env)->CallNonvirtualShortMethodA (env, obj, clazz, methodID, args);
}

jint CallNonvirtualIntMethodA(JNIEnv* env, jobject obj, jclass clazz, jmethodID methodID, jvalue* args) {
	return (*env)->CallNonvirtualIntMethodA (env, obj, clazz, methodID, args);
}

jlong CallNonvirtualLongMethodA(JNIEnv* env, jobject obj, jclass clazz, jmethodID methodID, jvalue* args) {
	return (*env)->CallNonvirtualLongMethodA (env, obj, clazz, methodID, args);
}

jfloat CallNonvirtualFloatMethodA(JNIEnv* env, jobject obj, jclass clazz, jmethodID methodID, jvalue* args) {
	return (*env)->CallNonvirtualFloatMethodA (env, obj, clazz, methodID, args);
}

jdouble CallNonvirtualDoubleMethodA(JNIEnv* env, jobject obj, jclass clazz, jmethodID methodID, jvalue* args) {
	return (*env)->CallNonvirtualDoubleMethodA (env, obj, clazz, methodID, args);
}

void CallNonvirtualVoidMethodA(JNIEnv* env, jobject obj, jclass clazz, jmethodID methodID, jvalue* args) {
	(*env)->CallNonvirtualVoidMethodA (env, obj, clazz, methodID, args);
}

*/
import "C"

import "unsafe"

const (
	JNI_VERSION_1_2 = C.JNI_VERSION_1_2
	JNI_VERSION_1_4 = C.JNI_VERSION_1_4
	JNI_VERSION_1_6 = C.JNI_VERSION_1_6
	JNI_VERSION_1_8 = C.JNI_VERSION_1_8

	DEFAULT_VERSION = JNI_VERSION_1_6
)

type (
	jfieldID      uintptr
	jbooleanArray uintptr
	jshortArray   uintptr
	jmethodID     uintptr
	jboolean      C.jboolean
	jdouble       C.jdouble
	jarray        uintptr
	jbyteArray    uintptr
	jclass        uintptr
	jthrowable    uintptr
	jobjectArray  uintptr
	jcharArray    uintptr
	jlongArray    uintptr
	jvalue        C.jvalue
	jshort        C.jshort
	jintArray     uintptr
	jlong         C.jlong
	jfloat        C.jfloat
	jstring       uintptr
	jchar         C.jchar
	jobject       uintptr
	jbyte         C.jbyte
	jsize         C.jsize
	jfloatArray   uintptr
	jdoubleArray  uintptr
	jint          C.jint
)

func findClass(env unsafe.Pointer, name unsafe.Pointer) jclass {
	return jclass(unsafe.Pointer(C.FindClass((*C.JNIEnv)(env), (*C.char)(name))))
}

func throw(env unsafe.Pointer, obj jthrowable) jint {
	return jint(C.Throw((*C.JNIEnv)(env), C.jthrowable(unsafe.Pointer(obj))))
}

func throwNew(env unsafe.Pointer, clazz jclass, msg unsafe.Pointer) jint {
	return jint(C.ThrowNew((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), (*C.char)(msg)))
}

func exceptionOccurred(env unsafe.Pointer) jthrowable {
	return jthrowable(unsafe.Pointer(C.ExceptionOccurred((*C.JNIEnv)(env))))
}

func exceptionDescribe(env unsafe.Pointer) {
	C.ExceptionDescribe((*C.JNIEnv)(env))
}

func exceptionClear(env unsafe.Pointer) {
	C.ExceptionClear((*C.JNIEnv)(env))
}

func fatalError(env unsafe.Pointer, msg unsafe.Pointer) {
	C.FatalError((*C.JNIEnv)(env), (*C.char)(msg))
}

func exceptionCheck(env unsafe.Pointer) jboolean {
	return jboolean(C.ExceptionCheck((*C.JNIEnv)(env)))
}

func pushLocalFrame(env unsafe.Pointer, capacity jint) jint {
	return jint(C.PushLocalFrame((*C.JNIEnv)(env), C.jint(capacity)))
}

func popLocalFrame(env unsafe.Pointer, result jobject) jobject {
	return jobject(unsafe.Pointer(C.PopLocalFrame((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(result)))))
}

func newGlobalRef(env unsafe.Pointer, lobj jobject) jobject {
	return jobject(unsafe.Pointer(C.NewGlobalRef((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(lobj)))))
}

func deleteGlobalRef(env unsafe.Pointer, gref jobject) {
	C.DeleteGlobalRef((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(gref)))
}

func deleteLocalRef(env unsafe.Pointer, obj jobject) {
	C.DeleteLocalRef((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)))
}

func isSameObject(env unsafe.Pointer, obj1 jobject, obj2 jobject) jboolean {
	return jboolean(C.IsSameObject((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj1)), C.jobject(unsafe.Pointer(obj2))))
}

func newLocalRef(env unsafe.Pointer, ref jobject) jobject {
	return jobject(unsafe.Pointer(C.NewLocalRef((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(ref)))))
}

func ensureLocalCapacity(env unsafe.Pointer, capacity jint) jint {
	return jint(C.EnsureLocalCapacity((*C.JNIEnv)(env), C.jint(capacity)))
}

func allocObject(env unsafe.Pointer, clazz jclass) jobject {
	return jobject(unsafe.Pointer(C.AllocObject((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)))))
}

func newObjectA(env unsafe.Pointer, clazz jclass, methodID jmethodID, args unsafe.Pointer) jobject {
	return jobject(unsafe.Pointer(C.NewObjectA((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args))))
}

func getObjectClass(env unsafe.Pointer, obj jobject) jclass {
	return jclass(unsafe.Pointer(C.GetObjectClass((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)))))
}

func isInstanceOf(env unsafe.Pointer, obj jobject, clazz jclass) jboolean {
	return jboolean(C.IsInstanceOf((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jclass(unsafe.Pointer(clazz))))
}

func getMethodID(env unsafe.Pointer, clazz jclass, name unsafe.Pointer, sig unsafe.Pointer) jmethodID {
	return jmethodID(unsafe.Pointer(C.GetMethodID((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), (*C.char)(name), (*C.char)(sig))))
}

func callObjectMethodA(env unsafe.Pointer, obj jobject, methodID jmethodID, args unsafe.Pointer) jobject {
	return jobject(unsafe.Pointer(C.CallObjectMethodA((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args))))
}

func callBooleanMethodA(env unsafe.Pointer, obj jobject, methodID jmethodID, args unsafe.Pointer) jboolean {
	return jboolean(C.CallBooleanMethodA((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callByteMethodA(env unsafe.Pointer, obj jobject, methodID jmethodID, args unsafe.Pointer) jbyte {
	return jbyte(C.CallByteMethodA((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callCharMethodA(env unsafe.Pointer, obj jobject, methodID jmethodID, args unsafe.Pointer) jchar {
	return jchar(C.CallCharMethodA((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callShortMethodA(env unsafe.Pointer, obj jobject, methodID jmethodID, args unsafe.Pointer) jshort {
	return jshort(C.CallShortMethodA((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callIntMethodA(env unsafe.Pointer, obj jobject, methodID jmethodID, args unsafe.Pointer) jint {
	return jint(C.CallIntMethodA((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callLongMethodA(env unsafe.Pointer, obj jobject, methodID jmethodID, args unsafe.Pointer) jlong {
	return jlong(C.CallLongMethodA((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callFloatMethodA(env unsafe.Pointer, obj jobject, methodID jmethodID, args unsafe.Pointer) jfloat {
	return jfloat(C.CallFloatMethodA((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callDoubleMethodA(env unsafe.Pointer, obj jobject, methodID jmethodID, args unsafe.Pointer) jdouble {
	return jdouble(C.CallDoubleMethodA((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callVoidMethodA(env unsafe.Pointer, obj jobject, methodID jmethodID, args unsafe.Pointer) {
	C.CallVoidMethodA((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args))
}

func getFieldID(env unsafe.Pointer, clazz jclass, name unsafe.Pointer, sig unsafe.Pointer) jfieldID {
	return jfieldID(unsafe.Pointer(C.GetFieldID((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), (*C.char)(name), (*C.char)(sig))))
}

func getObjectField(env unsafe.Pointer, obj jobject, fieldID jfieldID) jobject {
	return jobject(unsafe.Pointer(C.GetObjectField((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jfieldID(unsafe.Pointer(fieldID)))))
}

func getBooleanField(env unsafe.Pointer, obj jobject, fieldID jfieldID) jboolean {
	return jboolean(C.GetBooleanField((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jfieldID(unsafe.Pointer(fieldID))))
}

func getByteField(env unsafe.Pointer, obj jobject, fieldID jfieldID) jbyte {
	return jbyte(C.GetByteField((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jfieldID(unsafe.Pointer(fieldID))))
}

func getCharField(env unsafe.Pointer, obj jobject, fieldID jfieldID) jchar {
	return jchar(C.GetCharField((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jfieldID(unsafe.Pointer(fieldID))))
}

func getShortField(env unsafe.Pointer, obj jobject, fieldID jfieldID) jshort {
	return jshort(C.GetShortField((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jfieldID(unsafe.Pointer(fieldID))))
}

func getIntField(env unsafe.Pointer, obj jobject, fieldID jfieldID) jint {
	return jint(C.GetIntField((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jfieldID(unsafe.Pointer(fieldID))))
}

func getLongField(env unsafe.Pointer, obj jobject, fieldID jfieldID) jlong {
	return jlong(C.GetLongField((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jfieldID(unsafe.Pointer(fieldID))))
}

func getFloatField(env unsafe.Pointer, obj jobject, fieldID jfieldID) jfloat {
	return jfloat(C.GetFloatField((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jfieldID(unsafe.Pointer(fieldID))))
}

func getDoubleField(env unsafe.Pointer, obj jobject, fieldID jfieldID) jdouble {
	return jdouble(C.GetDoubleField((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jfieldID(unsafe.Pointer(fieldID))))
}

func setObjectField(env unsafe.Pointer, obj jobject, fieldID jfieldID, val jobject) {
	C.SetObjectField((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jfieldID(unsafe.Pointer(fieldID)), C.jobject(unsafe.Pointer(val)))
}

func setBooleanField(env unsafe.Pointer, obj jobject, fieldID jfieldID, val jboolean) {
	C.SetBooleanField((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jfieldID(unsafe.Pointer(fieldID)), C.jboolean(val))
}

func setByteField(env unsafe.Pointer, obj jobject, fieldID jfieldID, val jbyte) {
	C.SetByteField((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jfieldID(unsafe.Pointer(fieldID)), C.jbyte(val))
}

func setCharField(env unsafe.Pointer, obj jobject, fieldID jfieldID, val jchar) {
	C.SetCharField((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jfieldID(unsafe.Pointer(fieldID)), C.jchar(val))
}

func setShortField(env unsafe.Pointer, obj jobject, fieldID jfieldID, val jshort) {
	C.SetShortField((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jfieldID(unsafe.Pointer(fieldID)), C.jshort(val))
}

func setIntField(env unsafe.Pointer, obj jobject, fieldID jfieldID, val jint) {
	C.SetIntField((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jfieldID(unsafe.Pointer(fieldID)), C.jint(val))
}

func setLongField(env unsafe.Pointer, obj jobject, fieldID jfieldID, val jlong) {
	C.SetLongField((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jfieldID(unsafe.Pointer(fieldID)), C.jlong(val))
}

func setFloatField(env unsafe.Pointer, obj jobject, fieldID jfieldID, val jfloat) {
	C.SetFloatField((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jfieldID(unsafe.Pointer(fieldID)), C.jfloat(val))
}

func setDoubleField(env unsafe.Pointer, obj jobject, fieldID jfieldID, val jdouble) {
	C.SetDoubleField((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jfieldID(unsafe.Pointer(fieldID)), C.jdouble(val))
}

func getStaticMethodID(env unsafe.Pointer, clazz jclass, name unsafe.Pointer, sig unsafe.Pointer) jmethodID {
	return jmethodID(unsafe.Pointer(C.GetStaticMethodID((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), (*C.char)(name), (*C.char)(sig))))
}

func callStaticObjectMethodA(env unsafe.Pointer, clazz jclass, methodID jmethodID, args unsafe.Pointer) jobject {
	return jobject(unsafe.Pointer(C.CallStaticObjectMethodA((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args))))
}

func callStaticBooleanMethodA(env unsafe.Pointer, clazz jclass, methodID jmethodID, args unsafe.Pointer) jboolean {
	return jboolean(C.CallStaticBooleanMethodA((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callStaticByteMethodA(env unsafe.Pointer, clazz jclass, methodID jmethodID, args unsafe.Pointer) jbyte {
	return jbyte(C.CallStaticByteMethodA((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callStaticCharMethodA(env unsafe.Pointer, clazz jclass, methodID jmethodID, args unsafe.Pointer) jchar {
	return jchar(C.CallStaticCharMethodA((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callStaticShortMethodA(env unsafe.Pointer, clazz jclass, methodID jmethodID, args unsafe.Pointer) jshort {
	return jshort(C.CallStaticShortMethodA((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callStaticIntMethodA(env unsafe.Pointer, clazz jclass, methodID jmethodID, args unsafe.Pointer) jint {
	return jint(C.CallStaticIntMethodA((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callStaticLongMethodA(env unsafe.Pointer, clazz jclass, methodID jmethodID, args unsafe.Pointer) jlong {
	return jlong(C.CallStaticLongMethodA((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callStaticFloatMethodA(env unsafe.Pointer, clazz jclass, methodID jmethodID, args unsafe.Pointer) jfloat {
	return jfloat(C.CallStaticFloatMethodA((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callStaticDoubleMethodA(env unsafe.Pointer, clazz jclass, methodID jmethodID, args unsafe.Pointer) jdouble {
	return jdouble(C.CallStaticDoubleMethodA((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callStaticVoidMethodA(env unsafe.Pointer, cls jclass, methodID jmethodID, args unsafe.Pointer) {
	C.CallStaticVoidMethodA((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(cls)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args))
}

func getStaticFieldID(env unsafe.Pointer, clazz jclass, name unsafe.Pointer, sig unsafe.Pointer) jfieldID {
	return jfieldID(unsafe.Pointer(C.GetStaticFieldID((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), (*C.char)(name), (*C.char)(sig))))
}

func getStaticObjectField(env unsafe.Pointer, clazz jclass, fieldID jfieldID) jobject {
	return jobject(unsafe.Pointer(C.GetStaticObjectField((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jfieldID(unsafe.Pointer(fieldID)))))
}

func getStaticBooleanField(env unsafe.Pointer, clazz jclass, fieldID jfieldID) jboolean {
	return jboolean(C.GetStaticBooleanField((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jfieldID(unsafe.Pointer(fieldID))))
}

func getStaticByteField(env unsafe.Pointer, clazz jclass, fieldID jfieldID) jbyte {
	return jbyte(C.GetStaticByteField((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jfieldID(unsafe.Pointer(fieldID))))
}

func getStaticCharField(env unsafe.Pointer, clazz jclass, fieldID jfieldID) jchar {
	return jchar(C.GetStaticCharField((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jfieldID(unsafe.Pointer(fieldID))))
}

func getStaticShortField(env unsafe.Pointer, clazz jclass, fieldID jfieldID) jshort {
	return jshort(C.GetStaticShortField((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jfieldID(unsafe.Pointer(fieldID))))
}

func getStaticIntField(env unsafe.Pointer, clazz jclass, fieldID jfieldID) jint {
	return jint(C.GetStaticIntField((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jfieldID(unsafe.Pointer(fieldID))))
}

func getStaticLongField(env unsafe.Pointer, clazz jclass, fieldID jfieldID) jlong {
	return jlong(C.GetStaticLongField((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jfieldID(unsafe.Pointer(fieldID))))
}

func getStaticFloatField(env unsafe.Pointer, clazz jclass, fieldID jfieldID) jfloat {
	return jfloat(C.GetStaticFloatField((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jfieldID(unsafe.Pointer(fieldID))))
}

func getStaticDoubleField(env unsafe.Pointer, clazz jclass, fieldID jfieldID) jdouble {
	return jdouble(C.GetStaticDoubleField((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jfieldID(unsafe.Pointer(fieldID))))
}

func setStaticObjectField(env unsafe.Pointer, clazz jclass, fieldID jfieldID, value jobject) {
	C.SetStaticObjectField((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jfieldID(unsafe.Pointer(fieldID)), C.jobject(unsafe.Pointer(value)))
}

func setStaticBooleanField(env unsafe.Pointer, clazz jclass, fieldID jfieldID, value jboolean) {
	C.SetStaticBooleanField((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jfieldID(unsafe.Pointer(fieldID)), C.jboolean(value))
}

func setStaticByteField(env unsafe.Pointer, clazz jclass, fieldID jfieldID, value jbyte) {
	C.SetStaticByteField((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jfieldID(unsafe.Pointer(fieldID)), C.jbyte(value))
}

func setStaticCharField(env unsafe.Pointer, clazz jclass, fieldID jfieldID, value jchar) {
	C.SetStaticCharField((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jfieldID(unsafe.Pointer(fieldID)), C.jchar(value))
}

func setStaticShortField(env unsafe.Pointer, clazz jclass, fieldID jfieldID, value jshort) {
	C.SetStaticShortField((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jfieldID(unsafe.Pointer(fieldID)), C.jshort(value))
}

func setStaticIntField(env unsafe.Pointer, clazz jclass, fieldID jfieldID, value jint) {
	C.SetStaticIntField((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jfieldID(unsafe.Pointer(fieldID)), C.jint(value))
}

func setStaticLongField(env unsafe.Pointer, clazz jclass, fieldID jfieldID, value jlong) {
	C.SetStaticLongField((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jfieldID(unsafe.Pointer(fieldID)), C.jlong(value))
}

func setStaticFloatField(env unsafe.Pointer, clazz jclass, fieldID jfieldID, value jfloat) {
	C.SetStaticFloatField((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jfieldID(unsafe.Pointer(fieldID)), C.jfloat(value))
}

func setStaticDoubleField(env unsafe.Pointer, clazz jclass, fieldID jfieldID, value jdouble) {
	C.SetStaticDoubleField((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), C.jfieldID(unsafe.Pointer(fieldID)), C.jdouble(value))
}

func newString(env unsafe.Pointer, unicode unsafe.Pointer, len jsize) jstring {
	return jstring(unsafe.Pointer(C.NewString((*C.JNIEnv)(env), (*C.jchar)(unicode), C.jsize(len))))
}

func getStringLength(env unsafe.Pointer, str jstring) jsize {
	return jsize(C.GetStringLength((*C.JNIEnv)(env), C.jstring(unsafe.Pointer(str))))
}

func getStringChars(env unsafe.Pointer, str jstring, isCopy unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.GetStringChars((*C.JNIEnv)(env), C.jstring(unsafe.Pointer(str)), (*C.jboolean)(isCopy)))
}

func releaseStringChars(env unsafe.Pointer, str jstring, chars unsafe.Pointer) {
	C.ReleaseStringChars((*C.JNIEnv)(env), C.jstring(unsafe.Pointer(str)), (*C.jchar)(chars))
}

func newStringUTF(env unsafe.Pointer, utf unsafe.Pointer) jstring {
	return jstring(unsafe.Pointer(C.NewStringUTF((*C.JNIEnv)(env), (*C.char)(utf))))
}

func getStringUTFLength(env unsafe.Pointer, str jstring) jsize {
	return jsize(C.GetStringUTFLength((*C.JNIEnv)(env), C.jstring(unsafe.Pointer(str))))
}

func getStringUTFChars(env unsafe.Pointer, str jstring, isCopy unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.GetStringUTFChars((*C.JNIEnv)(env), C.jstring(unsafe.Pointer(str)), (*C.jboolean)(isCopy)))
}

func releaseStringUTFChars(env unsafe.Pointer, str jstring, chars unsafe.Pointer) {
	C.ReleaseStringUTFChars((*C.JNIEnv)(env), C.jstring(unsafe.Pointer(str)), (*C.char)(chars))
}

func getArrayLength(env unsafe.Pointer, array jarray) jsize {
	return jsize(C.GetArrayLength((*C.JNIEnv)(env), C.jarray(unsafe.Pointer(array))))
}

func newObjectArray(env unsafe.Pointer, len jsize, clazz jclass, init jobject) jobjectArray {
	return jobjectArray(unsafe.Pointer(C.NewObjectArray((*C.JNIEnv)(env), C.jsize(len), C.jclass(unsafe.Pointer(clazz)), C.jobject(unsafe.Pointer(init)))))
}

func getObjectArrayElement(env unsafe.Pointer, array jobjectArray, index jsize) jobject {
	return jobject(unsafe.Pointer(C.GetObjectArrayElement((*C.JNIEnv)(env), C.jobjectArray(unsafe.Pointer(array)), C.jsize(index))))
}

func setObjectArrayElement(env unsafe.Pointer, array jobjectArray, index jsize, val jobject) {
	C.SetObjectArrayElement((*C.JNIEnv)(env), C.jobjectArray(unsafe.Pointer(array)), C.jsize(index), C.jobject(unsafe.Pointer(val)))
}

func newBooleanArray(env unsafe.Pointer, len jsize) jbooleanArray {
	return jbooleanArray(unsafe.Pointer(C.NewBooleanArray((*C.JNIEnv)(env), C.jsize(len))))
}

func newByteArray(env unsafe.Pointer, len jsize) jbyteArray {
	return jbyteArray(unsafe.Pointer(C.NewByteArray((*C.JNIEnv)(env), C.jsize(len))))
}

func newCharArray(env unsafe.Pointer, len jsize) jcharArray {
	return jcharArray(unsafe.Pointer(C.NewCharArray((*C.JNIEnv)(env), C.jsize(len))))
}

func newShortArray(env unsafe.Pointer, len jsize) jshortArray {
	return jshortArray(unsafe.Pointer(C.NewShortArray((*C.JNIEnv)(env), C.jsize(len))))
}

func newIntArray(env unsafe.Pointer, len jsize) jintArray {
	return jintArray(unsafe.Pointer(C.NewIntArray((*C.JNIEnv)(env), C.jsize(len))))
}

func newLongArray(env unsafe.Pointer, len jsize) jlongArray {
	return jlongArray(unsafe.Pointer(C.NewLongArray((*C.JNIEnv)(env), C.jsize(len))))
}

func newFloatArray(env unsafe.Pointer, len jsize) jfloatArray {
	return jfloatArray(unsafe.Pointer(C.NewFloatArray((*C.JNIEnv)(env), C.jsize(len))))
}

func newDoubleArray(env unsafe.Pointer, len jsize) jdoubleArray {
	return jdoubleArray(unsafe.Pointer(C.NewDoubleArray((*C.JNIEnv)(env), C.jsize(len))))
}

func getBooleanArrayElements(env unsafe.Pointer, array jbooleanArray, isCopy unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.GetBooleanArrayElements((*C.JNIEnv)(env), C.jbooleanArray(unsafe.Pointer(array)), (*C.jboolean)(isCopy)))
}

func getByteArrayElements(env unsafe.Pointer, array jbyteArray, isCopy unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.GetByteArrayElements((*C.JNIEnv)(env), C.jbyteArray(unsafe.Pointer(array)), (*C.jboolean)(isCopy)))
}

func getCharArrayElements(env unsafe.Pointer, array jcharArray, isCopy unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.GetCharArrayElements((*C.JNIEnv)(env), C.jcharArray(unsafe.Pointer(array)), (*C.jboolean)(isCopy)))
}

func getShortArrayElements(env unsafe.Pointer, array jshortArray, isCopy unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.GetShortArrayElements((*C.JNIEnv)(env), C.jshortArray(unsafe.Pointer(array)), (*C.jboolean)(isCopy)))
}

func getIntArrayElements(env unsafe.Pointer, array jintArray, isCopy unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.GetIntArrayElements((*C.JNIEnv)(env), C.jintArray(unsafe.Pointer(array)), (*C.jboolean)(isCopy)))
}

func getLongArrayElements(env unsafe.Pointer, array jlongArray, isCopy unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.GetLongArrayElements((*C.JNIEnv)(env), C.jlongArray(unsafe.Pointer(array)), (*C.jboolean)(isCopy)))
}

func getFloatArrayElements(env unsafe.Pointer, array jfloatArray, isCopy unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.GetFloatArrayElements((*C.JNIEnv)(env), C.jfloatArray(unsafe.Pointer(array)), (*C.jboolean)(isCopy)))
}

func getDoubleArrayElements(env unsafe.Pointer, array jdoubleArray, isCopy unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.GetDoubleArrayElements((*C.JNIEnv)(env), C.jdoubleArray(unsafe.Pointer(array)), (*C.jboolean)(isCopy)))
}

func releaseBooleanArrayElements(env unsafe.Pointer, array jbooleanArray, elems unsafe.Pointer, mode jint) {
	C.ReleaseBooleanArrayElements((*C.JNIEnv)(env), C.jbooleanArray(unsafe.Pointer(array)), (*C.jboolean)(elems), C.jint(mode))
}

func releaseByteArrayElements(env unsafe.Pointer, array jbyteArray, elems unsafe.Pointer, mode jint) {
	C.ReleaseByteArrayElements((*C.JNIEnv)(env), C.jbyteArray(unsafe.Pointer(array)), (*C.jbyte)(elems), C.jint(mode))
}

func releaseCharArrayElements(env unsafe.Pointer, array jcharArray, elems unsafe.Pointer, mode jint) {
	C.ReleaseCharArrayElements((*C.JNIEnv)(env), C.jcharArray(unsafe.Pointer(array)), (*C.jchar)(elems), C.jint(mode))
}

func releaseShortArrayElements(env unsafe.Pointer, array jshortArray, elems unsafe.Pointer, mode jint) {
	C.ReleaseShortArrayElements((*C.JNIEnv)(env), C.jshortArray(unsafe.Pointer(array)), (*C.jshort)(elems), C.jint(mode))
}

func releaseIntArrayElements(env unsafe.Pointer, array jintArray, elems unsafe.Pointer, mode jint) {
	C.ReleaseIntArrayElements((*C.JNIEnv)(env), C.jintArray(unsafe.Pointer(array)), (*C.jint)(elems), C.jint(mode))
}

func releaseLongArrayElements(env unsafe.Pointer, array jlongArray, elems unsafe.Pointer, mode jint) {
	C.ReleaseLongArrayElements((*C.JNIEnv)(env), C.jlongArray(unsafe.Pointer(array)), (*C.jlong)(elems), C.jint(mode))
}

func releaseFloatArrayElements(env unsafe.Pointer, array jfloatArray, elems unsafe.Pointer, mode jint) {
	C.ReleaseFloatArrayElements((*C.JNIEnv)(env), C.jfloatArray(unsafe.Pointer(array)), (*C.jfloat)(elems), C.jint(mode))
}

func releaseDoubleArrayElements(env unsafe.Pointer, array jdoubleArray, elems unsafe.Pointer, mode jint) {
	C.ReleaseDoubleArrayElements((*C.JNIEnv)(env), C.jdoubleArray(unsafe.Pointer(array)), (*C.jdouble)(elems), C.jint(mode))
}

func getBooleanArrayRegion(env unsafe.Pointer, array jbooleanArray, start jsize, l jsize, buf unsafe.Pointer) {
	C.GetBooleanArrayRegion((*C.JNIEnv)(env), C.jbooleanArray(unsafe.Pointer(array)), C.jsize(start), C.jsize(l), (*C.jboolean)(buf))
}

func getByteArrayRegion(env unsafe.Pointer, array jbyteArray, start jsize, len jsize, buf unsafe.Pointer) {
	C.GetByteArrayRegion((*C.JNIEnv)(env), C.jbyteArray(unsafe.Pointer(array)), C.jsize(start), C.jsize(len), (*C.jbyte)(buf))
}

func getCharArrayRegion(env unsafe.Pointer, array jcharArray, start jsize, len jsize, buf unsafe.Pointer) {
	C.GetCharArrayRegion((*C.JNIEnv)(env), C.jcharArray(unsafe.Pointer(array)), C.jsize(start), C.jsize(len), (*C.jchar)(buf))
}

func getShortArrayRegion(env unsafe.Pointer, array jshortArray, start jsize, len jsize, buf unsafe.Pointer) {
	C.GetShortArrayRegion((*C.JNIEnv)(env), C.jshortArray(unsafe.Pointer(array)), C.jsize(start), C.jsize(len), (*C.jshort)(buf))
}

func getIntArrayRegion(env unsafe.Pointer, array jintArray, start jsize, len jsize, buf unsafe.Pointer) {
	C.GetIntArrayRegion((*C.JNIEnv)(env), C.jintArray(unsafe.Pointer(array)), C.jsize(start), C.jsize(len), (*C.jint)(buf))
}

func getLongArrayRegion(env unsafe.Pointer, array jlongArray, start jsize, len jsize, buf unsafe.Pointer) {
	C.GetLongArrayRegion((*C.JNIEnv)(env), C.jlongArray(unsafe.Pointer(array)), C.jsize(start), C.jsize(len), (*C.jlong)(buf))
}

func getFloatArrayRegion(env unsafe.Pointer, array jfloatArray, start jsize, len jsize, buf unsafe.Pointer) {
	C.GetFloatArrayRegion((*C.JNIEnv)(env), C.jfloatArray(unsafe.Pointer(array)), C.jsize(start), C.jsize(len), (*C.jfloat)(buf))
}

func getDoubleArrayRegion(env unsafe.Pointer, array jdoubleArray, start jsize, len jsize, buf unsafe.Pointer) {
	C.GetDoubleArrayRegion((*C.JNIEnv)(env), C.jdoubleArray(unsafe.Pointer(array)), C.jsize(start), C.jsize(len), (*C.jdouble)(buf))
}

func setBooleanArrayRegion(env unsafe.Pointer, array jbooleanArray, start jsize, l jsize, buf unsafe.Pointer) {
	C.SetBooleanArrayRegion((*C.JNIEnv)(env), C.jbooleanArray(unsafe.Pointer(array)), C.jsize(start), C.jsize(l), (*C.jboolean)(buf))
}

func setByteArrayRegion(env unsafe.Pointer, array jbyteArray, start jsize, len jsize, buf unsafe.Pointer) {
	C.SetByteArrayRegion((*C.JNIEnv)(env), C.jbyteArray(unsafe.Pointer(array)), C.jsize(start), C.jsize(len), (*C.jbyte)(buf))
}

func setCharArrayRegion(env unsafe.Pointer, array jcharArray, start jsize, len jsize, buf unsafe.Pointer) {
	C.SetCharArrayRegion((*C.JNIEnv)(env), C.jcharArray(unsafe.Pointer(array)), C.jsize(start), C.jsize(len), (*C.jchar)(buf))
}

func setShortArrayRegion(env unsafe.Pointer, array jshortArray, start jsize, len jsize, buf unsafe.Pointer) {
	C.SetShortArrayRegion((*C.JNIEnv)(env), C.jshortArray(unsafe.Pointer(array)), C.jsize(start), C.jsize(len), (*C.jshort)(buf))
}

func setIntArrayRegion(env unsafe.Pointer, array jintArray, start jsize, len jsize, buf unsafe.Pointer) {
	C.SetIntArrayRegion((*C.JNIEnv)(env), C.jintArray(unsafe.Pointer(array)), C.jsize(start), C.jsize(len), (*C.jint)(buf))
}

func setLongArrayRegion(env unsafe.Pointer, array jlongArray, start jsize, len jsize, buf unsafe.Pointer) {
	C.SetLongArrayRegion((*C.JNIEnv)(env), C.jlongArray(unsafe.Pointer(array)), C.jsize(start), C.jsize(len), (*C.jlong)(buf))
}

func setFloatArrayRegion(env unsafe.Pointer, array jfloatArray, start jsize, len jsize, buf unsafe.Pointer) {
	C.SetFloatArrayRegion((*C.JNIEnv)(env), C.jfloatArray(unsafe.Pointer(array)), C.jsize(start), C.jsize(len), (*C.jfloat)(buf))
}

func setDoubleArrayRegion(env unsafe.Pointer, array jdoubleArray, start jsize, len jsize, buf unsafe.Pointer) {
	C.SetDoubleArrayRegion((*C.JNIEnv)(env), C.jdoubleArray(unsafe.Pointer(array)), C.jsize(start), C.jsize(len), (*C.jdouble)(buf))
}

func registerNatives(env unsafe.Pointer, clazz jclass, methods unsafe.Pointer, nMethods jint) jint {
	return jint(C.RegisterNatives((*C.JNIEnv)(env), C.jclass(unsafe.Pointer(clazz)), (*C.JNINativeMethod)(methods), C.jint(nMethods)))
}

func getStringRegion(env unsafe.Pointer, str jstring, start jsize, len jsize, buf unsafe.Pointer) {
	C.GetStringRegion((*C.JNIEnv)(env), C.jstring(unsafe.Pointer(str)), C.jsize(start), C.jsize(len), (*C.jchar)(buf))
}

func getStringUTFRegion(env unsafe.Pointer, str jstring, start jsize, len jsize, buf unsafe.Pointer) {
	C.GetStringUTFRegion((*C.JNIEnv)(env), C.jstring(unsafe.Pointer(str)), C.jsize(start), C.jsize(len), (*C.char)(buf))
}

func getPrimitiveArrayCritical(env unsafe.Pointer, array jarray, isCopy unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.GetPrimitiveArrayCritical((*C.JNIEnv)(env), C.jarray(unsafe.Pointer(array)), (*C.jboolean)(isCopy)))
}

func releasePrimitiveArrayCritical(env unsafe.Pointer, array jarray, carray unsafe.Pointer, mode jint) {
	C.ReleasePrimitiveArrayCritical((*C.JNIEnv)(env), C.jarray(unsafe.Pointer(array)), (unsafe.Pointer)(carray), C.jint(mode))
}

func getStringCritical(env unsafe.Pointer, string jstring, isCopy unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(C.GetStringCritical((*C.JNIEnv)(env), C.jstring(unsafe.Pointer(string)), (*C.jboolean)(isCopy)))
}

func releaseStringCritical(env unsafe.Pointer, string jstring, cstring unsafe.Pointer) {
	C.ReleaseStringCritical((*C.JNIEnv)(env), C.jstring(unsafe.Pointer(string)), (*C.jchar)(cstring))
}

func attachCurrentThread(vm unsafe.Pointer, penv unsafe.Pointer, args unsafe.Pointer) jint {
	return jint(C.AttachCurrentThread((*C.JavaVM)(vm), (*unsafe.Pointer)(penv), (unsafe.Pointer)(args)))
}

func detachCurrentThread(vm unsafe.Pointer) jint {
	return jint(C.DetachCurrentThread((*C.JavaVM)(vm)))
}

func destroyJavaVM(vm unsafe.Pointer) jint {
	return jint(C.DestroyJavaVM((*C.JavaVM)(vm)))
}

func getJavaVM(env unsafe.Pointer, vm unsafe.Pointer) jint {
	return jint(C.GetJavaVM((*C.JNIEnv)(env), (**C.JavaVM)(vm)))
}

/* CallNonvirtual funcs... */

func callNonvirtualObjectMethodA(env unsafe.Pointer, obj jobject, clazz jclass, methodID jmethodID, args unsafe.Pointer) jobject {
	return jobject(unsafe.Pointer(C.CallNonvirtualObjectMethodA((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jclass(unsafe.Pointer(clazz)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args))))
}

func callNonvirtualBooleanMethodA(env unsafe.Pointer, obj jobject, clazz jclass, methodID jmethodID, args unsafe.Pointer) jboolean {
	return jboolean(C.CallNonvirtualBooleanMethodA((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jclass(unsafe.Pointer(clazz)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callNonvirtualByteMethodA(env unsafe.Pointer, obj jobject, clazz jclass, methodID jmethodID, args unsafe.Pointer) jbyte {
	return jbyte(C.CallNonvirtualByteMethodA((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jclass(unsafe.Pointer(clazz)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callNonvirtualCharMethodA(env unsafe.Pointer, obj jobject, clazz jclass, methodID jmethodID, args unsafe.Pointer) jchar {
	return jchar(C.CallNonvirtualCharMethodA((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jclass(unsafe.Pointer(clazz)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callNonvirtualShortMethodA(env unsafe.Pointer, obj jobject, clazz jclass, methodID jmethodID, args unsafe.Pointer) jshort {
	return jshort(C.CallNonvirtualShortMethodA((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jclass(unsafe.Pointer(clazz)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callNonvirtualIntMethodA(env unsafe.Pointer, obj jobject, clazz jclass, methodID jmethodID, args unsafe.Pointer) jint {
	return jint(C.CallNonvirtualIntMethodA((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jclass(unsafe.Pointer(clazz)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callNonvirtualLongMethodA(env unsafe.Pointer, obj jobject, clazz jclass, methodID jmethodID, args unsafe.Pointer) jlong {
	return jlong(C.CallNonvirtualLongMethodA((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jclass(unsafe.Pointer(clazz)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callNonvirtualFloatMethodA(env unsafe.Pointer, obj jobject, clazz jclass, methodID jmethodID, args unsafe.Pointer) jfloat {
	return jfloat(C.CallNonvirtualFloatMethodA((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jclass(unsafe.Pointer(clazz)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callNonvirtualDoubleMethodA(env unsafe.Pointer, obj jobject, clazz jclass, methodID jmethodID, args unsafe.Pointer) jdouble {
	return jdouble(C.CallNonvirtualDoubleMethodA((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jclass(unsafe.Pointer(clazz)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args)))
}

func callNonvirtualVoidMethodA(env unsafe.Pointer, obj jobject, clazz jclass, methodID jmethodID, args unsafe.Pointer) {
	C.CallNonvirtualVoidMethodA((*C.JNIEnv)(env), C.jobject(unsafe.Pointer(obj)), C.jclass(unsafe.Pointer(clazz)), C.jmethodID(unsafe.Pointer(methodID)), (*C.jvalue)(args))
}
