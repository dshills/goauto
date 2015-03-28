// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"os"
	"path/filepath"
)

// GoRelBase returns the file path relative to $GOPATH
func GoRelBase(f string) string {
	rd := GoRelDir(f)
	return filepath.Join(rd, filepath.Base(f))
}

// GoRelDir returns the path relative to $GOPATH
func GoRelDir(f string) string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return ""
	}
	rel, err := filepath.Rel(gopath, filepath.Dir(f))
	if err != nil {
		return ""
	}
	return rel
}

// GoRelSrcDir returns the path relative to $GOPATH/src
func GoRelSrcDir(f string) string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return ""
	}
	gopath = filepath.Join(gopath, "src")
	rel, err := filepath.Rel(gopath, filepath.Dir(f))
	if err != nil {
		return ""
	}
	return rel
}
