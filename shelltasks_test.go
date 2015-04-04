// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShellTask(t *testing.T) {
	info := TaskInfo{Src: "FROGS", Tout: ioutil.Discard, Terr: ioutil.Discard}

	tsk := NewShellTask("echo", "-n")
	err := tsk.Run(&info)
	assert.Nil(t, err)
	assert.Equal(t, "FROGS", info.Buf.String())

	tsk = NewShellTaskT(Identity, "echo", "-n")
	err = tsk.Run(&info)
	assert.Nil(t, err)
	assert.Equal(t, "FROGS", info.Buf.String())
}

func TestOSTasks(t *testing.T) {
	t0 := time.Now()

	tp := filepath.Join("src", "github.com", "dshills", "goauto", "testing")
	path, err := AbsPath(tp)
	assert.Nil(t, err)
	fname := filepath.Join(path, "testdata")
	assert.Nil(t, err)
	f, err := os.Create(fname)
	defer f.Close()
	f.WriteString("TESTING")
	_, err = os.Stat(fname)
	assert.Nil(t, err)

	info := TaskInfo{Src: fname, Tout: ioutil.Discard, Terr: ioutil.Discard}

	newPath := filepath.Join(path, "t")
	tsk := NewMkdirTask(func(f string) string { return newPath })
	err = tsk.Run(&info)
	assert.Nil(t, err)
	_, err = os.Stat(newPath)
	assert.Nil(t, err)

	info.Src = info.Target

	newfname := filepath.Join(newPath, "testdata2")
	tsk = NewCopyTask(func(string) string { return newfname })
	err = tsk.Run(&info)
	assert.Nil(t, err)
	assert.Equal(t, newfname, info.Target)
	_, err = os.Stat(newfname)
	assert.Nil(t, err)

	info.Src = info.Target

	tsk = NewRemoveTask(Identity)
	err = tsk.Run(&info)
	assert.Nil(t, err)
	_, err = os.Stat(newfname)
	assert.NotNil(t, err)

	info.Src = info.Target

	tsk = NewRemoveTask(func(f string) string { return fname })
	err = tsk.Run(&info)
	assert.Nil(t, err)
	_, err = os.Stat(fname)
	assert.NotNil(t, err)

	info.Src = info.Target

	tsk = NewRemoveTask(func(f string) string { return newPath })
	err = tsk.Run(&info)
	assert.Nil(t, err)
	_, err = os.Stat(newPath)
	assert.NotNil(t, err)

	t1 := time.Now()
	log.Printf("TestOSTasks finished in %v", t1.Sub(t0))
}
