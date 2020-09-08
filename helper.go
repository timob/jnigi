package jnigi

import (
	"os"
	"runtime"
	"path/filepath"
)

// AttemptToFindJVMLibPath tries to find the full path to the JVM shared library file
func AttemptToFindJVMLibPath() string {
	// for these linux is the "default"

	prefix := os.Getenv("JAVA_HOME")
	if prefix == "" {
		if runtime.GOOS == "windows" {
			prefix = filepath.Join("c:", "Program Files", "Java", "jdk")
		} else if runtime.GOOS == "darwin" {
			prefix = "/Library/Java/Home"
		} else {
			prefix = "/usr/lib/jvm/default-java"
		}
	}

	dirPath := prefix
	if runtime.GOOS == "windows" {
		dirPath = filepath.Join(dirPath, "bin", "server")
	} else if runtime.GOOS == "darwin" {
		dirPath = filepath.Join(dirPath, "lib", "server")
	} else {
		dirPath = filepath.Join(dirPath, "jre", "lib", runtime.GOARCH, "server")
	}

	var libPath string
	if runtime.GOOS == "windows" {
		libPath = filepath.Join(dirPath, "jvm.dll")
	} else if runtime.GOOS == "darwin" {
		libPath = filepath.Join(dirPath, "libjvm.dylib")
	} else {
		libPath = filepath.Join(dirPath, "libjvm.so")
	}
	return libPath
}