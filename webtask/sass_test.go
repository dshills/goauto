// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package webtask

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/dshills/goauto"
)

func TestSass(t *testing.T) {
	tp := filepath.Join("src", "github.com", "dshills", "goauto", "testing", "_sub.scss")
	p, err := goauto.AbsPath(tp)
	if err != nil {
		t.Error(err)
	}
	css := filepath.Join(filepath.Dir(p), "css")
	cache := filepath.Join(css, ".sass_cache")

	st := NewSassTask(css, cache, "compressed")
	ti := goauto.TaskInfo{
		Src:  p,
		Tout: ioutil.Discard,
		Terr: ioutil.Discard,
	}
	err = st.Run(&ti)
	if err != nil {
		t.Error(err)
	}

	_, err = os.Stat(css)
	if err != nil {
		t.Error(err)
	}

	_, err = os.Stat(cache)
	if err != nil {
		t.Error(err)
	}

	nc := filepath.Join(filepath.Dir(p), "css", "main.css")
	_, err = os.Stat(nc)
	if err != nil {
		t.Error(err)
	}
	os.Remove(nc)

	ncm := filepath.Join(filepath.Dir(p), "css", "main.css.map")
	_, err = os.Stat(ncm)
	if err != nil {
		t.Error(err)
	}
	os.Remove(ncm)
}
