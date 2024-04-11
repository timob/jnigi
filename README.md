# JNIGI
JNI Go Interface.

A package to access Java from Go code. Can be used from a Go executable or shared library.
This allows for Go to initiate the JVM or Java to start a Go runtime respectively.

[![Go Reference](https://pkg.go.dev/badge/github.com/timob/jnigi.svg)](https://pkg.go.dev/github.com/timob/jnigi)
[![Actions](https://github.com/timob/jnigi/actions/workflows/ci_test.yaml/badge.svg)](https://github.com/timob/jnigi/actions?query=branch%3Amaster)

## Module name change
The go module name is renamed to `github.com/timob/jnigi` in the branch. Checkout `v2` if you want to retain the old name.

## Install
``` bash
# In your apps Go module directory
go get github.com/timob/jnigi

# Add flags needed to include JNI header files, change this as appropriate for your JDK and OS
export CGO_CFLAGS="-I/usr/lib/jvm/default-java/include -I/usr/lib/jvm/default-java/include/linux"

# build your app
go build
```
The `compilevars.sh` (`compilevars.bat` on Windows) script can help setting the `CGO_CFLAGS` environment variable.

## Finding JVM at Runtime
The JVM library is dynamically linked at run time. Use the `LoadJVMLib(jvmLibPath string) error` function to load the shared library at run time.
There is a function `AttemptToFindJVMLibPath() string` to help to find the library path.

## Status
### Testing
Most of the code has tests. To run the tests using docker:
``` bash
# get source
git clone https://github.com/timob/jnigi.git

cd jnigi

# build image
docker build -t jnigi_test .

# run tests
docker run jnigi_test
```

### Uses
Has been used on Linux/Windows/MacOS (amd64) Android (arm64) multi threaded apps.

### Note about using on Windows
Because of the way the JVM triggers OS exceptions during `CreateJavaVM`, which the Go runtime treats as unhandled exceptions, the code for the Go runtime needs to be changed to allow the JVM to handle the exceptions. See https://github.com/timob/jnigi/issues/31#issuecomment-1668914368 for how to do this.

## Example

```` go
package main

import (
    "fmt"
    "github.com/timob/jnigi"
    "log"
    "runtime"
)

func main() {
    if err := jnigi.LoadJVMLib(jnigi.AttemptToFindJVMLibPath()); err != nil {
        log.Fatal(err)
    }

    runtime.LockOSThread()
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

    greeting := jnigi.NewObjectRef("java/lang/String")
    err = hello.CallMethod(env, "concat", greeting, world)
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
````
