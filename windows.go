// Copyright 2016 Tim O'Brien. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package jnigi

/*
#cgo CFLAGS:-I../include/ -Ic:/oraclejdk/include -Ic:/oraclejdk/include/win32
#cgo LDFLAGS:-ljvm -Lc:/oraclejdk/jre/bin/server
*/
import "C"
