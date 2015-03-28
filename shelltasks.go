// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"io"
	"os"
	"os/exec"
)

// NewShellTask returns a new task that will execute the named program with the given arguments
// the Target of the task will be the last argument to the shell command
// FileName will not be changed
func NewShellTask(cmd string, args ...string) *Task {
	return NewTaskType(func(t *Task, wout, werr io.Writer) error {
		t.Buffer.Reset()
		targs := append(args, t.Target)
		cmd := exec.Command(cmd, targs...)
		cmd.Stdout = &t.Buffer
		cmd.Stderr = werr
		defer func() {
			t.Buffer.WriteTo(wout)
		}()
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
	})
}

// NewCatTask returns a Task that will write the contents of a file to Buffer and wout
func NewCatTask() *Task {
	return NewTaskType(func(t *Task, wout, werr io.Writer) (err error) {
		in, err := os.Open(t.FileName)
		if err != nil {
			return
		}
		defer in.Close()
		t.Buffer.Reset()
		_, err = t.Buffer.ReadFrom(in)
		if err != nil {
			return
		}
		t.Buffer.WriteTo(wout)
		return nil
	})
}

// NewRemoveTask returns a Task that will remove a file
// FileName will be unchanged
func NewRemoveTask() *Task {
	return NewTaskType(func(t *Task, wout, werr io.Writer) error {
		err := os.Remove(t.Target)
		if err != nil {
			return err
		}
		return nil
	})
}

// NewRenameTask returns a Task that will rename a file using the t.TargetFunc
// FileName will be the new file name
func NewRenameTask() *Task {
	return NewTaskType(func(t *Task, wout, werr io.Writer) error {
		err := os.Rename(t.FileName, t.Target)
		if err != nil {
			return err
		}
		t.FileName = t.Target
		return nil
	})
}

// NewMkDirTask returns a Task that will create a new directory using the t.TargetFunc
// FileName will be unchanged
func NewMkDirTask() *Task {
	return NewTaskType(func(t *Task, wout, werr io.Writer) error {
		err := os.Mkdir(t.Target, 0755)
		if err != nil {
			return err
		}
		return nil
	})
}

// NewCopyTask returns a Task that will copy a t.FileName to t.Target
// FileName will be Target
func NewCopyTask() *Task {
	// This is a modified version from markc on StackOverflow
	return NewTaskType(func(t *Task, wout, werr io.Writer) (err error) {
		in, err := os.Open(t.FileName)
		if err != nil {
			return
		}
		defer in.Close()
		out, err := os.Create(t.Target)
		if err != nil {
			return
		}
		defer func() {
			cerr := out.Close()
			if err == nil {
				err = cerr
			}
		}()
		if _, err = io.Copy(out, in); err != nil {
			return
		}
		err = out.Sync()
		t.FileName = t.Target
		return
	})
}
