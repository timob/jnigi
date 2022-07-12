// Copyright 2016 Tim O'Brien. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	JNIGI (Java Native Interface Go Interface)

	A package to access Java from Go code.

	All constructor and method call functions convert parameter arguments and return values.

	Arguments are converted from Go to Java if:
	  - The type is Go built in type and there is an equivalent Java primitive type.
	  - The type is a slice of such a Go built in type.
	  - The type implements the ToJavaConverter interface
	Return values are converted from Java to Go if:
	  - The type is a Java primitive type.
	  - The type is a Java array of a primitive type.
	  - The type implements the ToGoConverter interface


	Go Builtin to/from Java Primitive:

		bool			Boolean
		byte			Byte
		int16			Short
		uint16			Char
		int				Int (also int32 -> Int)
		int64			Long
		float32			Float
		float64			Double

*/
package jnigi

import (
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

// copy arguments in to C memory before passing to jni functions
var copyToC bool = false

func toBool(b jboolean) bool {
	return b == 1
}

func fromBool(b bool) jboolean {
	if b {
		return 1
	} else {
		return 0
	}
}

// ObjectRef holds a reference to a Java Object
type ObjectRef struct {
	jobject   jobject
	className string
	isArray   bool
}

// NewObjectRef returns new *ObjectRef with Nil JNI object reference and class name set to className.
func NewObjectRef(className string) *ObjectRef {
	return &ObjectRef{0, className, false}
}

// NewObjectArrayRef returns new *ObjectRef with Nil JNI object array reference and class name set to className.
func NewObjectArrayRef(className string) *ObjectRef {
	return &ObjectRef{0, className, true}
}

// WrapJObject wraps a JNI object value in an ObjectRef
func WrapJObject(jobj uintptr, className string, isArray bool) *ObjectRef {
	return &ObjectRef{jobject(jobj), className, isArray}
}

// GetClassName returns class name of object reference.
func (o *ObjectRef) GetClassName() string {
	return o.className
}

// IsArray returns true if reference is to object array.
func (o *ObjectRef) IsArray() bool {
	return o.isArray
}

// Cast return a new *CastedObjectRef containing the receiver with casted class name set to className.
func (o *ObjectRef) Cast(className string) *CastedObjectRef {
	return &CastedObjectRef{o, className}
}

// IsNil is true if ObjectRef has a Nil Java value
func (o *ObjectRef) IsNil() bool {
	return o.jobject == 0
}

// IsInstanceOf returns true if o is an instance of className
func (o *ObjectRef) IsInstanceOf(env *Env, className string) (bool, error) {
	class, err := env.callFindClass(className)
	if err != nil {
		return false, err
	}
	return toBool(isInstanceOf(env.jniEnv, o.jobject, class)), nil
}

func (o *ObjectRef) jobj() jobject {
	return o.jobject
}

// JObject gets JNI object value of o
func (o *ObjectRef) JObject() jobject {
	return o.jobj()
}

type jobj interface {
	jobj() jobject
}

// CastedObjectRef represents an object reference casted to a super class.
// This is used to create method signatures for generic classes.
type CastedObjectRef struct {
	*ObjectRef
	Cast string
}

// GetClassName returns class name of the cast.
func (c *CastedObjectRef) GetClassName() string {
	return c.Cast
}

// ExceptionHandler is used to convert a thrown exception (java.lang.Throwable) to a Go error.
type ExceptionHandler interface {
	CatchException(env *Env, exception *ObjectRef) error
}

// ExceptionHandlerFunc is an adapter to allow use of ordinary functions as an
// ExceptionHandler. If f is a function with the appropriate signature, ExceptionHandlerFunc(f)
// is an ExceptionHandler object that calls f.
type ExceptionHandlerFunc func(env *Env, exception *ObjectRef) error

// CatchException calls f to implement ExceptionHandler.
func (f ExceptionHandlerFunc) CatchException(env *Env, exception *ObjectRef) error {
	return f(env, exception)
}

// Env holds a JNIEnv value. Methods in this package often require an *Env pointer to specify
// the JNI Env to run in, so it might be good to store the Env as a package variable.
type Env struct {
	jniEnv           unsafe.Pointer
	preCalcSig       string
	classCache       map[string]jclass
	ExceptionHandler ExceptionHandler
}

// WrapEnv wraps an JNI Env value in an Env
func WrapEnv(envPtr unsafe.Pointer) *Env {
	return &Env{jniEnv: envPtr, classCache: make(map[string]jclass)}
}

// JVM holds a JavaVM value, you only need one of these in your app.
type JVM struct {
	javaVM unsafe.Pointer
}

// JVMInitArgs holds a JavaVMInitArgs value
type JVMInitArgs struct {
	javaVMInitArgs unsafe.Pointer
}

// CreateJVM calls JNI CreateJavaVM and returns references to the JVM and the initial environment.
// Use NewJVMInitArgs to create jvmInitArgs.
//
// On Android, call AndroidJVM to return references to the Android Application's JVM instead.
//
// Must call runtime.LockOSThread() first.
func CreateJVM(jvmInitArgs *JVMInitArgs) (*JVM, *Env, error) {
	p := malloc(unsafe.Sizeof((unsafe.Pointer)(nil)))
	p2 := malloc(unsafe.Sizeof((unsafe.Pointer)(nil)))

	if jni_CreateJavaVM(p2, p, jvmInitArgs.javaVMInitArgs) < 0 {
		return nil, nil, errors.New("Couldn't instantiate JVM")
	}
	jvm := &JVM{*(*unsafe.Pointer)(p2)}
	env := &Env{jniEnv: *(*unsafe.Pointer)(p), classCache: make(map[string]jclass)}

	free(p)
	free(p2)
	return jvm, env, nil
}

// AttachCurrentThread calls JNI AttachCurrentThread.
// Must call runtime.LockOSThread() first.
func (j *JVM) AttachCurrentThread() *Env {
	p := malloc(unsafe.Sizeof((unsafe.Pointer)(nil)))

	//	p := (**C.JNIEnv)(malloc(unsafe.Sizeof((*C.JNIEnv)(nil))))

	if attachCurrentThread(j.javaVM, p, nil) < 0 {
		panic("AttachCurrentThread failed")
	}

	return &Env{jniEnv: *(*unsafe.Pointer)(p), classCache: make(map[string]jclass)}
}

// DetachCurrentThread calls JNI DetachCurrentThread
func (j *JVM) DetachCurrentThread() error {
	if detachCurrentThread(j.javaVM) < 0 {
		return errors.New("JNIGI: detachCurrentThread error")
	}
	return nil
}

// Destroy calls JNI DestroyJavaVM
func (j *JVM) Destroy() error {
	if destroyJavaVM(j.javaVM) < 0 {
		return errors.New("JNIGI: destroyJavaVM error")
	}
	return nil
}

// GetJVM Calls JNI GetJavaVM. Needs to be called on an attached thread.
func (j *Env) GetJVM() (*JVM, error) {
	p := malloc(unsafe.Sizeof((unsafe.Pointer)(nil)))

	if getJavaVM(j.jniEnv, p) < 0 {
		return nil, errors.New("Couldn't get JVM")
	}

	jvm := &JVM{*(*unsafe.Pointer)(p)}

	free(p)

	return jvm, nil
}

func (j *Env) exceptionCheck() bool {
	return toBool(exceptionCheck(j.jniEnv))
}

func (j *Env) describeException() {
	exceptionDescribe(j.jniEnv)
}

func (j *Env) handleException() error {

	e := exceptionOccurred(j.jniEnv)
	if e == 0 {
		return errors.New("Java JNI function returned error but JNI indicates no current exception")
	}

	defer deleteLocalRef(j.jniEnv, jobject(e))

	ref := WrapJObject(uintptr(e), "java/lang/Throwable", false)

	if j.ExceptionHandler == nil {
		return DefaultExceptionHandler.CatchException(j, ref)
	}

	// Temporarily disable handler in the event exception rises during handling.
	// By setting it to the DescribeExceptionHandler, exceptions will get printed
	// and cleared.
	handler := j.ExceptionHandler
	j.ExceptionHandler = DescribeExceptionHandler
	defer func() {
		j.ExceptionHandler = handler
	}()

	return handler.CatchException(j, ref)
}

// NewObject calls JNI NewObjectA, className class name of new object, args arguments to constructor.
func (j *Env) NewObject(className string, args ...interface{}) (*ObjectRef, error) {
	class, err := j.callFindClass(className)
	if err != nil {
		return nil, err
	}

	if err := replaceConvertedArgs(args); err != nil {
		return nil, err
	}
	var methodSig string
	if j.preCalcSig != "" {
		methodSig = j.preCalcSig
		j.preCalcSig = ""
	} else {
		calcSig, err := sigForMethod(Void, "", args)
		if err != nil {
			return nil, err
		}
		methodSig = calcSig
	}

	mid, err := j.callGetMethodID(false, class, "<init>", methodSig)
	if err != nil {
		return nil, err
	}

	// create args for jni call
	jniArgs, refs, err := j.createArgs(args)
	if err != nil {
		return nil, err
	}
	defer func() {
		cleanUpArgs(jniArgs)
		for _, ref := range refs {
			deleteLocalRef(j.jniEnv, ref)
		}
	}()

	obj := newObjectA(j.jniEnv, class, mid, jniArgs)
	if obj == 0 {
		return nil, j.handleException()
	}

	return &ObjectRef{obj, className, false}, nil
}

func (j *Env) callFindClass(className string) (jclass, error) {
	if v, ok := j.classCache[className]; ok {
		return v, nil
	}
	cnCstr := cString(className)
	defer free(cnCstr)
	class := findClass(j.jniEnv, cnCstr)
	if class == 0 {
		return 0, j.handleException()
	}
	ref := newGlobalRef(j.jniEnv, jobject(class))
	deleteLocalRef(j.jniEnv, jobject(class))
	j.classCache[className] = jclass(ref)

	return jclass(ref), nil
}

func (j *Env) callGetMethodID(static bool, class jclass, name, sig string) (jmethodID, error) {
	mnCstr := cString(name)
	defer free(mnCstr)

	sigCstr := cString(sig)
	defer free(sigCstr)

	var mid jmethodID
	if static {
		mid = getStaticMethodID(j.jniEnv, class, mnCstr, sigCstr)
	} else {
		mid = getMethodID(j.jniEnv, class, mnCstr, sigCstr)
	}
	//	fmt.Printf("sig = %s\n", sig)
	if mid == 0 {
		return 0, j.handleException()
	}

	return mid, nil
}

// PrecalculateSignature sets the signature of the next call to sig, disables automatic signature building.
func (j *Env) PrecalculateSignature(sig string) {
	j.preCalcSig = sig
}

const big = 1024 * 1024 * 100

// FromObjectArray converts an Java array of objects objRef in to a slice of *ObjectRef which is returned.
func (j *Env) FromObjectArray(objRef *ObjectRef) []*ObjectRef {
	len := int(getArrayLength(j.jniEnv, jarray(objRef.jobject)))
	// exception check?

	v := make([]*ObjectRef, len)
	for i := 0; i < len; i++ {
		jobj := getObjectArrayElement(j.jniEnv, jobjectArray(objRef.jobject), jsize(i))
		if j.exceptionCheck() {
			panic(j.handleException())
		}
		v[i] = &ObjectRef{jobj, objRef.className, false}
	}

	return v
}

func (j *Env) toGoArray(array jobject, aType Type) (interface{}, error) {
	len := int(getArrayLength(j.jniEnv, jarray(array)))
	// exception check?

	switch aType.baseType() {
	case Boolean:
		v := make([]bool, len)
		if len >= 0 {
			ptr := getBooleanArrayElements(j.jniEnv, jbooleanArray(array), nil)
			if j.exceptionCheck() {
				return nil, j.handleException()
			}
			elems := (*(*[big]byte)(ptr))[0:len]
			for i := 0; i < len; i++ {
				v[i] = (elems[i] == 1)
			}
			releaseBooleanArrayElements(j.jniEnv, jbooleanArray(array), ptr, jint(jni_abort))
		}
		return v, nil
	case Byte:
		v := make([]byte, len)
		if len >= 0 {
			ptr := getByteArrayElements(j.jniEnv, jbyteArray(array), nil)
			if j.exceptionCheck() {
				return nil, j.handleException()
			}
			elems := (*(*[big]byte)(ptr))[0:len]
			copy(v, elems)
			releaseByteArrayElements(j.jniEnv, jbyteArray(array), ptr, jint(jni_abort))
		}
		return v, nil
	case Short:
		v := make([]int16, len)
		if len >= 0 {
			ptr := getShortArrayElements(j.jniEnv, jshortArray(array), nil)
			if j.exceptionCheck() {
				return nil, j.handleException()
			}
			elems := (*(*[big]int16)(ptr))[0:len]
			copy(v, elems)
			releaseShortArrayElements(j.jniEnv, jshortArray(array), ptr, jint(jni_abort))
		}
		return v, nil
	case Char:
		v := make([]uint16, len)
		if len >= 0 {
			ptr := getCharArrayElements(j.jniEnv, jcharArray(array), nil)
			if j.exceptionCheck() {
				return nil, j.handleException()
			}
			elems := (*(*[big]uint16)(ptr))[0:len]
			copy(v, elems)
			releaseCharArrayElements(j.jniEnv, jcharArray(array), ptr, jint(jni_abort))
		}
		return v, nil
	case Int:
		v := make([]int, len)
		if len >= 0 {
			ptr := getIntArrayElements(j.jniEnv, jintArray(array), nil)
			if j.exceptionCheck() {
				return nil, j.handleException()
			}
			elems := (*(*[big]int32)(ptr))[0:len]
			//copy(v, elems)
			for i := 0; i < len; i++ {
				v[i] = int(elems[i])
			}
			releaseIntArrayElements(j.jniEnv, jintArray(array), ptr, jint(jni_abort))
		}
		return v, nil
	case Long:
		v := make([]int64, len)
		if len >= 0 {
			ptr := getLongArrayElements(j.jniEnv, jlongArray(array), nil)
			if j.exceptionCheck() {
				return nil, j.handleException()
			}
			elems := (*(*[big]int64)(ptr))[0:len]
			copy(v, elems)
			releaseLongArrayElements(j.jniEnv, jlongArray(array), ptr, jint(jni_abort))
		}
		return v, nil
	case Float:
		v := make([]float32, len)
		if len >= 0 {
			ptr := getFloatArrayElements(j.jniEnv, jfloatArray(array), nil)
			if j.exceptionCheck() {
				return nil, j.handleException()
			}
			elems := (*(*[big]float32)(ptr))[0:len]
			copy(v, elems)
			releaseFloatArrayElements(j.jniEnv, jfloatArray(array), ptr, jint(jni_abort))
		}
		return v, nil
	case Double:
		v := make([]float64, len)
		if len >= 0 {
			ptr := getDoubleArrayElements(j.jniEnv, jdoubleArray(array), nil)
			if j.exceptionCheck() {
				return nil, j.handleException()
			}
			elems := (*(*[big]float64)(ptr))[0:len]
			copy(v, elems)
			releaseDoubleArrayElements(j.jniEnv, jdoubleArray(array), ptr, jint(jni_abort))
		}
		return v, nil
	default:
		return nil, errors.New("JNIGI unsupported array type")
	}
}

// ToObjectArray converts slice of ObjectRef objRefs of class name className, to a new Java object
// array returning a reference to this array.
func (j *Env) ToObjectArray(objRefs []*ObjectRef, className string) (arrayRef *ObjectRef) {
	arrayRef = &ObjectRef{className: className, isArray: true}
	class, err := j.callFindClass(className)
	if err != nil {
		j.describeException()
		exceptionClear(j.jniEnv)
		return
	}

	oa := newObjectArray(j.jniEnv, jsize(len(objRefs)), class, 0)
	if oa == 0 {
		panic(j.handleException())
	}
	arrayRef.jobject = jobject(oa)

	for i, obj := range objRefs {
		setObjectArrayElement(j.jniEnv, oa, jsize(i), obj.jobject)
		if j.exceptionCheck() {
			j.describeException()
			exceptionClear(j.jniEnv)
		}
	}
	return
}

// ByteArray holds a JNI JbyteArray
type ByteArray struct {
	arr jbyteArray
	n   int
}

// NewByteArray calls JNI NewByteArray
func (j *Env) NewByteArray(n int) *ByteArray {
	a := newByteArray(j.jniEnv, jsize(n))
	return &ByteArray{a, n}
}

// NewByteArrayFromSlice calls JNI NewByteArray and GetCritical, copies src to byte array,
// calls JNI Release Critical. Returns new byte array.
func (j *Env) NewByteArrayFromSlice(src []byte) *ByteArray {
	b := j.NewByteArray(len(src))
	if len(src) > 0 {
		bytes := b.GetCritical(j)
		copy(bytes, src)
		b.ReleaseCritical(j, bytes)
	}
	return b
}

// NewByteArrayFromObject creates new ByteArray and sets it from ObjectRef o.
func (j *Env) NewByteArrayFromObject(o *ObjectRef) *ByteArray {
	ba := &ByteArray{}
	ba.SetObject(o)
	ba.n = int(getArrayLength(j.jniEnv, jarray(ba.arr)))
	return ba
}

func (b *ByteArray) jobj() jobject {
	return jobject(b.arr)
}

func (b *ByteArray) getType() Type {
	return Byte | Array
}

// GetCritical calls JNI GetPrimitiveArrayCritical
func (b *ByteArray) GetCritical(env *Env) []byte {
	if b.n == 0 {
		return nil
	}
	ptr := getPrimitiveArrayCritical(env.jniEnv, jarray(b.arr), nil)
	return (*(*[big]byte)(ptr))[0:b.n]
}

// GetCritical calls JNI ReleasePrimitiveArrayCritical
func (b *ByteArray) ReleaseCritical(env *Env, bytes []byte) {
	if len(bytes) == 0 {
		return
	}
	ptr := unsafe.Pointer(&bytes[0])
	releasePrimitiveArrayCritical(env.jniEnv, jarray(b.arr), ptr, 0)
}

// GetObject returns byte array as *ObjectRef.
func (b *ByteArray) GetObject() *ObjectRef {
	return &ObjectRef{jobject(b.arr), "java/lang/Object", false}
}

// SetObject sets byte array from o.
func (b *ByteArray) SetObject(o *ObjectRef) {
	b.arr = jbyteArray(o.jobject)
}

// CopyBytes creates a go slice of bytes of same length as byte array, calls GetCritical,
// copies byte array into go slice, calls ReleaseCritical, returns go slice.
func (b *ByteArray) CopyBytes(env *Env) []byte {
	r := make([]byte, b.n)
	src := b.GetCritical(env)
	copy(r, src)
	b.ReleaseCritical(env, src)
	return r
}

// this copies slice contents in to C memory before passing this pointer to JNI array function
// if copy var is set to true
func (j *Env) toJavaArray(src interface{}) (jobject, error) {
	switch v := src.(type) {
	case []bool:
		ba := newBooleanArray(j.jniEnv, jsize(len(v)))
		if ba == 0 {
			return 0, j.handleException()
		}
		if len(v) == 0 {
			return jobject(ba), nil
		}
		src := make([]byte, len(v))
		for i, vset := range v {
			if vset {
				src[i] = 1
			}
		}
		var ptr unsafe.Pointer
		if copyToC {
			ptr = malloc(uintptr(len(v)))
			defer free(ptr)
			data := (*(*[big]byte)(ptr))[:len(v)]
			copy(data, src)
		} else {
			ptr = unsafe.Pointer(&src[0])
		}
		setBooleanArrayRegion(j.jniEnv, ba, jsize(0), jsize(len(v)), ptr)
		if j.exceptionCheck() {
			return 0, j.handleException()
		}
		return jobject(ba), nil
	case []byte:
		ba := newByteArray(j.jniEnv, jsize(len(v)))
		if ba == 0 {
			return 0, j.handleException()
		}
		if len(v) == 0 {
			return jobject(ba), nil
		}
		var ptr unsafe.Pointer
		if copyToC {
			ptr = malloc(uintptr(len(v)))
			defer free(ptr)
			data := (*(*[big]byte)(ptr))[:len(v)]
			copy(data, v)
		} else {
			ptr = unsafe.Pointer(&v[0])
		}
		setByteArrayRegion(j.jniEnv, ba, jsize(0), jsize(len(v)), ptr)
		if j.exceptionCheck() {
			return 0, j.handleException()
		}
		return jobject(ba), nil
	case []int16:
		array := newShortArray(j.jniEnv, jsize(len(v)))
		if array == 0 {
			return 0, j.handleException()
		}
		if len(v) == 0 {
			return jobject(array), nil
		}
		var ptr unsafe.Pointer
		if copyToC {
			ptr = malloc(unsafe.Sizeof(int16(0)) * uintptr(len(v)))
			defer free(ptr)
			data := (*(*[big]int16)(ptr))[:len(v)]
			copy(data, v)
		} else {
			ptr = unsafe.Pointer(&v[0])
		}
		setShortArrayRegion(j.jniEnv, array, jsize(0), jsize(len(v)), ptr)
		if j.exceptionCheck() {
			return 0, j.handleException()
		}
		return jobject(array), nil
	case []uint16:
		array := newCharArray(j.jniEnv, jsize(len(v)))
		if array == 0 {
			return 0, j.handleException()
		}
		if len(v) == 0 {
			return jobject(array), nil
		}
		var ptr unsafe.Pointer
		if copyToC {
			ptr = malloc(unsafe.Sizeof(uint16(0)) * uintptr(len(v)))
			defer free(ptr)
			data := (*(*[big]uint16)(ptr))[:len(v)]
			copy(data, v)
		} else {
			ptr = unsafe.Pointer(&v[0])
		}
		setCharArrayRegion(j.jniEnv, array, jsize(0), jsize(len(v)), ptr)
		if j.exceptionCheck() {
			return 0, j.handleException()
		}
		return jobject(array), nil
	case []int32:
		array := newIntArray(j.jniEnv, jsize(len(v)))
		if array == 0 {
			return 0, j.handleException()
		}
		if len(v) == 0 {
			return jobject(array), nil
		}
		var ptr unsafe.Pointer
		if copyToC {
			ptr = malloc(unsafe.Sizeof(int32(0)) * uintptr(len(v)))
			defer free(ptr)
			data := (*(*[big]int32)(ptr))[:len(v)]
			copy(data, v)
		} else {
			ptr = unsafe.Pointer(&v[0])
		}
		setIntArrayRegion(j.jniEnv, array, jsize(0), jsize(len(v)), ptr)
		if j.exceptionCheck() {
			return 0, j.handleException()
		}
		return jobject(array), nil
	case []int:
		array := newIntArray(j.jniEnv, jsize(len(v)))
		if array == 0 {
			return 0, j.handleException()
		}
		if len(v) == 0 {
			return jobject(array), nil
		}
		var ptr unsafe.Pointer
		if copyToC {
			ptr = malloc(unsafe.Sizeof(int32(0)) * uintptr(len(v)))
			defer free(ptr)
			data := (*(*[big]int32)(ptr))[:len(v)]
			//copy(data, v)
			for i := 0; i < len(data); i++ {
				data[i] = int32(v[i])
			}
		} else {
			data := make([]int32, len(v))
			for i := 0; i < len(v); i++ {
				data[i] = int32(v[i])
			}
			ptr = unsafe.Pointer(&data[0])
		}
		setIntArrayRegion(j.jniEnv, array, jsize(0), jsize(len(v)), ptr)
		if j.exceptionCheck() {
			return 0, j.handleException()
		}
		return jobject(array), nil
	case []int64:
		array := newLongArray(j.jniEnv, jsize(len(v)))
		if array == 0 {
			return 0, j.handleException()
		}
		if len(v) == 0 {
			return jobject(array), nil
		}
		var ptr unsafe.Pointer
		if copyToC {
			ptr = malloc(unsafe.Sizeof(int64(0)) * uintptr(len(v)))
			defer free(ptr)
			data := (*(*[big]int64)(ptr))[:len(v)]
			copy(data, v)
		} else {
			ptr = unsafe.Pointer(&v[0])
		}
		setLongArrayRegion(j.jniEnv, array, jsize(0), jsize(len(v)), ptr)
		if j.exceptionCheck() {
			return 0, j.handleException()
		}
		return jobject(array), nil
	case []float32:
		array := newFloatArray(j.jniEnv, jsize(len(v)))
		if array == 0 {
			return 0, j.handleException()
		}
		if len(v) == 0 {
			return jobject(array), nil
		}
		var ptr unsafe.Pointer
		if copyToC {
			ptr = malloc(unsafe.Sizeof(float32(0)) * uintptr(len(v)))
			defer free(ptr)
			data := (*(*[big]float32)(ptr))[:len(v)]
			copy(data, v)
		} else {
			ptr = unsafe.Pointer(&v[0])
		}
		setFloatArrayRegion(j.jniEnv, array, jsize(0), jsize(len(v)), ptr)
		if j.exceptionCheck() {
			return 0, j.handleException()
		}
		return jobject(array), nil
	case []float64:
		array := newDoubleArray(j.jniEnv, jsize(len(v)))
		if array == 0 {
			return 0, j.handleException()
		}
		if len(v) == 0 {
			return jobject(array), nil
		}
		var ptr unsafe.Pointer
		if copyToC {
			ptr = malloc(unsafe.Sizeof(float64(0)) * uintptr(len(v)))
			defer free(ptr)
			data := (*(*[big]float64)(ptr))[:len(v)]
			copy(data, v)
		} else {
			ptr = unsafe.Pointer(&v[0])
		}
		setDoubleArrayRegion(j.jniEnv, array, jsize(0), jsize(len(v)), ptr)
		if j.exceptionCheck() {
			return 0, j.handleException()
		}
		return jobject(array), nil
	default:
		return 0, errors.New("JNIGI unsupported array type")
	}
}

// pointer should be freed, refs should be deleted
// jvalue 64 bit
func (j *Env) createArgs(args []interface{}) (ptr unsafe.Pointer, refs []jobject, err error) {
	if len(args) == 0 {
		return nil, nil, nil
	}

	argList := make([]uint64, len(args))
	refs = make([]jobject, 0)

	for i, arg := range args {
		switch v := arg.(type) {
		case *convertedArg:
			argList[i] = uint64(v.ObjectRef.jobject)
			refs = append(refs, v.ObjectRef.jobject)
		case jobj:
			argList[i] = uint64(v.jobj())
		case bool:
			if v {
				argList[i] = uint64(jboolean(1))
			} else {
				argList[i] = uint64(jboolean(0))
			}
		case byte:
			argList[i] = uint64(jbyte(v))
		case uint16:
			argList[i] = uint64(jchar(v))
		case int16:
			argList[i] = uint64(jshort(v))
		case int32:
			argList[i] = uint64(jint(v))
		case int:
			argList[i] = uint64(jint(int32(v)))
		case int64:
			argList[i] = uint64(jlong(v))
		case float32:
			argList[i] = uint64(jfloat(v))
		case float64:
			argList[i] = uint64(jdouble(v))
		case []bool, []byte, []int16, []uint16, []int32, []int, []int64, []float32, []float64:
			if array, arrayErr := j.toJavaArray(v); arrayErr == nil {
				argList[i] = uint64(array)
				refs = append(refs, array)
			} else {
				err = arrayErr
			}
		default:
			err = fmt.Errorf("JNIGI: argument not a valid value %T (%v)", args[i], args[i])
		}

		if err != nil {
			break
		}
	}

	if err != nil {
		for _, ref := range refs {
			deleteLocalRef(j.jniEnv, ref)
		}
		refs = nil
		return
	}

	if copyToC {
		ptr = malloc(unsafe.Sizeof(uint64(0)) * uintptr(len(args)))
		data := (*(*[big]uint64)(ptr))[:len(args)]
		copy(data, argList)
	} else {
		ptr = unsafe.Pointer(&argList[0])
	}
	return
}

// TypeSpec is implemented by Type, ObjectType and ObjectArrayType.
type TypeSpec interface {
	internal()
}

// Type is used to specify return types and field types. Array value can be ORed with primitive type.
// Implements TypeSpec. See package constants for values.
type Type uint32

const (
	Void = Type(1 << iota)
	Boolean
	Byte
	Char
	Short
	Int
	Long
	Float
	Double
	Object
	Array
)

func (t Type) baseType() Type {
	return t &^ Array
}

func (t Type) isArray() bool {
	return t&Array > 0
}

func (t Type) internal() {}

// ObjectType is treated as Object Type. It's value is used to specify the class of the object.
// For example jnigi.ObjectType("java/lang/string").
// Implements TypeSpec.
type ObjectType string

func (o ObjectType) internal() {}

// ObjectArrayType is treated as Object | Array Type. It's value specify the class of the elements.
// Implements TypeSpec.
type ObjectArrayType string

func (o ObjectArrayType) internal() {}

type convertedArray interface {
	getType() Type
}

// Allow return types to be a string that specifies an object type. This is to
// retain compatiblity.
func typeOfReturnValue(value interface{}) (t Type, className string, err error) {
	if v, ok := value.(string); ok {
		return typeOfValue(ObjectType(v))
	}
	return typeOfValue(value)
}

// ClassInfoGetter is implemented by *ObjectRef, *CastedObjectRef to get type info for object values.
type ClassInfoGetter interface {
	GetClassName() string
	IsArray() bool
}

// TypeGetter can be implemented to control which primitive type a value is treated as.
type TypeGetter interface {
	GetType() Type
}

func typeOfValue(value interface{}) (t Type, className string, err error) {
	switch v := value.(type) {
	case Type:
		t = v
		if t.baseType() == Object {
			className = "java/lang/Object"
		}
	case ObjectType:
		t = Object
		className = string(v)
	case ObjectArrayType:
		t = Object | Array
		className = string(v)
	case TypeGetter:
		t = v.GetType()

	// This is implemented by *ObjectRef, *CastedObjectRef, *convertedArg
	case ClassInfoGetter:
		t = Object
		if v.IsArray() {
			t = t | Array
		}
		className = v.GetClassName()
	case bool, *bool:
		t = Boolean
	case byte, *byte:
		t = Byte
	case int16, *int16:
		t = Short
	case uint16, *uint16:
		t = Char
	case int32, *int32:
		t = Int
	case int, *int:
		t = Int
	case int64, *int64:
		t = Long
	case float32, *float32:
		t = Float
	case float64, *float64:
		t = Double
	case []bool, *[]bool:
		t = Boolean | Array
		className = "java/lang/Object"
	case []byte, *[]byte:
		t = Byte | Array
		className = "java/lang/Object"
	case []uint16, *[]uint16:
		t = Char | Array
		className = "java/lang/Object"
	case []int16, *[]int16:
		t = Short | Array
		className = "java/lang/Object"
	case []int32, *[]int32:
		t = Int | Array
		className = "java/lang/Object"
	case []int, *[]int:
		t = Int | Array
		className = "java/lang/Object"
	case []int64, *[]int64:
		t = Long | Array
		className = "java/lang/Object"
	case []float32, *[]float32:
		t = Float | Array
		className = "java/lang/Object"
	case []float64, *[]float64:
		t = Double | Array
		className = "java/lang/Object"
	case convertedArray:
		t = v.getType()
		className = "java/lang/Object"
	default:
		err = fmt.Errorf("JNIGI: unknown type %T (value = %v)", v, v)
	}
	return
}

func typeSignature(t Type, className string) (sig string) {
	if t.isArray() {
		sig = "["
	}
	base := t.baseType()
	switch {
	case base == Object:
		sig += "L" + className + ";"
	case base == Void:
		sig += "V"
	case base == Boolean:
		sig += "Z"
	case base == Byte:
		sig += "B"
	case base == Char:
		sig += "C"
	case base == Short:
		sig += "S"
	case base == Int:
		sig += "I"
	case base == Long:
		sig += "J"
	case base == Float:
		sig += "F"
	case base == Double:
		sig += "D"
	}
	return
}

func sigForMethod(returnType Type, returnClass string, args []interface{}) (string, error) {
	var paramStr string
	for i := range args {
		t, c, err := typeOfValue(args[i])
		if err != nil {
			return "", err
		}
		paramStr += typeSignature(t, c)
	}
	return fmt.Sprintf("(%s)%s", paramStr, typeSignature(returnType, returnClass)), nil
}

func cleanUpArgs(ptr unsafe.Pointer) {
	if copyToC {
		free(ptr)
	}
}

func (o *ObjectRef) getClass(env *Env) (class jclass, err error) {
	class, err = env.callFindClass(o.className)
	if err != nil {
		return 0, err
	}

	// if object is java/lang/Object try to up class it
	// there is an odd way to get the class name see: http://stackoverflow.com/questions/12719766/can-i-know-the-name-of-the-class-that-calls-a-jni-c-method
	if o.className == "java/lang/Object" {
		mid, err := env.callGetMethodID(false, class, "getClass", "()Ljava/lang/Class;")
		if err != nil {
			return 0, err
		}
		obj := callObjectMethodA(env.jniEnv, o.jobject, mid, nil)
		if env.exceptionCheck() {
			return 0, env.handleException()
		}
		defer deleteLocalRef(env.jniEnv, obj)
		objClass := getObjectClass(env.jniEnv, obj)
		if objClass == 0 {
			return 0, env.handleException()
		}
		defer deleteLocalRef(env.jniEnv, jobject(objClass))
		mid, err = env.callGetMethodID(false, objClass, "getName", "()Ljava/lang/String;")
		if err != nil {
			return 0, err
		}
		obj2 := callObjectMethodA(env.jniEnv, obj, mid, nil)
		if env.exceptionCheck() {
			return 0, env.handleException()
		}
		strObj := WrapJObject(uintptr(obj2), "java/lang/String", false)
		if strObj.IsNil() {
			return 0, errors.New("unexpected error getting object class name")
		}
		defer env.DeleteLocalRef(strObj)
		var b []byte
		if err := strObj.CallMethod(env, "getBytes", &b, env.GetUTF8String()); err != nil {
			return 0, err
		}
		gotClass := string(b)

		// note uses . for class name separator
		if gotClass != "java.lang.Object" {
			gotClass = strings.Replace(gotClass, ".", "/", -1)
			class, err = env.callFindClass(gotClass)
			if err != nil {
				return 0, err
			}
			o.className = gotClass
			return class, err
		}
	}
	return
}

// CallMethod calls method methodName on o with arguments args and stores return value in dest.
func (o *ObjectRef) CallMethod(env *Env, methodName string, dest interface{}, args ...interface{}) error {
	rType, rClassName, err := typeOfReturnValue(dest)
	if err != nil {
		return err
	}

	retVal, err := o.genericCallMethod(env, methodName, rType, rClassName, args...)
	if err != nil {
		return err
	}

	if v, ok := dest.(ToGoConverter); ok && (rType&Object == Object || rType&Array == Array) {
		return v.ConvertToGo(retVal.(*ObjectRef))
	} else if rType.isArray() && rType != Object|Array {
		// If return type is an array of convertable java to go types, do the conversion
		converted, err := env.toGoArray(retVal.(*ObjectRef).jobject, rType)
		deleteLocalRef(env.jniEnv, retVal.(*ObjectRef).jobject)
		if err != nil {
			return err
		}

		return assignDest(converted, dest)
	} else {
		return assignDest(retVal, dest)
	}

}

func (o *ObjectRef) genericCallMethod(env *Env, methodName string, rType Type, rClassName string, args ...interface{}) (interface{}, error) {
	class, err := o.getClass(env)
	if err != nil {
		return nil, err
	}

	if err := replaceConvertedArgs(args); err != nil {
		return nil, err
	}
	var methodSig string
	if env.preCalcSig != "" {
		methodSig = env.preCalcSig
		env.preCalcSig = ""
	} else {
		calcSig, err := sigForMethod(rType, rClassName, args)
		if err != nil {
			return nil, err
		}
		methodSig = calcSig
	}

	mid, err := env.callGetMethodID(false, class, methodName, methodSig)
	if err != nil {
		return nil, err
	}

	// create args for jni call
	jniArgs, refs, err := env.createArgs(args)
	if err != nil {
		return nil, err
	}
	defer func() {
		cleanUpArgs(jniArgs)
		for _, ref := range refs {
			deleteLocalRef(env.jniEnv, ref)
		}
	}()

	var retVal interface{}

	switch {
	case rType == Void:
		callVoidMethodA(env.jniEnv, o.jobject, mid, jniArgs)
	case rType == Boolean:
		retVal = toBool(callBooleanMethodA(env.jniEnv, o.jobject, mid, jniArgs))
	case rType == Byte:
		retVal = byte(callByteMethodA(env.jniEnv, o.jobject, mid, jniArgs))
	case rType == Char:
		retVal = uint16(callCharMethodA(env.jniEnv, o.jobject, mid, jniArgs))
	case rType == Short:
		retVal = int16(callShortMethodA(env.jniEnv, o.jobject, mid, jniArgs))
	case rType == Int:
		retVal = int(callIntMethodA(env.jniEnv, o.jobject, mid, jniArgs))
	case rType == Long:
		retVal = int64(callLongMethodA(env.jniEnv, o.jobject, mid, jniArgs))
	case rType == Float:
		retVal = float32(callFloatMethodA(env.jniEnv, o.jobject, mid, jniArgs))
	case rType == Double:
		retVal = float64(callDoubleMethodA(env.jniEnv, o.jobject, mid, jniArgs))
	case rType == Object || rType.isArray():
		obj := callObjectMethodA(env.jniEnv, o.jobject, mid, jniArgs)
		retVal = &ObjectRef{obj, rClassName, rType.isArray()}
	default:
		return nil, errors.New("JNIGI unknown return type")
	}

	if env.exceptionCheck() {
		return nil, env.handleException()
	}

	return retVal, nil
}

// CallNonvirtualMethod calls non virtual method methodName on o with arguments args and stores return value in dest.
func (o *ObjectRef) CallNonvirtualMethod(env *Env, className string, methodName string, dest interface{}, args ...interface{}) error {
	rType, rClassName, err := typeOfReturnValue(dest)
	if err != nil {
		return err
	}

	retVal, err := o.genericCallNonvirtualMethod(env, className, methodName, rType, rClassName, args...)
	if err != nil {
		return err
	}

	if v, ok := dest.(ToGoConverter); ok && (rType&Object == Object || rType&Array == Array) {
		return v.ConvertToGo(retVal.(*ObjectRef))
	} else if rType.isArray() && rType != Object|Array {
		// If return type is an array of convertable java to go types, do the conversion
		converted, err := env.toGoArray(retVal.(*ObjectRef).jobject, rType)
		deleteLocalRef(env.jniEnv, retVal.(*ObjectRef).jobject)
		if err != nil {
			return err
		}

		return assignDest(converted, dest)
	} else {
		return assignDest(retVal, dest)
	}
}

func (o *ObjectRef) genericCallNonvirtualMethod(env *Env, className string, methodName string, rType Type, rClassName string, args ...interface{}) (interface{}, error) {
	class, err := env.callFindClass(className)
	if err != nil {
		return nil, err
	}

	if err := replaceConvertedArgs(args); err != nil {
		return nil, err
	}
	var methodSig string
	if env.preCalcSig != "" {
		methodSig = env.preCalcSig
		env.preCalcSig = ""
	} else {
		calcSig, err := sigForMethod(rType, rClassName, args)
		if err != nil {
			return nil, err
		}
		methodSig = calcSig
	}

	mid, err := env.callGetMethodID(false, class, methodName, methodSig)
	if err != nil {
		return nil, err
	}

	// create args for jni call
	jniArgs, refs, err := env.createArgs(args)
	if err != nil {
		return nil, err
	}
	defer func() {
		cleanUpArgs(jniArgs)
		for _, ref := range refs {
			deleteLocalRef(env.jniEnv, ref)
		}
	}()

	var retVal interface{}

	switch {
	case rType == Void:
		callNonvirtualVoidMethodA(env.jniEnv, o.jobject, class, mid, jniArgs)
	case rType == Boolean:
		retVal = toBool(callNonvirtualBooleanMethodA(env.jniEnv, o.jobject, class, mid, jniArgs))
	case rType == Byte:
		retVal = byte(callNonvirtualByteMethodA(env.jniEnv, o.jobject, class, mid, jniArgs))
	case rType == Char:
		retVal = uint16(callNonvirtualCharMethodA(env.jniEnv, o.jobject, class, mid, jniArgs))
	case rType == Short:
		retVal = int16(callNonvirtualShortMethodA(env.jniEnv, o.jobject, class, mid, jniArgs))
	case rType == Int:
		retVal = int(callNonvirtualIntMethodA(env.jniEnv, o.jobject, class, mid, jniArgs))
	case rType == Long:
		retVal = int64(callNonvirtualLongMethodA(env.jniEnv, o.jobject, class, mid, jniArgs))
	case rType == Float:
		retVal = float32(callNonvirtualFloatMethodA(env.jniEnv, o.jobject, class, mid, jniArgs))
	case rType == Double:
		retVal = float64(callNonvirtualDoubleMethodA(env.jniEnv, o.jobject, class, mid, jniArgs))
	case rType == Object || rType.isArray():
		obj := callNonvirtualObjectMethodA(env.jniEnv, o.jobject, class, mid, jniArgs)
		retVal = &ObjectRef{obj, rClassName, rType.isArray()}
	default:
		return nil, errors.New("JNIGI unknown return type")
	}

	if env.exceptionCheck() {
		return nil, env.handleException()
	}

	return retVal, nil
}

// CallStaticMethod calls static method methodName in class className with arguments args and stores return value in dest.
func (j *Env) CallStaticMethod(className string, methodName string, dest interface{}, args ...interface{}) error {
	rType, rClassName, err := typeOfReturnValue(dest)
	if err != nil {
		return err
	}

	retVal, err := j.genericCallStaticMethod(className, methodName, rType, rClassName, args...)
	if err != nil {
		return err
	}

	if v, ok := dest.(ToGoConverter); ok && (rType&Object == Object || rType&Array == Array) {
		return v.ConvertToGo(retVal.(*ObjectRef))
	} else if rType.isArray() && rType != Object|Array {
		// If return type is an array of convertable java to go types, do the conversion
		converted, err := j.toGoArray(retVal.(*ObjectRef).jobject, rType)
		deleteLocalRef(j.jniEnv, retVal.(*ObjectRef).jobject)
		if err != nil {
			return err
		}

		return assignDest(converted, dest)
	} else {
		return assignDest(retVal, dest)
	}
}

func (j *Env) genericCallStaticMethod(className string, methodName string, rType Type, rClassName string, args ...interface{}) (interface{}, error) {
	class, err := j.callFindClass(className)
	if err != nil {
		return nil, err
	}

	if err := replaceConvertedArgs(args); err != nil {
		return nil, err
	}
	var methodSig string
	if j.preCalcSig != "" {
		methodSig = j.preCalcSig
		j.preCalcSig = ""
	} else {
		calcSig, err := sigForMethod(rType, rClassName, args)
		if err != nil {
			return nil, err
		}
		methodSig = calcSig
	}

	mid, err := j.callGetMethodID(true, class, methodName, methodSig)
	if err != nil {
		return nil, err
	}

	// create args for jni call
	jniArgs, refs, err := j.createArgs(args)
	if err != nil {
		return nil, err
	}
	defer func() {
		cleanUpArgs(jniArgs)
		for _, ref := range refs {
			deleteLocalRef(j.jniEnv, ref)
		}
	}()

	var retVal interface{}

	switch {
	case rType == Void:
		callStaticVoidMethodA(j.jniEnv, class, mid, jniArgs)
	case rType == Boolean:
		retVal = toBool(callStaticBooleanMethodA(j.jniEnv, class, mid, jniArgs))
	case rType == Byte:
		retVal = byte(callStaticByteMethodA(j.jniEnv, class, mid, jniArgs))
	case rType == Char:
		retVal = uint16(callStaticCharMethodA(j.jniEnv, class, mid, jniArgs))
	case rType == Short:
		retVal = int16(callStaticShortMethodA(j.jniEnv, class, mid, jniArgs))
	case rType == Int:
		retVal = int(callStaticIntMethodA(j.jniEnv, class, mid, jniArgs))
	case rType == Long:
		retVal = int64(callStaticLongMethodA(j.jniEnv, class, mid, jniArgs))
	case rType == Float:
		retVal = float32(callStaticFloatMethodA(j.jniEnv, class, mid, jniArgs))
	case rType == Double:
		retVal = float64(callStaticDoubleMethodA(j.jniEnv, class, mid, jniArgs))
	case rType == Object || rType.isArray():
		obj := callStaticObjectMethodA(j.jniEnv, class, mid, jniArgs)
		retVal = &ObjectRef{obj, rClassName, rType.isArray()}
	default:
		return nil, errors.New("JNIGI unknown return type")
	}

	if j.exceptionCheck() {
		return nil, j.handleException()
	}

	return retVal, nil
}

func (j *Env) callGetFieldID(static bool, class jclass, name, sig string) (jfieldID, error) {
	fnCstr := cString(name)
	defer free(fnCstr)

	sigCstr := cString(sig)
	defer free(sigCstr)

	var fid jfieldID
	if static {
		fid = getStaticFieldID(j.jniEnv, class, fnCstr, sigCstr)
	} else {
		fid = getFieldID(j.jniEnv, class, fnCstr, sigCstr)
	}
	if fid == 0 {
		return 0, j.handleException()
	}

	return fid, nil
}

// GetField gets field fieldName in o and stores value in dest.
func (o *ObjectRef) GetField(env *Env, fieldName string, dest interface{}) error {
	fType, fClassName, err := typeOfReturnValue(dest)
	if err != nil {
		return err
	}

	fieldVal, err := o.genericGetField(env, fieldName, fType, fClassName)
	if err != nil {
		return err
	}

	if v, ok := dest.(ToGoConverter); ok && (fType&Object == Object || fType&Array == Array) {
		return v.ConvertToGo(fieldVal.(*ObjectRef))
	} else if fType.isArray() && fType != Object|Array {
		// If return type is an array of convertable java to go types, do the conversion
		converted, err := env.toGoArray(fieldVal.(*ObjectRef).jobject, fType)
		deleteLocalRef(env.jniEnv, fieldVal.(*ObjectRef).jobject)
		if err != nil {
			return err
		}

		return assignDest(converted, dest)
	} else {
		return assignDest(fieldVal, dest)
	}
}

func (o *ObjectRef) genericGetField(env *Env, fieldName string, fType Type, fClassName string) (interface{}, error) {
	class, err := o.getClass(env)
	if err != nil {
		return nil, err
	}

	var fieldSig string
	if env.preCalcSig != "" {
		fieldSig = env.preCalcSig
		env.preCalcSig = ""
	} else {
		fieldSig = typeSignature(fType, fClassName)
	}

	fid, err := env.callGetFieldID(false, class, fieldName, fieldSig)
	if err != nil {
		return nil, err
	}

	var retVal interface{}

	switch {
	case fType == Boolean:
		retVal = toBool(getBooleanField(env.jniEnv, o.jobject, fid))
	case fType == Byte:
		retVal = byte(getByteField(env.jniEnv, o.jobject, fid))
	case fType == Char:
		retVal = uint16(getCharField(env.jniEnv, o.jobject, fid))
	case fType == Short:
		retVal = int16(getShortField(env.jniEnv, o.jobject, fid))
	case fType == Int:
		retVal = int(getIntField(env.jniEnv, o.jobject, fid))
	case fType == Long:
		retVal = int64(getLongField(env.jniEnv, o.jobject, fid))
	case fType == Float:
		retVal = float32(getFloatField(env.jniEnv, o.jobject, fid))
	case fType == Double:
		retVal = float64(getDoubleField(env.jniEnv, o.jobject, fid))
	case fType == Object || fType.isArray():
		obj := getObjectField(env.jniEnv, o.jobject, fid)
		retVal = &ObjectRef{obj, fClassName, fType.isArray()}
	default:
		return nil, errors.New("JNIGI unknown field type")
	}

	if env.exceptionCheck() {
		return nil, env.handleException()
	}

	return retVal, nil
}

// SetField sets field fieldName in o to value.
func (o *ObjectRef) SetField(env *Env, fieldName string, value interface{}) error {
	class, err := o.getClass(env)
	if err != nil {
		return err
	}

	vType, vClassName, err := typeOfValue(value)
	if err != nil {
		return err
	}

	var fieldSig string
	if env.preCalcSig != "" {
		fieldSig = env.preCalcSig
		env.preCalcSig = ""
	} else {
		fieldSig = typeSignature(vType, vClassName)
	}

	fid, err := env.callGetFieldID(false, class, fieldName, fieldSig)
	if err != nil {
		return err
	}

	switch v := value.(type) {
	case bool:
		setBooleanField(env.jniEnv, o.jobject, fid, fromBool(v))
	case byte:
		setByteField(env.jniEnv, o.jobject, fid, jbyte(v))
	case uint16:
		setCharField(env.jniEnv, o.jobject, fid, jchar(v))
	case int16:
		setShortField(env.jniEnv, o.jobject, fid, jshort(v))
	case int32:
		setIntField(env.jniEnv, o.jobject, fid, jint(v))
	case int:
		setIntField(env.jniEnv, o.jobject, fid, jint(int32(v)))
	case int64:
		setLongField(env.jniEnv, o.jobject, fid, jlong(v))
	case float32:
		setFloatField(env.jniEnv, o.jobject, fid, jfloat(v))
	case float64:
		setDoubleField(env.jniEnv, o.jobject, fid, jdouble(v))
	case jobj:
		setObjectField(env.jniEnv, o.jobject, fid, v.jobj())
	case []bool, []byte, []int16, []uint16, []int32, []int, []int64, []float32, []float64:
		array, err := env.toJavaArray(v)
		if err != nil {
			return err
		}
		defer deleteLocalRef(env.jniEnv, array)
		setObjectField(env.jniEnv, o.jobject, fid, jobject(array))
	default:
		return errors.New("JNIGI unknown field value")
	}

	if env.exceptionCheck() {
		return env.handleException()
	}

	return nil
}

// GetField gets field fieldName in class className, stores value in dest.
func (j *Env) GetStaticField(className string, fieldName string, dest interface{}) error {
	fType, fClassName, err := typeOfReturnValue(dest)
	if err != nil {
		return err
	}

	fieldVal, err := j.genericGetStaticField(className, fieldName, fType, fClassName)
	if err != nil {
		return err
	}

	if v, ok := dest.(ToGoConverter); ok && (fType&Object == Object || fType&Array == Array) {
		return v.ConvertToGo(fieldVal.(*ObjectRef))
	} else if fType.isArray() && fType != Object|Array {
		// If return type is an array of convertable java to go types, do the conversion
		converted, err := j.toGoArray(fieldVal.(*ObjectRef).jobject, fType)
		deleteLocalRef(j.jniEnv, fieldVal.(*ObjectRef).jobject)
		if err != nil {
			return err
		}

		return assignDest(converted, dest)
	} else {
		return assignDest(fieldVal, dest)
	}
}

func (j *Env) genericGetStaticField(className string, fieldName string, fType Type, fClassName string) (interface{}, error) {
	class, err := j.callFindClass(className)
	if err != nil {
		return nil, err
	}

	var fieldSig string
	if j.preCalcSig != "" {
		fieldSig = j.preCalcSig
		j.preCalcSig = ""
	} else {
		fieldSig = typeSignature(fType, fClassName)
	}

	fid, err := j.callGetFieldID(true, class, fieldName, fieldSig)
	if err != nil {
		return nil, err
	}

	var retVal interface{}

	switch {
	case fType == Boolean:
		retVal = toBool(getStaticBooleanField(j.jniEnv, class, fid))
	case fType == Byte:
		retVal = byte(getStaticByteField(j.jniEnv, class, fid))
	case fType == Char:
		retVal = uint16(getStaticCharField(j.jniEnv, class, fid))
	case fType == Short:
		retVal = int16(getStaticShortField(j.jniEnv, class, fid))
	case fType == Int:
		retVal = int(getStaticIntField(j.jniEnv, class, fid))
	case fType == Long:
		retVal = int64(getStaticLongField(j.jniEnv, class, fid))
	case fType == Float:
		retVal = float32(getStaticFloatField(j.jniEnv, class, fid))
	case fType == Double:
		retVal = float64(getStaticDoubleField(j.jniEnv, class, fid))
	case fType == Object || fType.isArray():
		obj := getStaticObjectField(j.jniEnv, class, fid)
		retVal = &ObjectRef{obj, fClassName, fType.isArray()}
	default:
		return nil, errors.New("JNIGI unknown field type")
	}

	if j.exceptionCheck() {
		return nil, j.handleException()
	}

	return retVal, nil
}

// SetField sets field fieldName in class className to value.
func (j *Env) SetStaticField(className string, fieldName string, value interface{}) error {
	class, err := j.callFindClass(className)
	if err != nil {
		return err
	}

	vType, vClassName, err := typeOfValue(value)
	if err != nil {
		return err
	}

	var fieldSig string
	if j.preCalcSig != "" {
		fieldSig = j.preCalcSig
		j.preCalcSig = ""
	} else {
		fieldSig = typeSignature(vType, vClassName)
	}

	fid, err := j.callGetFieldID(true, class, fieldName, fieldSig)
	if err != nil {
		return err
	}

	switch v := value.(type) {
	case bool:
		setStaticBooleanField(j.jniEnv, class, fid, fromBool(v))
	case byte:
		setStaticByteField(j.jniEnv, class, fid, jbyte(v))
	case uint16:
		setStaticCharField(j.jniEnv, class, fid, jchar(v))
	case int16:
		setStaticShortField(j.jniEnv, class, fid, jshort(v))
	case int32:
		setStaticIntField(j.jniEnv, class, fid, jint(v))
	case int:
		setStaticIntField(j.jniEnv, class, fid, jint(int32(v)))
	case int64:
		setStaticLongField(j.jniEnv, class, fid, jlong(v))
	case float32:
		setStaticFloatField(j.jniEnv, class, fid, jfloat(v))
	case float64:
		setStaticDoubleField(j.jniEnv, class, fid, jdouble(v))
	case jobj:
		setStaticObjectField(j.jniEnv, class, fid, v.jobj())
	case []bool, []byte, []int16, []uint16, []int32, []int, []int64, []float32, []float64:
		array, err := j.toJavaArray(v)
		if err != nil {
			return err
		}
		defer deleteLocalRef(j.jniEnv, array)
		setStaticObjectField(j.jniEnv, class, fid, jobject(array))
	default:
		return errors.New("JNIGI unknown field value")
	}

	if j.exceptionCheck() {
		return j.handleException()
	}

	return nil
}

// RegisterNative calls JNI RegisterNative for class className, method methodName with return type returnType and parameters params,
// fptr is used as native function.
func (j *Env) RegisterNative(className, methodName string, returnType TypeSpec, params []interface{}, fptr interface{}) error {
	class, err := j.callFindClass(className)
	if err != nil {
		return err
	}

	mnCstr := cString(methodName)
	defer free(mnCstr)
	rType, rClassName, err := typeOfReturnValue(returnType)
	if err != nil {
		return err
	}

	// Convert strings in params to ObjectType. This to retain compaibility.
	// with code that assumes strings signify ObjectType.
	compatPrm := make([]interface{}, len(params))
	for i, param := range params {
		if v, ok := param.(string); ok {
			compatPrm[i] = ObjectType(v)
		} else {
			compatPrm[i] = param
		}
	}

	sig, err := sigForMethod(rType, rClassName, compatPrm)
	if err != nil {
		return err
	}
	sigCstr := cString(sig)
	defer free(sigCstr)

	if registerNative(j.jniEnv, class, mnCstr, sigCstr, fptr.(unsafe.Pointer)) < 0 {
		return j.handleException()
	}

	return nil
}

// NewGlobalRef creates a new object reference to o in Env j.
func (j *Env) NewGlobalRef(o *ObjectRef) *ObjectRef {
	g := newGlobalRef(j.jniEnv, o.jobject)
	return &ObjectRef{g, o.className, o.isArray}
}

// DeleteGlobalRef deletes global object reference o.
func (j *Env) DeleteGlobalRef(o *ObjectRef) {
	deleteGlobalRef(j.jniEnv, o.jobject)
	o.jobject = 0
}

// DeleteLocalRef deletes object reference o in Env j.
func (j *Env) DeleteLocalRef(o *ObjectRef) {
	deleteLocalRef(j.jniEnv, o.jobject)
	o.jobject = 0
}

// EnsureLocalCapacity calls JNI EnsureLocalCapacity on Env j
func (j *Env) EnsureLocalCapacity(capacity int32) error {
	success := ensureLocalCapacity(j.jniEnv, jint(capacity)) == 0
	if j.exceptionCheck() {
		return j.handleException()
	}
	if !success {
		return errors.New("JNIGI: ensureLocalCapacity error")
	}
	return nil
}

// PushLocalFrame calls JNI PushLocalFrame on Env j
func (j *Env) PushLocalFrame(capacity int32) error {
	success := pushLocalFrame(j.jniEnv, jint(capacity)) == 0
	if j.exceptionCheck() {
		return j.handleException()
	}
	if !success {
		return errors.New("JNIGI: pushLocalFrame error")
	}
	return nil
}

// PopLocalFrame calls JNI popLocalFrame on Env j
func (j *Env) PopLocalFrame(result *ObjectRef) *ObjectRef {
	if result == nil {
		result = &ObjectRef{}
	}
	o := popLocalFrame(j.jniEnv, result.jobject)
	result.jobject = 0
	return &ObjectRef{o, result.className, result.isArray}
}

var utf8 *ObjectRef

// GetUTF8String return global reference to java/lang/String containing "UTF-8"
func (j *Env) GetUTF8String() *ObjectRef {
	if utf8 == nil {
		str, err := j.NewObject("java/lang/String", []byte("UTF-8"))
		if err != nil {
			panic(err)
		}
		global := j.NewGlobalRef(str)
		j.DeleteLocalRef(str)
		utf8 = global
	}

	return utf8
}

// StackTraceElement is a struct holding the contents of java.lang.StackTraceElement
// for use in a ThrowableError.
type StackTraceElement struct {
	ClassName      string
	FileName       string
	LineNumber     int
	MethodName     string
	IsNativeMethod bool
	AsString       string
}

func (el StackTraceElement) String() string {
	return el.AsString
}

// ThrowableError is an error struct that holds the relevant contents of a
// java.lang.Throwable. This is the returned error from ThrowableErrorExceptionHandler.
type ThrowableError struct {
	ClassName        string
	LocalizedMessage string
	Message          string
	StackTrace       []StackTraceElement
	AsString         string
	Cause            *ThrowableError
}

func (e ThrowableError) String() string {
	return e.AsString
}

func (e ThrowableError) Error() string {
	return e.AsString
}

func stringFromJavaLangString(env *Env, ref *ObjectRef) string {
	if ref.IsNil() {
		return ""
	}
	env.PrecalculateSignature("(Ljava/lang/String;)[B")
	var ret []byte
	err := ref.CallMethod(env, "getBytes", &ret, env.GetUTF8String())
	if err != nil {
		return ""
	}
	return string(ret)
}

func callStringMethodAndAssign(env *Env, obj *ObjectRef, method string, assign func(s string)) error {
	env.PrecalculateSignature("()Ljava/lang/String;")
	strref := NewObjectRef("java/lang/String")
	err := obj.CallMethod(env, method, strref)
	if err != nil {
		return err
	}
	defer env.DeleteLocalRef(strref)

	assign(stringFromJavaLangString(env, strref))

	return nil
}

// NewStackTraceElementFromObject creates a new StackTraceElement with its contents
// set from the values provided in stackTraceElement's methods.
func NewStackTraceElementFromObject(env *Env, stackTraceElement *ObjectRef) (*StackTraceElement, error) {

	if stackTraceElement.IsNil() {
		return nil, nil
	}

	getStringAndAssign := func(method string, assign func(s string)) error {
		return callStringMethodAndAssign(env, stackTraceElement, method, assign)
	}

	out := StackTraceElement{}

	// ClassName
	if err := getStringAndAssign("getClassName", func(s string) {
		out.ClassName = s
	}); err != nil {
		return nil, err
	}

	// FileName
	if err := getStringAndAssign("getFileName", func(s string) {
		out.FileName = s
	}); err != nil {
		return nil, err
	}

	// MethodName
	if err := getStringAndAssign("getMethodName", func(s string) {
		out.MethodName = s
	}); err != nil {
		return nil, err
	}

	// ToString
	if err := getStringAndAssign("toString", func(s string) {
		out.AsString = s
	}); err != nil {
		return nil, err
	}

	// LineNumber
	{
		env.PrecalculateSignature("()I")
		var lineNum int
		err := stackTraceElement.CallMethod(env, "getLineNumber", &lineNum)
		if err != nil {
			return nil, err
		}
		out.LineNumber = lineNum
	}

	// IsNativeMethod
	{
		env.PrecalculateSignature("()Z")
		var isNative bool
		err := stackTraceElement.CallMethod(env, "isNativeMethod", &isNative)
		if err != nil {
			return nil, err
		}
		out.IsNativeMethod = isNative
	}

	return &out, nil
}

// NewThrowableErrorFromObject creates a new ThrowableError with its contents
// set from the values provided in throwable's methods.
func NewThrowableErrorFromObject(env *Env, throwable *ObjectRef) (*ThrowableError, error) {

	if throwable.IsNil() {
		return nil, nil
	}

	getStringAndAssign := func(obj *ObjectRef, method string, assign func(s string)) error {
		return callStringMethodAndAssign(env, obj, method, assign)
	}

	out := &ThrowableError{}

	// ClassName
	{
		objClass := getObjectClass(env.jniEnv, throwable.jobject)
		if objClass == 0 {
			return nil, fmt.Errorf("unable to get throwable class")
		}

		clsref := WrapJObject(uintptr(objClass), "java/lang/Class", false)
		defer env.DeleteLocalRef(clsref)

		if err := getStringAndAssign(clsref, "getName", func(s string) {
			out.ClassName = s
		}); err != nil {
			return nil, err
		}
	}

	// AsString
	if err := getStringAndAssign(throwable, "toString", func(s string) {
		out.AsString = s
	}); err != nil {
		return nil, err
	}

	// From this point on, return throwableError if a call fails, since we have some basic information.

	// LocalizedMessage
	if err := getStringAndAssign(throwable, "getLocalizedMessage", func(s string) {
		out.LocalizedMessage = s
	}); err != nil {
		return out, err
	}

	// Message
	if err := getStringAndAssign(throwable, "getMessage", func(s string) {
		out.Message = s
	}); err != nil {
		return out, err
	}

	// StackTrace
	{
		env.PrecalculateSignature("()[Ljava/lang/StackTraceElement;")
		stkTrcArr := NewObjectRef("java/lang/StackTraceElement")
		err := throwable.CallMethod(env, "getStackTrace", stkTrcArr)
		if err != nil {
			return out, err
		}
		defer env.DeleteLocalRef(stkTrcArr)

		if !stkTrcArr.IsNil() {
			stkTrcSlc := env.FromObjectArray(stkTrcArr)
			stackTrace := make([]StackTraceElement, 0, len(stkTrcSlc))
			for _, stkTrc := range stkTrcSlc {
				if stkTrc.IsNil() {
					continue
				}
				defer env.DeleteLocalRef(stkTrc)
				stackTraceElement, err := NewStackTraceElementFromObject(env, stkTrc)
				if err != nil || stackTraceElement == nil {
					continue
				}
				stackTrace = append(stackTrace, *stackTraceElement)
			}

			out.StackTrace = stackTrace
		}
	}

	// Cause
	{
		env.PrecalculateSignature("()Ljava/lang/Throwable;")
		obj := NewObjectRef("java/lang/Throwable")
		err := throwable.CallMethod(env, "getCause", obj)
		if err != nil {
			return out, err
		}
		defer env.DeleteLocalRef(obj)

		out.Cause, _ = NewThrowableErrorFromObject(env, obj)
	}

	return out, nil
}

var (
	errThrowableConvertFail = fmt.Errorf("Java exception occured")

	// DefaultExceptionHandler is an alias for DescribeExceptionHandler, which is the default.
	DefaultExceptionHandler = DescribeExceptionHandler

	// DescribeExceptionHandler calls the JNI exceptionDescribe function.
	DescribeExceptionHandler ExceptionHandler = ExceptionHandlerFunc(func(env *Env, exception *ObjectRef) error {
		exceptionDescribe(env.jniEnv)
		exceptionClear(env.jniEnv)
		return errors.New("Java exception occured. check stderr")
	})

	// ThrowableToStringExceptionHandler calls ToString on the exception and returns an error
	// with the returned value as its Error message.
	// If exception is nil or the toString() call fails, a generic default error is returned.
	ThrowableToStringExceptionHandler ExceptionHandler = ExceptionHandlerFunc(func(env *Env, exception *ObjectRef) error {
		exceptionClear(env.jniEnv)
		if exception.IsNil() {
			return errThrowableConvertFail
		}
		msg := "Java exception occured"
		callStringMethodAndAssign(env, exception, "toString", func(s string) {
			if s == "" {
				return
			}
			msg = s
		})
		return errors.New(msg)
	})

	// ThrowableErrorExceptionHandler populates a new ThrowableError with the values of exception.
	// If exception is nil, the getClass().getName(), or the toString call fails, a generic default
	// error is returned.
	ThrowableErrorExceptionHandler ExceptionHandler = ExceptionHandlerFunc(func(env *Env, exception *ObjectRef) error {
		exceptionClear(env.jniEnv)
		if exception.IsNil() {
			return errThrowableConvertFail
		}
		throwableError, _ := NewThrowableErrorFromObject(env, exception)
		if throwableError == nil {
			return errThrowableConvertFail
		}
		return *throwableError
	})
)
