// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"os"
	"path/filepath"
	"strings"
)

// A Transformer is a function that changes a string to a string
// Transforms are used to build Targets in Tasks
// A common example would be to change myfile.scss to myfile.css
type Transformer func(string) string

// Identity returns itself
func Identity(f string) string {
	return f
}

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

// ExtTransformer returns a Transformer for chaning file extensions
func ExtTransformer(newExt string) Transformer {
	return func(f string) string {
		b := strings.TrimSuffix(f, filepath.Ext(f))
		return b + "." + newExt
	}
}
