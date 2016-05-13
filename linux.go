// Copyright 2016 Tim O'Brien. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux darwin

package jnigi

/*
#cgo CFLAGS:-I../include/ -I/usr/lib/jvm/default-java/include
#cgo LDFLAGS:-ljvm -L/usr/lib/jvm/default-java/jre/lib/amd64/server

*/
import "C"
