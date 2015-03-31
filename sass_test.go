// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSass(t *testing.T) {

	p, err := AbsPath("src/github.com/dshills/goauto/testing/_sub.scss")
	assert.Nil(t, err)
	css := filepath.Join(filepath.Dir(p), "css")
	cache := filepath.Join(css, ".sass_cache")

	st := NewSassTask(css, cache, "compressed")
	ti := TaskInfo{
		Src:  p,
		Tout: os.Stdout,
		Terr: os.Stderr,
	}
	err = st.Run(&ti)
	assert.Nil(t, err)

	_, err = os.Stat(css)
	assert.Nil(t, err)

	_, err = os.Stat(cache)
	assert.Nil(t, err)

	nc := filepath.Join(filepath.Dir(p), "css", "main.css")
	_, err = os.Stat(nc)
	assert.Nil(t, err)
	os.Remove(nc)

	ncm := filepath.Join(filepath.Dir(p), "css", "main.css.map")
	_, err = os.Stat(ncm)
	assert.Nil(t, err)
	os.Remove(ncm)
}
