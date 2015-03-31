// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransforms(t *testing.T) {
	i := "TRANSFORM"
	o := Identity(i)
	assert.Equal(t, i, o)

	/* Only works local, obviously
	i = "/Users/dshills/Development/Go/src/github.com/dshills/goauto/transform.go"
	o = GoRelBase(i)
	assert.Equal(t, "src/github.com/dshills/goauto/transform.go", o)

	o = GoRelDir(i)
	assert.Equal(t, "src/github.com/dshills/goauto", o)

	o = GoRelSrcDir(i)
	assert.Equal(t, "github.com/dshills/goauto", o)
	*/

	i = "transform.go"
	fx := ExtTransformer("js")
	o = fx(i)
	assert.Equal(t, "transform.js", o)
}
