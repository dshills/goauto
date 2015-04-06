// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"path/filepath"
	"testing"
)

func TestGoPaths(t *testing.T) {
	gps := GoPaths()
	if len(gps) < 1 {
		t.Errorf("GoPaths returned no data\n")
	}
}

func TestAbsPath(t *testing.T) {
	tp := filepath.Join("src", "github.com", "dshills", "goauto")
	ap, err := AbsPath(tp)
	if err != nil {
		t.Errorf("AbsPath error: %v\n", err)
	}
	if ap == tp {
		t.Errorf("AbsPath: %v should not equal %v\n", ap, tp)
	}

	expect := "/Users/dshills/Development/Go/src/github.com/dshills/goauto"
	if ap != expect {
		t.Errorf("Expected %v Got %v\n", expect, ap)
	}

}
