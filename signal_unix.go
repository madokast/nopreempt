// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build unix

package nopreempt

// gsignalStack saves the fields of the gsignal stack changed by
// setGsignalStack.
type gsignalStack struct {
	stack       stack
	stackguard0 uintptr
	stackguard1 uintptr
	stktopsp    uintptr
}
