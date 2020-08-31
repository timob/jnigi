# JNIGI
Java Native Interface Golang Interface.

A package to access Java from Golang code. Can be used from a Golang executable or shared library.
This allows for Golang to initiate the JVM or Java to start a Golang runtime respectively.

## Compile
The `CGO_CFLAGS` needs to be set to add the JNI C header files. The `compilevars.sh` script will do
this.
```
# put this in your build script
source <gopath>/src/github.com/timob/jnigi/compilevars.sh <root path of jdk>
```

On Windows you can use `compilevars.bat` in the same way (but you don't need `source` at the begining).


## Finding JVM at Runtime
Use the `LoadJVMLib(jvmLibPath string) error` function to load the shared library at run time.
There is a function `AttemptToFindJVMLibPath() string` to help to find the library path.

## Notes
### Signals
The JVM calls sigaction, Golang requires that SA_ONSTACK flag be passed.
Without this Golang will not be able to print the exception information eg. type, stack trace, line number.
On Linux a solution is using LD_PRELOAD with a library that intercepts sigaction and adds the flag. (Code that does this: https://gist.github.com/timob/5d3032b54ed6ba2dc6de34b245c556c7)

### MacOS AWT/Swing

In order to run AWT/Swing on MacOS (only) is requires to run a platform specific Cocoa Run Loop code otherwise java gui becomes irresponsive.

You would need to:

1) Spin go routine in a *separate* thread from which is initialized JVM 
2) In the application's main thread run `RunJVMGUILoop` function which configures and executes Cocoa run loop. 
Running the loop is a blocking opperation for the current (main) thread. 

Note it's not JNIGI limitation but a specific Apple's requirement on MacOS. If you do not execute AWT/Swing you do not need such setup.

On other platforms `RunJVMGUILoop` is emulating MacOS behavior for *consistency* by blocking the current execution thread.

## Status
* Has been used in Golang (many versions since 1.6) executable multi threaded applications on Linux / Windows.
* Tests for main functions tests are present.
* Documentation needed.

## Changes
* 2020-08-31 Add AWT/Swing MacOS support, workaround for JDK-7131356 "No Java runtime present" 
* 2020-08-11 Add DestroyJavaVM support, JNI_VERSION_1_8 const
* 2019-05-29 Better multiplatform support, dynamic loading of JVM library.
* 2016-08-01 Initial version.

## Example

```` go
package main

import (
    "fmt"
    "github.com/timob/jnigi"
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
    obj, err := env.NewObject("java/lang/Object")
    if err != nil {
        log.Fatal(err)
    }
    v, err := obj.CallMethod(env, "hashCode", jnigi.Int)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("object hash code: %d\n", v.(int))
    if err := jvm.Destroy(); err != nil {
        log.Fatal(err)
    }
}

````
