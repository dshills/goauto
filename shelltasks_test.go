// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShellTask(t *testing.T) {
	info := TaskInfo{Src: "FROGS", Tout: os.Stdout, Terr: os.Stderr}

	tsk := NewShellTask(Identity, "echo")
	err := tsk.Run(&info)
	assert.Nil(t, err)
	assert.Equal(t, "FROGS\n", info.Buf.String())
}

func TestOSTasks(t *testing.T) {
	path, err := AbsPath("src/github.com/dshills/goauto")
	assert.Nil(t, err)
	fname := filepath.Join(path, "testdata")
	assert.Nil(t, err)
	f, err := os.Create(fname)
	defer f.Close()
	f.WriteString("TESTING")

	info := TaskInfo{Src: fname, Tout: os.Stdout, Terr: os.Stderr}

	newPath := filepath.Join(path, "t")
	tsk := NewMkdirTask(func(f string) string { return newPath })
	err = tsk.Run(&info)
	assert.Nil(t, err)

	info.Src = info.Target

	newfname := filepath.Join(newPath, "testdata2")
	tsk = NewCopyTask(func(string) string { return newfname })
	err = tsk.Run(&info)
	assert.Nil(t, err)
	assert.Equal(t, newfname, info.Target)

	info.Src = info.Target

	tsk = NewRemoveTask(Identity)
	err = tsk.Run(&info)
	assert.Nil(t, err)

	info.Src = info.Target

	tsk = NewRemoveTask(func(f string) string { return fname })
	err = tsk.Run(&info)
	assert.Nil(t, err)

	info.Src = info.Target

	tsk = NewRemoveTask(func(f string) string { return newPath })
	err = tsk.Run(&info)
	assert.Nil(t, err)
}
