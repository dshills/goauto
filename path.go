// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AbsPath is a utility function to get the absolute path from a path
// It will first check for an absolute path then GOPATH relative then a pwd relative
// will return an error for a path that does not exist
func AbsPath(path string) (ap string, err error) {

	// Check for absolute path
	if filepath.IsAbs(path) {
		ap = filepath.Clean(path)
		if _, err := os.Stat(ap); err == nil {
			return ap, nil
		}
	}

	// Check for GOPATH relative
	for _, gp := range GoPaths() {
		ap = filepath.Clean(filepath.Join(gp, path))
		_, err = os.Stat(ap)
		if err == nil {
			return
		}
	}

	ap, err = filepath.Abs(path)
	if err == nil {
		if _, err = os.Stat(ap); err == nil {
			return
		}
	}

	return "", fmt.Errorf("%q: no such file or directory", path)
}

// GoPaths retreives the GOPATH env and splits it []string
func GoPaths() []string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return []string{}
	}
	return strings.Split(gopath, ":")
}

// IsHidden is a HACKY check for hidden directory name
func IsHidden(d string) bool {
	if d[:1] == "." {
		return true
	}
	return false
}
