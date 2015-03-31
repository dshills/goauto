// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

type sassTask struct {
	cssDir string
	args   []string
}

// NewSassTask returns a Task that will run command line sass in the directory of the file change
// sass must be in the PATH
// Blank strings for cssDir, cachDir or style will use sass defaults
// TaskInfo.Target will not be updated
func NewSassTask(cssDir, cacheDir, style string) Tasker {
	st := sassTask{cssDir: cssDir}
	if cacheDir != "" {
		st.args = append(st.args, "--cache-location", cacheDir)
	}
	if style != "" {
		st.args = append(st.args, "--style", style)
	}
	return st
}

func (st sassTask) Run(info *TaskInfo) (err error) {
	dir := filepath.Dir(info.Src)
	info.Buf.Reset()
	fmt.Fprintln(info.Tout, "Sass ...", dir)
	if st.cssDir != "" {
		dir += ":" + st.cssDir
	}
	targs := append(st.args, "--update", dir)
	fmt.Fprintln(info.Tout, targs)
	cmd := exec.Command("sass", targs...)
	cmd.Stdout = &info.Buf
	cmd.Stderr = info.Terr

	defer func() {
		info.Tout.Write(info.Buf.Bytes())
	}()
	return cmd.Run()
}
