// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.
// +build !darwin

package goauto

// NewWatchOSX returns a standard WatchFS, if OS is not OSX
func NewWatchOSX() Watcher {
	return NewWatchFS()
}
