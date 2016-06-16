// Copyright 2016 Tim O'Brien. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jnigi

import (
	"errors"
	"fmt"
	"runtime"
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

type ObjectRef struct {
	jobject   jobject
	className string
	isArray   bool
}

func WrapJObject(jobj uintptr, className string, isArray bool) *ObjectRef {
	return &ObjectRef{jobject(jobj), className, isArray}
}

func (o *ObjectRef) Cast(className string) *ObjectRef {
	if className == o.className {
		return o
	} else {
		return &ObjectRef{o.jobject, className, o.isArray}
	}
}

func (o *ObjectRef) IsNil() bool {
	return o.jobject == 0
}

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

type jobj interface {
	jobj() jobject
}

var classCache map[string]jclass = make(map[string]jclass)

type Env struct {
	jniEnv     unsafe.Pointer
	preCalcSig string
	noReturnConvert bool
}

type JVM struct {
	javaVM unsafe.Pointer
}

type JVMInitArgs struct {
	javaVMInitArgs unsafe.Pointer
}

func CreateJVM(jvmInitArgs *JVMInitArgs) (*JVM, *Env, error) {
	runtime.LockOSThread()

	p := malloc(unsafe.Sizeof((unsafe.Pointer)(nil)))
	p2 := malloc(unsafe.Sizeof((unsafe.Pointer)(nil)))

	if jni_CreateJavaVM(p2, p, jvmInitArgs.javaVMInitArgs) < 0 {
		return nil, nil, errors.New("Couldn't instantiate JVM")
	}
	jvm := &JVM{*(*unsafe.Pointer)(p2)}
	env := &Env{jniEnv: *(*unsafe.Pointer)(p)}

	free(p)
	free(p2)
	return jvm, env, nil
}

func (j *JVM) AttachCurrentThread() *Env {
	runtime.LockOSThread()
	p := malloc(unsafe.Sizeof((unsafe.Pointer)(nil)))

	//	p := (**C.JNIEnv)(malloc(unsafe.Sizeof((*C.JNIEnv)(nil))))

	if attachCurrentThread(j.javaVM, p, nil) < 0 {
		panic("AttachCurrentThread failed")
	}

	return &Env{jniEnv: *(*unsafe.Pointer)(p)}
}

func (j *JVM) DetachCurrentThread() error {
	if detachCurrentThread(j.javaVM) < 0 {
		return errors.New("JNIGI: detachCurrentThread error")
	}
	return nil
}

func (j *Env) exceptionCheck() bool {
	return toBool(exceptionCheck(j.jniEnv))
}

func (j *Env) describeException() {
	exceptionDescribe(j.jniEnv)
}

func (j *Env) handleException() error {
	var eStr string
	if e := exceptionOccurred(j.jniEnv); e == 0 {
		eStr = "Java JNI function returned error but JNI indicates no current exception"
	} else {
		//TODO: return exception string here instead of just printing to stderr with "exceptionDescribe"
		eStr = "Java exception occured. check stderr"
		exceptionDescribe(j.jniEnv)
		exceptionClear(j.jniEnv)
		defer deleteLocalRef(j.jniEnv, jobject(e))
	}
	return errors.New(eStr)
}

func (j *Env) NewObject(className string, args ...interface{}) (*ObjectRef, error) {
	class, err := j.callFindClass(className)
	if err != nil {
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
	if v, ok := classCache[className]; ok {
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

func (j *Env) PrecalculateSignature(sig string) {
	j.preCalcSig = sig
}

func (j *Env) NoReturnConvert() {
	j.noReturnConvert = true
}

const big = 1024 * 1024 * 100

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

type ByteArray struct {
	arr jbyteArray
	n int
}

func (j *Env) NewByteArray(n int) *ByteArray {
	a := newByteArray(j.jniEnv, jsize(n))
	return &ByteArray{a, n}
}

func (j *Env) NewByteArrayFromSlice(src []byte) *ByteArray {
	b := j.NewByteArray(len(src))
	if len(src) > 0 {
		bytes := b.GetCritical(j)
		copy(bytes, src)
		b.ReleaseCritical(j, bytes)
	}
	return b
}

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

func (b *ByteArray) GetCritical(env *Env) []byte {
	if b.n == 0 {
		return nil
	}
	ptr := getPrimitiveArrayCritical(env.jniEnv, jarray(b.arr), nil)
	return (*(*[big]byte)(ptr))[0:b.n]
}

func (b *ByteArray) ReleaseCritical(env *Env, bytes []byte) {
	if len(bytes) == 0 {
		return
	}
	ptr := unsafe.Pointer(&bytes[0])
	releasePrimitiveArrayCritical(env.jniEnv, jarray(b.arr), ptr, 0)
}

//returns jlo
func (b *ByteArray) GetObject() *ObjectRef {
	return &ObjectRef{jobject(b.arr), "java/lang/Object", false}
}

func (b *ByteArray) SetObject(o *ObjectRef) {
	b.arr = jbyteArray(o.jobject)
}

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
			err = fmt.Errorf("JNIGI: argument not a valid value %t (%v)", args[i], args[i])
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

type ObjectType string

type ObjectArrayType string

type convertedArray interface {
	getType() Type
}

func typeOfValue(value interface{}) (t Type, className string, err error) {
	switch v := value.(type) {
	case Type:
		t = v
		if t.baseType() == Object {
			className = "java/lang/Object"
		}
	case string:
		t = Object
		className = v
	case ObjectType:
		t = Object
		className = string(v)
	case ObjectArrayType:
		t = Object | Array
		className = string(v)
	case *ObjectRef:
		t = Object
		if v.isArray {
			t = t | Array
		}
		className = v.className
	case bool:
		t = Boolean
	case byte:
		t = Byte
	case int16:
		t = Short
	case uint16:
		t = Char
	case int32:
		t = Int
	case int:
		t = Int
	case int64:
		t = Long
	case float32:
		t = Float
	case float64:
		t = Double
	case []bool:
		t = Boolean | Array
		className = "java/lang/Object"
	case []byte:
		t = Byte | Array
		className = "java/lang/Object"
	case []uint16:
		t = Char | Array
		className = "java/lang/Object"
	case []int16:
		t = Short | Array
		className = "java/lang/Object"
	case []int32:
		t = Int | Array
		className = "java/lang/Object"
	case []int:
		t = Int | Array
		className = "java/lang/Object"
	case []int64:
		t = Long | Array
		className = "java/lang/Object"
	case []float32:
		t = Float | Array
		className = "java/lang/Object"
	case []float64:
		t = Double | Array
		className = "java/lang/Object"
	case convertedArray:
		t = v.getType()
		className = "java/lang/Object"
	default:
		err = fmt.Errorf("JNIGI: unknown type %t (%v)", v, v)
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

func (o *ObjectRef) CallMethod(env *Env, methodName string, returnType interface{}, args ...interface{}) (interface{}, error) {
	class, err := env.callFindClass(o.className)
	if err != nil {
		return nil, err
	}

	rType, rClassName, err := typeOfValue(returnType)
	if err != nil {
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

	var arrayToConvert jobject
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
		if rType == Object || rType == Object|Array || env.noReturnConvert {
			retVal = &ObjectRef{obj, rClassName, rType.isArray()}
		} else {
			arrayToConvert = obj
		}
	default:
		return nil, errors.New("JNIGI unknown return type")
	}

	env.noReturnConvert = false

	if env.exceptionCheck() {
		return nil, env.handleException()
	}

	if arrayToConvert != 0 {
		retVal, err = env.toGoArray(arrayToConvert, rType)
		if err != nil {
			return nil, err
		}
	}

	return retVal, nil
}

func (j *Env) CallStaticMethod(className string, methodName string, returnType interface{}, args ...interface{}) (interface{}, error) {
	class, err := j.callFindClass(className)
	if err != nil {
		return nil, err
	}

	rType, rClassName, err := typeOfValue(returnType)
	if err != nil {
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

	var arrayToConvert jobject
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
		if rType == Object || rType == Object|Array || j.noReturnConvert {
			retVal = &ObjectRef{obj, rClassName, rType.isArray()}
		} else {
			arrayToConvert = obj
		}
	default:
		return nil, errors.New("JNIGI unknown return type")
	}

	j.noReturnConvert = false

	if j.exceptionCheck() {
		return nil, j.handleException()
	}

	if arrayToConvert != 0 {
		retVal, err = j.toGoArray(arrayToConvert, rType)
		if err != nil {
			return nil, err
		}
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

func (o *ObjectRef) GetField(env *Env, fieldName string, fieldType interface{}) (interface{}, error) {
	class := getObjectClass(env.jniEnv, o.jobject)
	if class == 0 {
		return nil, env.handleException()
	}
	defer deleteLocalRef(env.jniEnv, jobject(class))

	fType, fClassName, err := typeOfValue(fieldType)
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

	var arrayToConvert jobject
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
		if fType == Object || fType == Object|Array || env.noReturnConvert {
			retVal = &ObjectRef{obj, fClassName, fType.isArray()}
		} else {
			arrayToConvert = obj
		}
	default:
		return nil, errors.New("JNIGI unknown field type")
	}

	env.noReturnConvert = false

	if env.exceptionCheck() {
		return nil, env.handleException()
	}

	if arrayToConvert != 0 {
		retVal, err = env.toGoArray(arrayToConvert, fType)
		if err != nil {
			return nil, err
		}
	}

	return retVal, nil
}

func (o *ObjectRef) SetField(env *Env, fieldName string, value interface{}) error {
	class := getObjectClass(env.jniEnv, o.jobject)
	if class == 0 {
		return env.handleException()
	}
	defer deleteLocalRef(env.jniEnv, jobject(class))

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

func (j *Env) GetStaticField(className string, fieldName string, fieldType interface{}) (interface{}, error) {
	class, err := j.callFindClass(className)
	if err != nil {
		return nil, err
	}

	fType, fClassName, err := typeOfValue(fieldType)
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

	var arrayToConvert jobject
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
		if fType == Object || fType == Object|Array || j.noReturnConvert {
			retVal = &ObjectRef{obj, fClassName, fType.isArray()}
		} else {
			arrayToConvert = obj
		}
	default:
		return nil, errors.New("JNIGI unknown field type")
	}

	j.noReturnConvert = false

	if j.exceptionCheck() {
		return nil, j.handleException()
	}

	if arrayToConvert != 0 {
		retVal, err = j.toGoArray(arrayToConvert, fType)
		if err != nil {
			return nil, err
		}
	}

	return retVal, nil
}

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

func (j *Env) RegisterNative(className, methodName string, returnType interface{}, params []interface{}, fptr interface{}) error {
	class, err := j.callFindClass(className)
	if err != nil {
		return err
	}

	mnCstr := cString(methodName)
	defer free(mnCstr)
	rType, rClassName, err := typeOfValue(returnType)
	if err != nil {
		return err
	}
	sig, err := sigForMethod(rType, rClassName, params)
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

func (j *Env) NewGlobalRef(o *ObjectRef) *ObjectRef {
	g := newGlobalRef(j.jniEnv, o.jobject)
	return &ObjectRef{g, o.className, o.isArray}
}

func (j *Env) DeleteGlobalRef(o *ObjectRef) {
	deleteGlobalRef(j.jniEnv, o.jobject)
	o.jobject = 0
}

func (j *Env) DeleteLocalRef(o *ObjectRef) {
	deleteLocalRef(j.jniEnv, o.jobject)
	o.jobject = 0
}

var utf8 *ObjectRef

// return global reference to java/lang/String containing "UTF-8"
func (j *Env) GetUTF8String() *ObjectRef {
	if utf8 == nil {
		cStr := cString("UTF-8")
		local := newStringUTF(j.jniEnv, cStr)
		if local == 0 {
			panic(j.handleException())
		}
		global := jstring(newGlobalRef(j.jniEnv, jobject(local)))
		deleteLocalRef(j.jniEnv, jobject(local))
		free(cStr)
		utf8 = &ObjectRef{jobject: jobject(global), isArray: false, className: "java/lang/String"}
	}

	return utf8
}
