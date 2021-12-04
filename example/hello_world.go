package main

import (
	"fmt"
	"tekao.net/jnigi"
	"log"
)

func main() {
	if err := jnigi.LoadJVMLib(jnigi.AttemptToFindJVMLibPath()); err != nil {
		log.Fatal(err)
	}
	jvm, env, err := jnigi.CreateJVM(jnigi.NewJVMInitArgs(false, true, jnigi.DEFAULT_VERSION, []string{"-Xcheck:jni"}))
	if err != nil {
		log.Fatal(err)
	}

	hello, err := env.NewObject("java/lang/String", []byte("Hello "))
	if err != nil {
		log.Fatal(err)
	}

	world, err := env.NewObject("java/lang/String", []byte("World!"))
	if err != nil {
		log.Fatal(err)
	}

	greeting := NewObjectRef("java/lang/String")
	err = hello.CallMethod(env, "concat", &greeting, world)
	if err != nil {
		log.Fatal(err)
	}

	var goGreeting []byte
	err = greeting.CallMethod(env, "getBytes", &goGreeting)
	if err != nil {
		log.Fatal(err)
	}

	// Prints "Hello World!"
	fmt.Printf("%s\n", goGreeting)

	if err := jvm.Destroy(); err != nil {
		log.Fatal(err)
	}
}
