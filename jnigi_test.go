// Copyright 2016 Tim O'Brien. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jnigi

import (
	"testing"
)

var env *Env
var jvm *JVM

func TestInit(t *testing.T) {
	jvm2, e2, err := CreateJVM(NewJVMInitArgs(false, true, DEFAULT_VERSION, []string{"-Xcheck:jni"}))
	if err != nil {
		t.Fatal(err)
	}
	env = e2
	jvm = jvm2

	t.Logf("%x", e2.jniEnv)

}

func toGoStr(t *testing.T, o *ObjectRef) string {
	v, err := o.CallMethod(env, "getBytes", Byte|Array)
	if err != nil {
		t.Fatal(err)
	}
	return string(v.([]byte))
}

func fromGoStr(t *testing.T, str string) *ObjectRef {
	jstr, err := env.NewObject("java/lang/String", []byte(str))
	if err != nil {
		t.Fatal(err)
	}
	return jstr
}

func TestBasic(t *testing.T) {
	// new object, int method
	obj, err := env.NewObject("java/lang/Object")
	if err != nil {
		t.Fatal(err)
	}
	v, err := obj.CallMethod(env, "hashCode", Int)
	if err != nil {
		t.Fatal(err)
	}

	// byte array argument, byte array method
	testStr := "hello world"
	str, err := env.NewObject("java/lang/String", []byte(testStr))
	if err != nil {
		t.Fatal(err)
	}
	v, err = str.CallMethod(env, "getBytes", Byte|Array, env.GetUTF8String())
	if err != nil {
		t.Fatal(err)
	}
	if b, ok := v.([]byte); !ok || string(b) != testStr {
		t.Logf("basic test failed")
	}

	// object method, int arg, object arg
	v, err = str.CallMethod(env, "substring", "java/lang/String", 6)
	if err != nil {
		t.Fatal(err)
	}
	str2 := v.(*ObjectRef)
	v, err = str2.CallMethod(env, "getBytes", Byte|Array)
	if err != nil {
		t.Fatal(err)
	}

	env.PrecalculateSignature("(Ljava/lang/String;)Z")
	v, err = str.CallMethod(env, "endsWith", Boolean, str2)
	if err != nil {
		t.Fatal(err)
	}
	if b, ok := v.(bool); !ok || !b {
		t.Logf("basic test failed")
	}

	// call static method
	v, err = env.CallStaticMethod("java/lang/System", "getProperty", "java/lang/String", fromGoStr(t, "java.vm.version"))
	t.Logf(toGoStr(t, v.(*ObjectRef)))

	// get static field
	v, err = env.GetStaticField("java/util/Calendar", "APRIL", Int)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("april = %d", v.(int))

	// set/get object field
	pt, err := env.NewObject("java/awt/Point")
	if err != nil {
		t.Fatal(err)
	}
	err = pt.SetField(env, "x", 5)
	if err != nil {
		t.Fatal(err)
	}
	v, err = pt.GetField(env, "x", Int)
	if err != nil {
		t.Fatal(err)
	}
	if i, ok := v.(int); !ok || i != 5 {
		t.Logf("basic test failed")
	}
}

func TestAttach(t *testing.T) {
	x := make(chan byte)

	obj, err := env.NewObject("java/lang/Object")
	if err != nil {
		t.Fatal(err)
	}
	gObj := env.NewGlobalRef(obj)

	go func() {
		nenv := jvm.AttachCurrentThread()
		t.Logf("%x", nenv.jniEnv)

		v, err := gObj.CallMethod(nenv, "hashCode", Int)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%d", v.(int))

		x <- 4
	}()

	<-x
}

func TestObjectArrays(t *testing.T) {
	str, err := env.NewObject("java/lang/String", []byte("splitXme"))
	if err != nil {
		t.Fatal(err)
	}

	regex, err := env.NewObject("java/lang/String", []byte("X"))
	if err != nil {
		t.Fatal(err)
	}

	v, err := str.CallMethod(env, "split", ObjectArrayType("java/lang/String"), regex)
	if err != nil {
		t.Fatal(err)
	}

	parts := env.FromObjectArray(v.(*ObjectRef))
	for _, p := range parts {
		v, err = p.CallMethod(env, "getBytes", Byte|Array)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%s", string(v.([]byte)))
	}

	array := env.ToObjectArray(parts, "java/lang/String")

	v, err = array.CallMethod(env, "getClass", "java/lang/Class")
	if err != nil {
		t.Fatal(err)
	}
	v, err = v.(*ObjectRef).CallMethod(env, "getName", "java/lang/String")
	if err != nil {
		t.Fatal(err)
	}
	v, err = v.(*ObjectRef).CallMethod(env, "getBytes", Byte|Array)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", v.([]byte))
}

func TestInstanceOf(t *testing.T) {
	alist, err := env.NewObject("java/util/ArrayList")
	if err != nil {
		t.Fatal(err)
	}

	str, err := env.NewObject("java/lang/String")
	if err != nil {
		t.Fatal(err)
	}
	_, err = alist.CallMethod(env, "add", Boolean, str.Cast("java/lang/Object"))
	if err != nil {
		t.Fatal(err)
	}

	v, err := alist.CallMethod(env, "get", "java/lang/Object", 0)
	if err != nil {
		t.Fatal(err)
	}
	obj := v.(*ObjectRef)

	if v, err := obj.IsInstanceOf(env, "java/lang/String"); err != nil {
		t.Fatal(err)
	} else if !v {
		t.Fatal("instanceof test failed")
	}
}