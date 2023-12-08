// Copyright 2016 Tim O'Brien. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jnigi

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var env *Env
var jvm *JVM

// Run them all here so we can be sure they run on the same Goroutine
func TestAll(t *testing.T) {
	PTestInit(t)
	PTestBasic(t)
	PTestObjectArrays(t)
	PTestConvert(t)
	PTestInstanceOf(t)
	PTestByteArray(t)
	PTestAttach(t)
	PTestGetJVM(t)
	PTestEnsureLocalCapacity(t)
	PTestPushPopLocalFrame(t)
	PTestHandleException(t)
	PTestCast(t)
	PTestRegisterNative(t)
	PTestDestroy(t)
}

func PTestInit(t *testing.T) {
	libPath := AttemptToFindJVMLibPath()
	if err := LoadJVMLib(libPath); err != nil {
		t.Logf("library path = %s", libPath)
		t.Log("can use JAVA_HOME environment variable to set JRE root directory")
		t.Fatal(err)
	}
	runtime.LockOSThread()
	cwd, _ := os.Getwd()
	jvm2, e2, err := CreateJVM(NewJVMInitArgs(false, true, DEFAULT_VERSION, []string{"-Xcheck:jni", "-Djava.class.path=" + filepath.Join(cwd, "java/test/out")}))
	if err != nil {
		t.Fatal(err)
	}
	env = e2
	jvm = jvm2

	t.Logf("%x", e2.jniEnv)
}

func PTestBasic(t *testing.T) {
	// new object, int method
	obj, err := env.NewObject("java/lang/Object")
	if err != nil {
		t.Fatal(err)
	}
	var v int
	if err := obj.CallMethod(env, "hashCode", &v); err != nil {
		t.Fatal(err)
	}
	env.DeleteLocalRef(obj)

	// byte array argument, byte array method
	var testStr string = "hello world"
	str, err := env.NewObject("java/lang/String", []byte(testStr))
	if err != nil {
		t.Fatal(err)
	}
	defer env.DeleteLocalRef(str)
	var goBytes []byte
	if err := str.CallMethod(env, "getBytes", &goBytes, env.GetUTF8String()); err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, testStr, string(goBytes)) {
		t.Fail()
	}

	// object method, int arg, object arg
	str2 := NewObjectRef("java/lang/String")
	if err := str.CallMethod(env, "substring", str2, 6); err != nil {
		t.Fatal(err)
	}
	defer env.DeleteLocalRef(str2)

	var dummy []byte
	if err := str2.CallMethod(env, "getBytes", &dummy); err != nil {
		t.Fatal(err)
	}
	// test precalculated signature
	env.PrecalculateSignature("(Ljava/lang/String;)Z")
	var same bool
	if err := str.CallMethod(env, "endsWith", &same, str2); err != nil {
		t.Fatal(err)
	}
	if !same {
		t.Errorf("basic test failed")
	}

	// call static method
	jvmVer := NewObjectRef("java/lang/String")
	if err := env.CallStaticMethod("java/lang/System", "getProperty", jvmVer, fromGoStr(t, "java.vm.version")); err != nil {
		t.Fatal(err)
	}

	// get static field
	err = env.SetStaticField("java/util/Calendar", "APRIL", 5)
	if err != nil {
		t.Fatal(err)
	}
	var calPos int
	err = env.GetStaticField("java/util/Calendar", "APRIL", &calPos)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 5, calPos) {
		t.Fail()
	}

	// set/get object field
	pt, err := env.NewObject("java/awt/Point")
	if err != nil {
		t.Fatal(err)
	}
	defer env.DeleteLocalRef(pt)

	err = pt.SetField(env, "x", 5)
	if err != nil {
		t.Fatal(err)
	}
	var gotX int
	err = pt.GetField(env, "x", &gotX)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, 5, gotX) {
		t.Fail()
	}

	src := "fromChar"
	dest := make([]uint16, len(src))
	for i, c := range src {
		dest[i] = uint16(c)
	}
	str, err = env.NewObject("java/lang/String", dest)
	if err != nil {
		t.Fatal(err)
	}
	defer env.DeleteLocalRef(str)

	if err := str.CallMethod(env, "getBytes", &goBytes, env.GetUTF8String()); err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, src, string(goBytes)) {
		t.Fail()
	}
}

func PTestAttach(t *testing.T) {
	x := make(chan byte)

	obj, err := env.NewObject("java/lang/Object")
	if err != nil {
		t.Fatal(err)
	}
	gObj := env.NewGlobalRef(obj)

	go func() {
		runtime.LockOSThread()
		nenv := jvm.AttachCurrentThread()
		t.Logf("%x", nenv.jniEnv)

		var v int
		if err := gObj.CallMethod(nenv, "hashCode", &v); err != nil {
			t.Fatal(err)
		}
		t.Logf("%d", v)
		if err := jvm.DetachCurrentThread(nenv); err != nil {
			t.Fatal(err)
		}
		runtime.UnlockOSThread()

		x <- 4
	}()

	<-x
	env.DeleteGlobalRef(gObj)
	env.DeleteLocalRef(obj)
}

func PTestObjectArrays(t *testing.T) {
	var subject = "splitXme"
	str, err := env.NewObject("java/lang/String", []byte(subject))
	if err != nil {
		t.Fatal(err)
	}
	defer env.DeleteLocalRef(str)
	regex, err := env.NewObject("java/lang/String", []byte("X"))
	if err != nil {
		t.Fatal(err)
	}
	defer env.DeleteLocalRef(regex)

	v := NewObjectArrayRef("java/lang/String")
	if err := str.CallMethod(env, "split", v, regex); err != nil {
		t.Fatal(err)
	}

	parts := env.FromObjectArray(v)
	var got []string
	for _, p := range parts {
		var part []byte
		if err := p.CallMethod(env, "getBytes", &part); err != nil {
			t.Fatal(err)
		}
		got = append(got, string(part))
	}
	if !assert.Equal(t, subject, strings.Join(got, "X")) {
		t.Fail()
	}

	array := env.ToObjectArray(parts, "java/lang/String")

	class := NewObjectRef("java/lang/Class")
	if err := array.CallMethod(env, "getClass", class); err != nil {
		t.Fatal(err)
	}
	defer env.DeleteLocalRef(class)
	jClassName := NewObjectRef("java/lang/String")
	if err := class.CallMethod(env, "getName", jClassName); err != nil {
		t.Fatal(err)
	}
	defer env.DeleteLocalRef(jClassName)
	var className []byte
	if err := jClassName.CallMethod(env, "getBytes", &className); err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, "[Ljava.lang.String;", string(className)) {
		t.Fail()
	}
}

type GoString string

func (g *GoString) ConvertToGo(obj *ObjectRef) error {
	defer env.DeleteLocalRef(obj)
	var goBytes []byte
	if err := obj.CallMethod(env, "getBytes", &goBytes); err != nil {
		return err
	}
	*g = GoString(goBytes)
	return nil
}

func (g *GoString) ConvertToJava() (obj *ObjectRef, err error) {
	return env.NewObject("java/lang/String", []byte(string(*g)))
}

func (g *GoString) GetClassName() string {
	return "java/lang/String"
}

func (g *GoString) IsArray() bool {
	return false
}

func PTestConvert(t *testing.T) {
	var testString GoString = "test string"
	str, err := env.NewObject("java/lang/String", &testString)
	if err != nil {
		t.Fatal(err)
	}
	defer env.DeleteLocalRef(str)

	var firstWord GoString
	if err := str.CallMethod(env, "substring", &firstWord, 0, 4); err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, GoString("test"), firstWord) {
		t.Fail()
	}
}

func PTestInstanceOf(t *testing.T) {
	alist, err := env.NewObject("java/util/ArrayList")
	if err != nil {
		t.Fatal(err)
	}
	defer env.DeleteLocalRef(alist)

	str, err := env.NewObject("java/lang/String")
	if err != nil {
		t.Fatal(err)
	}
	defer env.DeleteLocalRef(str)

	var dummy bool
	if err := alist.CallMethod(env, "add", &dummy, str.Cast("java/lang/Object")); err != nil {
		t.Fatal(err)
	}

	obj := NewObjectRef("java/lang/Object")
	if err := alist.CallMethod(env, "get", obj, 0); err != nil {
		t.Fatal(err)
	}
	defer env.DeleteLocalRef(obj)

	if isInstance, err := obj.IsInstanceOf(env, "java/lang/String"); err != nil {
		t.Fatal(err)
	} else if !isInstance {
		t.Error("InstanceOf test failed")
	}
}

func PTestByteArray(t *testing.T) {
	ba := env.NewByteArray(5)
	bytes := ba.GetCritical(env)
	copy(bytes, []byte("hello"))
	ba.ReleaseCritical(env, bytes)
	str, err := env.NewObject("java/lang/String", ba)
	if err != nil {
		t.Fatal(err)
	}
	if !assert.Equal(t, "hello", toGoStr(t, str)) {
		t.Fail()
	}
	env.DeleteLocalRef(str)

	testStr := "hello world"
	str, err = env.NewObject("java/lang/String", []byte(testStr))
	if err != nil {
		t.Fatal(err)
	}
	defer env.DeleteLocalRef(str)

	arr := NewArrayRef(Byte | Array)
	if err := str.CallMethod(env, "getBytes", arr, env.GetUTF8String()); err != nil {
		t.Fatal(err)
	}

	ba2 := env.NewByteArrayFromObject(arr.ObjectRef)

	bytes = ba2.CopyBytes(env)
	if !assert.Equal(t, "hello world", string(bytes)) {
		t.Fail()
	}

	ba3 := env.NewByteArrayFromSlice([]byte("hello world!"))
	bytes = ba3.CopyBytes(env)
	if !assert.Equal(t, "hello world!", string(bytes)) {
		t.Fail()
	}
}

func PTestGetJVM(t *testing.T) {
	_, err := env.GetJVM()
	if err != nil {
		t.Fatalf("GetJavaVM failed %s", err)
	}
	t.Logf("Call GetJavaJVM: passed")
}

func PTestDestroy(t *testing.T) {
	err := jvm.Destroy()
	if err != nil {
		t.Fatalf("DestroyJVM failed %s", err)
	}
}

func PTestEnsureLocalCapacity(t *testing.T) {
	if err := env.EnsureLocalCapacity(256); err != nil {
		t.Fatalf("EnsureLocalCapacity failed %s", err)
	}
	t.Logf("Call EnsureLocalCapacity: passed")
}

func PTestPushPopLocalFrame(t *testing.T) {
	if err := env.PushLocalFrame(64); err != nil {
		t.Fatalf("PushLocalFrame failed %s", err)
	}
	t.Logf("Call PushLocalFrame: passed")

	obj, err := env.NewObject("java/lang/Object")
	if err != nil {
		t.Fatal(err)
	}

	var dummy int
	if err := obj.CallMethod(env, "hashCode", &dummy); err != nil {
		t.Fatal(err)
	}

	// Pop local frame with obj reference; obj should now be in previous frame
	obj = env.PopLocalFrame(obj)
	t.Logf("Call PopLocalFrame: passed")

	if err := obj.CallMethod(env, "hashCode", &dummy); err != nil {
		t.Fatalf("hashCode after PopLocalFrame failed %s", err)
	}

	env.DeleteLocalRef(obj)

	// Now do again with nil argument to pop
	if err := env.PushLocalFrame(32); err != nil {
		t.Fatalf("PushLocalFrame failed %s", err)
	}
	t.Logf("Call PushLocalFrame: passed")

	obj, err = env.NewObject("java/lang/Object")
	if err != nil {
		t.Fatal(err)
	}
	if err := obj.CallMethod(env, "hashCode", &dummy); err != nil {
		t.Fatal(err)
	}

	// Pop local frame with nil
	obj = env.PopLocalFrame(nil)
	t.Logf("Call PopLocalFrame: passed")

	if !obj.IsNil() {
		t.Fatal("PopLocalFrame return value is not nil")
	}
}

func PTestHandleException(t *testing.T) {
	jexceptErrMsg := "Java exception occurred. check stderr/logcat"
	if _, err := env.NewObject("java/foo/bar"); err == nil {
		t.Fatal("did not return error")
	} else if !assert.Equal(t, jexceptErrMsg, err.Error()) {
		t.Fatal("did not return standard error")
	}

	env.ExceptionHandler = ThrowableToStringExceptionHandler
	if _, err := env.NewObject("java/foo/bar"); err == nil {
		t.Fatal("did not return error")
	} else if !assert.Equal(t, "java.lang.NoClassDefFoundError: java/foo/bar", err.Error()) {
		t.Fatal("did not return standard error")
	}

	env.ExceptionHandler = ThrowableErrorExceptionHandler
	if _, err := env.NewObject("java/foo/bar"); err == nil {
		t.Fatal("did not return error")
	} else {
		throwableError, ok := err.(ThrowableError)
		if !ok {
			t.Fatalf("expected ThrowableError, but got %T", err)
		}

		// get the cause stack trace, then set this to nil so rest of structure can be tested
		var causeST []StackTraceElement
		if v := throwableError.Cause; v != nil {
			causeST = v.StackTrace
			throwableError.Cause.StackTrace = nil
		}

		want := ThrowableError{
			ClassName:        "java.lang.NoClassDefFoundError",
			LocalizedMessage: "java/foo/bar",
			Message:          "java/foo/bar",
			StackTrace:       []StackTraceElement{},
			AsString:         "java.lang.NoClassDefFoundError: java/foo/bar",
			Cause: &ThrowableError{
				ClassName:        "java.lang.ClassNotFoundException",
				LocalizedMessage: "java.foo.bar",
				Message:          "java.foo.bar",
				StackTrace:       nil,
				AsString:         "java.lang.ClassNotFoundException: java.foo.bar",
				Cause:            (*ThrowableError)(nil),
			},
		}

		if !assert.Equal(t, want, throwableError) {
			t.Fail()
		}
		if !assert.Equal(t, "java.lang.NoClassDefFoundError: java/foo/bar", throwableError.Error()) {
			t.Fail()
		}
		if !assert.Equal(t, "java.lang.ClassNotFoundException: java.foo.bar", throwableError.Cause.Error()) {
			t.Fail()
		}

		if !assert.NotNil(t, causeST) {
			t.Fail()
		} else {
			for _, v := range causeST {
				if !assert.NotEmpty(t, v) {
					t.Fail()
				}
			}
		}

	}

	env.ExceptionHandler = nil

	if _, err := env.NewObject("java/foo/bar"); err == nil {
		t.Fatal("did not return error")
	} else if !assert.Equal(t, jexceptErrMsg, err.Error()) {
		t.Error("did not return standard error")
	}
}

func PTestCast(t *testing.T) {
	str, err := env.NewObject("java/lang/String", []byte("hello world"))
	if err != nil {
		t.Fatal(err)
	}
	// create new object ref with class name java/lang/Foo
	c := &ObjectRef{str.JObject(), "java/lang/Foo", false}
	var goBytes []byte
	if err := c.Cast("java/lang/String").CallMethod(env, "getBytes", &goBytes, env.GetUTF8String()); err != nil {
		t.Fatal(err)
	}
}

func PTestRegisterNative(t *testing.T) {
	if err := env.RegisterNative("local/JnigiTesting", "Greet", ObjectType("java/lang/String"), []interface{}{ObjectType("java/lang/String")}, c_go_callback_Greet); err != nil {
		t.Fatal(err)
	}
	objRef, err := env.NewObject("local/JnigiTesting")
	if err != nil {
		t.Fatal(err)
	}
	nameRef, err := env.NewObject("java/lang/String", []byte("World"))
	if err != nil {
		t.Fatal(err)
	}
	strRef := NewObjectRef("java/lang/String")
	if err := objRef.CallMethod(env, "Greet", strRef, nameRef); err != nil {
		t.Fatal(err)
	}
	goStr := toGoStr(t, strRef)
	if !assert.Equal(t, "Hello World!", goStr) {
		t.Fail()
	}
}

func toGoStr(t *testing.T, o *ObjectRef) string {
	var goBytes []byte
	if err := o.CallMethod(env, "getBytes", &goBytes); err != nil {
		t.Fatal(err)
	}
	return string(goBytes)
}

func fromGoStr(t *testing.T, str string) *ObjectRef {
	jstr, err := env.NewObject("java/lang/String", []byte(str))
	if err != nil {
		t.Fatal(err)
	}
	return jstr
}
