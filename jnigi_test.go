// Copyright 2016 Tim O'Brien. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jnigi

import (
	"testing"
	"runtime"
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
	jvm2, e2, err := CreateJVM(NewJVMInitArgs(false, true, DEFAULT_VERSION, []string{"-Xcheck:jni"}))
	if err != nil {
		t.Fatal(err)
	}
	env = e2
	jvm = jvm2

	t.Logf("%x", e2.jniEnv)

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

	// byte array argument, byte array method
	testStr := "hello world"
	str, err := env.NewObject("java/lang/String", []byte(testStr))
	if err != nil {
		t.Fatal(err)
	}

	var goBytes []byte
	if err := str.CallMethod(env, "getBytes", &goBytes, env.GetUTF8String()); err != nil {
		t.Fatal(err)
	}
	if string(goBytes) != testStr {
		t.Errorf("basic test failed")
	}

	// object method, int arg, object arg
	str2 := NewObjectRef("java/lang/String")
	if err := str.CallMethod(env, "substring", str2, 6); err != nil {
		t.Fatal(err)
	}

	var dummy []byte
	if err := str2.CallMethod(env, "getBytes", &dummy); err != nil {
		t.Fatal(err)
	}

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
	t.Logf(toGoStr(t, jvmVer))

	// get static field
	var calPos int
	err = env.GetStaticField("java/util/Calendar", "APRIL", &calPos)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("april = %d", calPos)

	// set/get object field
	pt, err := env.NewObject("java/awt/Point")
	if err != nil {
		t.Fatal(err)
	}
	err = pt.SetField(env, "x", 5)
	if err != nil {
		t.Fatal(err)
	}
	var gotX int
	err = pt.GetField(env, "x", &gotX)
	if err != nil {
		t.Fatal(err)
	}
	if gotX != 5 {
		t.Errorf("basic test failed")
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

	if err := str.CallMethod(env, "getBytes", &goBytes, env.GetUTF8String()); err != nil {
		t.Fatal(err)
	}
	if string(goBytes) != src {
		t.Errorf("basic test failed")
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
}

func PTestObjectArrays(t *testing.T) {
	str, err := env.NewObject("java/lang/String", []byte("splitXme"))
	if err != nil {
		t.Fatal(err)
	}

	regex, err := env.NewObject("java/lang/String", []byte("X"))
	if err != nil {
		t.Fatal(err)
	}

	v := NewObjectArrayRef("java/lang/String")
	if err := str.CallMethod(env, "split", v, regex); err != nil {
		t.Fatal(err)
	}

	parts := env.FromObjectArray(v)
	for _, p := range parts {
		var part []byte
		if err := p.CallMethod(env, "getBytes", &part); err != nil {
			t.Fatal(err)
		}
		t.Logf("%s", string(part))
	}

	array := env.ToObjectArray(parts, "java/lang/String")

	class := NewObjectRef("java/lang/Class")
	if err := array.CallMethod(env, "getClass", class); err != nil {
		t.Fatal(err)
	}
	jClassName := NewObjectRef("java/lang/String")
	if err := class.CallMethod(env, "getName", jClassName); err != nil {
		t.Fatal(err)
	}
	var className []byte
	if err := jClassName.CallMethod(env, "getBytes", &className); err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", string(className))
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

	var firstWord GoString
	if err := str.CallMethod(env, "substring", &firstWord, 0, 4); err != nil {
		t.Fatal(err)
	}
	if firstWord != "test" {
		t.Errorf("convert test failed got %s", firstWord)
	}
}

func PTestInstanceOf(t *testing.T) {
	alist, err := env.NewObject("java/util/ArrayList")
	if err != nil {
		t.Fatal(err)
	}

	str, err := env.NewObject("java/lang/String")
	if err != nil {
		t.Fatal(err)
	}
	var dummy bool
	if err := alist.CallMethod(env, "add", &dummy, str.Cast("java/lang/Object")); err != nil {
		t.Fatal(err)
	}

	obj := NewObjectRef("java/lang/Object")
	if err := alist.CallMethod(env, "get", obj, 0); err != nil {
		t.Fatal(err)
	}

	if isInstance, err := obj.IsInstanceOf(env, "java/lang/String"); err != nil {
		t.Fatal(err)
	} else if !isInstance {
		t.Fatal("instanceof test failed")
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
	if toGoStr(t, str) != "hello" {
		t.Fatal("ByteArray test failed")
	}

	testStr := "hello world"
	str, err = env.NewObject("java/lang/String", []byte(testStr))
	if err != nil {
		t.Fatal(err)
	}

	arr := NewArrayRef(Byte|Array)
	if err := str.CallMethod(env, "getBytes", arr, env.GetUTF8String()); err != nil {
		t.Fatal(err)
	}

	ba2 := env.NewByteArrayFromObject(arr.ObjectRef)
	bytes = ba2.GetCritical(env)
	if string(bytes) != "hello world" {
		t.Logf("ByteArray test failed")
	}
	ba2.ReleaseCritical(env, bytes)
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
	t.Logf("Call DestroyJVM: passed")
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

	if _, err := env.NewObject("java/foo/bar"); err == nil {
		t.Fatal("did not return error")
	} else if err.Error() != "Java exception occured. check stderr/logcat" {
		t.Fatalf("did not return standard error: %v", err)
	}

	env.ExceptionHandler = ThrowableToStringExceptionHandler
	if _, err := env.NewObject("java/foo/bar"); err == nil {
		t.Fatal("did not return error")
	} else if err.Error() != "java.lang.NoClassDefFoundError: java/foo/bar" {
		t.Fatalf("unexpected result of ToString: %v", err)
	}

	env.ExceptionHandler = ThrowableErrorExceptionHandler
	if _, err := env.NewObject("java/foo/bar"); err == nil {
		t.Fatal("did not return error")
	} else {

		throwableError, ok := err.(ThrowableError)
		if !ok {
			t.Fatalf("expected ThrowableError, but got %T", err)
		}

		if err.Error() != "java.lang.NoClassDefFoundError: java/foo/bar" {
			t.Fatalf("unexpected error message: %v", err)
		}

		if v := throwableError.ClassName; v != "java.lang.NoClassDefFoundError" {
			t.Fatalf("unexpected class name: %s", v)
		}

		if v := throwableError.LocalizedMessage; v != "java/foo/bar" {
			t.Fatalf("unexpected localized message: %s", v)
		}

		if v := throwableError.Message; v != "java/foo/bar" {
			t.Fatalf("unexpected message: %s", v)
		}

		if v := throwableError.AsString; v != "java.lang.NoClassDefFoundError: java/foo/bar" {
			t.Fatalf("unexpected toString value: %s", v)
		}

		if v := throwableError.StackTrace; len(v) > 0 {
			t.Fatal("expect empty stack trace")
		}

		if throwableError.Cause == nil {
			t.Fatal("expected a cause")
		}

		cause := throwableError.Cause

		if cause.Error() != "java.lang.ClassNotFoundException: java.foo.bar" {
			t.Fatalf("unexpected error message: %v", cause)
		}

		if v := cause.ClassName; v != "java.lang.ClassNotFoundException" {
			t.Fatalf("unexpected class name: %s", v)
		}

		if v := cause.LocalizedMessage; v != "java.foo.bar" {
			t.Fatalf("unexpected localized message: %s", v)
		}

		if v := cause.Message; v != "java.foo.bar" {
			t.Fatalf("unexpected message: %s", v)
		}

		if v := cause.AsString; v != "java.lang.ClassNotFoundException: java.foo.bar" {
			t.Fatalf("unexpected toString value: %s", v)
		}

		if v := cause.StackTrace; v == nil {
			t.Fatal("expected a stack trace")
		} else if len(cause.StackTrace) == 0 {
			t.Fatal("expected stack trace entries")
		}

		for i, v := range cause.StackTrace {
			if v.AsString == "" {
				t.Fatalf("stack trace index %d: no AsString value", i)
			}
			if v.ClassName == "" {
				t.Fatalf("stack trace index %d: no ClassName value", i)
			}
			if v.FileName == "" {
				t.Fatalf("stack trace index %d: no FileName value", i)
			}
			if v.MethodName == "" {
				t.Fatalf("stack trace index %d: no MethodName value", i)
			}
			if v.LineNumber == 0 {
				t.Fatalf("stack trace index %d: no LineNumber value", i)
			}
		}
	}

	env.ExceptionHandler = nil

	if _, err := env.NewObject("java/foo/bar"); err == nil {
		t.Fatal("did not return error")
	} else if err.Error() != "Java exception occured. check stderr/logcat" {
		t.Fatalf("did not return standard error: %v", err)
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