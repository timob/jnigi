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

	v, err := hello.CallMethod(env, "concat", jnigi.ObjectType("java/lang/String"), world)
	if err != nil {
		log.Fatal(err)
	}
	greeting := v.(*jnigi.ObjectRef)

	v, err = greeting.CallMethod(env, "getBytes", jnigi.Byte|jnigi.Array)
	if err != nil {
		log.Fatal(err)
	}
	goGreeting := string(v.([]byte))

	// Prints "Hello World!"
	fmt.Printf("%s", goGreeting)

	if err := jvm.Destroy(); err != nil {
		log.Fatal(err)
	}
}
