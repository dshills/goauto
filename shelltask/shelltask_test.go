// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package shelltask

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/dshills/goauto"
)

func TestShellTask(t *testing.T) {
	e := "FROGS"

	info := goauto.TaskInfo{Src: e, Tout: ioutil.Discard, Terr: ioutil.Discard}

	tsk := NewShellTask("echo", "-n")
	err := tsk.Run(&info)
	if err != nil {
		t.Error(err)
	}
	if info.Buf.String() != e {
		t.Errorf("Expected %v got %v\n", e, info.Buf.String())
	}

	tsk = NewShellTaskT(goauto.Identity, "echo", "-n")
	err = tsk.Run(&info)
	if err != nil {
		t.Error(err)
	}
	if info.Buf.String() != e {
		t.Errorf("Expected %v got %v\n", e, info.Buf.String())
	}
}

func TestOSTasks(t *testing.T) {
	tp := filepath.Join("src", "github.com", "dshills", "goauto", "testing")
	path, err := goauto.AbsPath(tp)
	if err != nil {
		t.Error(err)
	}
	fname := filepath.Join(path, "testdata")
	f, err := os.Create(fname)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	f.WriteString("TESTING")
	_, err = os.Stat(fname)
	if err != nil {
		t.Fatal(err)
	}

	info := goauto.TaskInfo{Src: fname, Tout: ioutil.Discard, Terr: ioutil.Discard}

	newPath := filepath.Join(path, "t")
	tsk := NewMkdirTask(func(f string) string { return newPath })
	err = tsk.Run(&info)
	if err != nil {
		t.Error(err)
	}
	_, err = os.Stat(newPath)
	if err != nil {
		t.Error(err)
	}

	info.Src = info.Target

	newfname := filepath.Join(newPath, "testdata2")
	tsk = NewCopyTask(func(string) string { return newfname })
	err = tsk.Run(&info)
	if err != nil {
		t.Error(err)
	}
	if newfname != info.Target {
		t.Errorf("Expected %v got %v\n", newfname, info.Target)
	}
	_, err = os.Stat(newfname)
	if err != nil {
		t.Error(err)
	}

	info.Src = info.Target

	tsk = NewRemoveTask(goauto.Identity)
	err = tsk.Run(&info)
	if err != nil {
		t.Error(err)
	}
	_, err = os.Stat(newfname)
	if err == nil {
		t.Errorf("Expected error, file should be gone\n")
	}

	info.Src = info.Target

	tsk = NewRemoveTask(func(f string) string { return fname })
	err = tsk.Run(&info)
	if err != nil {
		t.Error(err)
	}
	_, err = os.Stat(fname)
	if err == nil {
		t.Errorf("Expected error, file should be gone\n")
	}

	info.Src = info.Target

	tsk = NewRemoveTask(func(f string) string { return newPath })
	err = tsk.Run(&info)
	if err != nil {
		t.Error(err)
	}
	_, err = os.Stat(newPath)
	if err == nil {
		t.Errorf("Expected error, file should be gone\n")
	}
}
