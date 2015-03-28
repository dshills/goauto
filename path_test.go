// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoPaths(t *testing.T) {
	gps := GoPaths()
	assert.NotEmpty(t, gps)
}

func TestAbsPath(t *testing.T) {
	var ap, tp string

	tp = "/usr/bin"
	ap, err := AbsPath(tp)
	assert.Nil(t, err)
	assert.Equal(t, ap, tp)

	tp = "src/github.com/dshills/goauto"
	ap, err = AbsPath(tp)
	assert.Nil(t, err)
	assert.NotEqual(t, ap, tp)

}
