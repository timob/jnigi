package jnigi

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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
		dirPath = filepath.Join(dirPath, "lib", "server")
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

func assignDest(val interface{}, dest interface{}) error {
	if dest == nil {
		return nil
	}

	var assigned bool
	switch v := val.(type) {
	case bool:
		if dv, ok := dest.(*bool); ok {
			*dv = v
			assigned = true
		}
	case byte:
		if dv, ok := dest.(*byte); ok {
			*dv = v
			assigned = true
		}
	case uint16:
		if dv, ok := dest.(*uint16); ok {
			*dv = v
			assigned = true
		}
	case int16:
		if dv, ok := dest.(*int16); ok {
			*dv = v
			assigned = true
		}
	case int:
		if dv, ok := dest.(*int); ok {
			*dv = v
			assigned = true
		}
	case int64:
		if dv, ok := dest.(*int64); ok {
			*dv = v
			assigned = true
		}
	case float32:
		if dv, ok := dest.(*float32); ok {
			*dv = v
			assigned = true
		}
	case float64:
		if dv, ok := dest.(*float64); ok {
			*dv = v
			assigned = true
		}
	case []bool:
		if dv, ok := dest.(*[]bool); ok {
			*dv = v
			assigned = true
		}
	case []byte:
		if dv, ok := dest.(*[]byte); ok {
			*dv = v
			assigned = true
		}
	case []uint16:
		if dv, ok := dest.(*[]uint16); ok {
			*dv = v
			assigned = true
		}
	case []int16:
		if dv, ok := dest.(*[]int16); ok {
			*dv = v
			assigned = true
		}
	case []int:
		if dv, ok := dest.(*[]int); ok {
			*dv = v
			assigned = true
		}
	case []int64:
		if dv, ok := dest.(*[]int64); ok {
			*dv = v
			assigned = true
		}
	case []float32:
		if dv, ok := dest.(*[]float32); ok {
			*dv = v
			assigned = true
		}
	case []float64:
		if dv, ok := dest.(*[]float64); ok {
			*dv = v
			assigned = true
		}
	case *ObjectRef:
		if dv, ok := dest.(*ObjectRef); ok {
			*dv = *v
			assigned = true
		}
	}

	if !assigned {
		return fmt.Errorf("JNIGI Error: expected dest argument to be %T (not %T) or nil", val, dest)
	}
	return nil
}
