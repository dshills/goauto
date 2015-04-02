// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoPaths(t *testing.T) {
	gps := GoPaths()
	assert.NotEmpty(t, gps)
}

func TestAbsPath(t *testing.T) {
	tp := filepath.Join("src", "github.com", "dshills", "goauto")
	ap, err := AbsPath(tp)
	assert.Nil(t, err)
	assert.NotEqual(t, ap, tp)

	/* local test
	assert.Equal(t, "/Users/dshills/Development/Go/src/github.com/dshills/goauto", ap)
	*/

}
