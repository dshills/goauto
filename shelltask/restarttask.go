// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package shelltask

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/dshills/goauto"
)

// A RestartTask represents a task to launch or relaunch an executable file
type RestartTask struct {
	ecmd *exec.Cmd
	Cmd  string
	Args []string
}

// NewRestartTask returns a ReloadTask
func NewRestartTask(cmd string, args ...string) *RestartTask {
	return &RestartTask{Cmd: cmd, Args: args}
}

// Restart will launch or relaunch the application
func (r *RestartTask) Restart(t *goauto.TaskInfo) (err error) {
	if r.Cmd == "" {
		return errors.New("Cmd not set, Nothing to run")
	}

	err = r.Kill(t)
	if err != nil {
		return
	}

	r.ecmd = exec.Command(r.Cmd, r.Args...)
	r.ecmd.Stdout = t.Tout
	r.ecmd.Stderr = t.Terr

	err = r.ecmd.Start()
	go r.ecmd.Wait()
	if t.Verbose {
		fmt.Fprintf(t.Tout, "Process %v started\n", r.Cmd)
	}
	return
}

// Kill will stop the running task
func (r *RestartTask) Kill(t *goauto.TaskInfo) (err error) {
	defer func() { r.ecmd = nil }()
	if r.ecmd == nil || r.ecmd.Process == nil {
		return
	}
	if r.ecmd.ProcessState != nil && r.ecmd.ProcessState.Exited() {
		if t.Verbose {
			fmt.Fprintf(t.Tout, "Process %v already exited\n", r.Cmd)
		}
		return
	}

	if err = r.ecmd.Process.Kill(); err != nil {
		return
	}

	err = r.ecmd.Process.Release()
	if t.Verbose {
		fmt.Fprintf(t.Tout, "Process %v killed\n", r.Cmd)
	}
	return
}

// Run will restart the application in Cmd
func (r *RestartTask) Run(t *goauto.TaskInfo) (err error) {
	return r.Restart(t)
}
