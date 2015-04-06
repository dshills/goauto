// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import "testing"

func TestTransforms(t *testing.T) {
	var e, i, o string

	i = "TRANSFORM"
	o = Identity(i)
	if i != o {
		t.Errorf("Expected %v got %v", i, o)
	}

	/* Only works local, obviously
	i = "/Users/dshills/Development/Go/src/github.com/dshills/goauto/transform.go"

	e = "src/github.com/dshills/goauto/transform.go"
	o = GoRelBase(i)
	if o != e {
		t.Errorf("Expected %v got %v", e, o)
	}

	o = GoRelDir(i)
	e = "src/github.com/dshills/goauto"
	if o != e {
		t.Errorf("Expected %v got %v", e, o)
	}

	o = GoRelSrcDir(i)
	e = "github.com/dshills/goauto"
	if o != e {
		t.Errorf("Expected %v got %v", e, o)
	}
	*/

	i = "transform.go"
	fx := ExtTransformer("js")
	o = fx(i)
	e = "transform.js"
	if o != e {
		t.Errorf("Expected %v got %v", e, o)
	}
}
