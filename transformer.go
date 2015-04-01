// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
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
	gopaths := GoPaths()
	for _, gp := range gopaths {
		rel, err := filepath.Rel(gp, filepath.Dir(f))
		if err == nil {
			return rel
		}
	}
	return ""
}

// GoRelSrcDir returns the path relative to $GOPATH/src
func GoRelSrcDir(f string) string {
	gopaths := GoPaths()
	for _, gp := range gopaths {
		gp = filepath.Join(gp, "src")
		rel, err := filepath.Rel(gp, filepath.Dir(f))
		if err == nil {
			return rel
		}
	}
	return ""
}

// ExtTransformer returns a Transformer for chaning file extensions
func ExtTransformer(newExt string) Transformer {
	return func(f string) string {
		b := strings.TrimSuffix(f, filepath.Ext(f))
		return b + "." + newExt
	}
}
