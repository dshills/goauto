// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package shelltask

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/dshills/goauto"
)

type shellTask struct {
	cmd       string
	args      []string
	transform goauto.Transformer
}

// NewShellTask returns a goauto.Tasker which calls out to the shell with the command cmd and arguments arg
// the goauto.TaskInfo.Target will be the last argument and is equal to goauto.TaskInfo.Src
// goauto.TaskInfo.Target is set to goauto.TaskInfo.Src
func NewShellTask(cmd string, args ...string) goauto.Tasker {
	return &shellTask{cmd: cmd, args: args, transform: goauto.Identity}
}

// NewShellTaskT returns a goauto.Tasker which calls out to the shell with the command cmd and arguments arg
// the goauto.TaskInfo.Target will be the last argument and is generated by transform(goauto.TaskInfo.Src)
// goauto.TaskInfo.Target is set to transform(goauto.TaskInfo.Src)
func NewShellTaskT(transform goauto.Transformer, cmd string, args ...string) goauto.Tasker {
	return &shellTask{cmd: cmd, args: args, transform: transform}
}

// Run will execute the task
func (st *shellTask) Run(info *goauto.TaskInfo) (err error) {
	t0 := time.Now()
	info.Target = st.transform(info.Src)
	info.Buf.Reset()
	targs := append(st.args, info.Target)
	cmd := exec.Command(st.cmd, targs...)
	cmd.Stdout = &info.Buf
	cmd.Stderr = info.Terr

	defer func() {
		fmt.Fprint(info.Tout, info.Buf)
		if err != nil && info.Verbose {
			t1 := time.Now()
			fmt.Fprintf(info.Tout, ">>> %v %v %v\n", st.cmd, st.args, t1.Sub(t0))
		}
	}()
	err = cmd.Run()
	return
}

// NewCatTask returns a goauto.Tasker which writes the file contents to goauto.TaskInfo.Buf and goauto.TaskInfo.Tout
// after running transform(goauto.TaskInfo.Src)
// goauto.TaskInfo.Target is set to transform(goauto.TaskInfo.Src)
func NewCatTask(t goauto.Transformer) goauto.Tasker {
	return goauto.NewTask(t, cat)
}

func cat(info *goauto.TaskInfo) (err error) {
	in, err := os.Open(info.Target)
	if err != nil {
		return
	}
	defer in.Close()
	info.Buf.Reset()
	_, err = info.Buf.ReadFrom(in)
	if err != nil {
		return
	}
	fmt.Fprint(info.Tout, info.Buf)
	if err != nil && info.Verbose {
		fmt.Fprintf(info.Tout, ">>> cat %v", info.Target)
	}
	return
}

// NewRemoveTask returns a goauto.Tasker which will delete the file named transform(goauto.TaskInfo.Src)
// goauto.TaskInfo.Target is set to transform(goauto.TaskInfo.Src)
func NewRemoveTask(t goauto.Transformer) goauto.Tasker {
	return goauto.NewTask(t, remove)
}

func remove(info *goauto.TaskInfo) (err error) {
	err = os.Remove(info.Target)
	if err != nil && info.Verbose {
		fmt.Fprintf(info.Tout, ">>> Remove %v", info.Target)
	}
	return
}

// NewMoveTask returns a goauto.Tasker which will rename a file from goauto.TaskInfo.Target to transform(goauto.TaskInfo.Src)
// goauto.TaskInfo.Target is set to transform(goauto.TaskInfo.Src)
func NewMoveTask(t goauto.Transformer) goauto.Tasker {
	return goauto.NewTask(t, move)
}

func move(info *goauto.TaskInfo) (err error) {
	err = os.Rename(info.Src, info.Target)
	if err != nil && info.Verbose {
		fmt.Fprintf(info.Tout, ">>> Renaming %v to %v", info.Src, info.Target)
	}
	return
}

// NewMkdirTask returns a goauto.Tasker which makes a new dir named transform(goauto.TaskInfo.Src)
// goauto.TaskInfo.Target is not reset
func NewMkdirTask(t goauto.Transformer) goauto.Tasker {
	return goauto.NewTask(t, mkdir)
}

func mkdir(info *goauto.TaskInfo) (err error) {
	dir := info.Target
	info.Target = info.Src
	if err = os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
		return
	}
	if err != nil && info.Verbose {
		fmt.Fprintf(info.Tout, ">>> mkdir %v\n", dir)
	}
	return
}

// NewCopyTask returns a goauto.Tasker that copies the file contents of goauto.TaskInfo.Src to transform(goauto.TaskInfo.Src)
// goauto.TaskInfo.Target is set to transform(goauto.TaskInfo.Src)
func NewCopyTask(t goauto.Transformer) goauto.Tasker {
	return goauto.NewTask(t, fcopy)
}

func fcopy(info *goauto.TaskInfo) (err error) {
	in, err := os.Open(info.Src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(info.Target)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
		if err != nil && info.Verbose {
			fmt.Fprintf(info.Tout, ">>> Copy %v to %v\n", info.Src, info.Target)
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}

	return out.Sync()
}