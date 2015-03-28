// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

// Package goauto implements a set of tools for building workflow automation tools.
// These tools can be as simple as running a compiler when a source file changes to complex chains of tasks doing almost any action required within a development environment
// See README.md for more details on usage
package goauto

// Verbose is a global var that will print a lot of debug info during processing
// This is handy for debugging. By default it is off
var Verbose bool

func init() {
	Verbose = false
}

/*

TODO

More built in Go and Shell tasks
Write more tests. Can always use more tests
Call your mother







Really need to find some cool ASCII art for this file









*/
