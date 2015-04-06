// Copyright 2015 Davin Hills. All rights reserved.
// MIT license. License details can be found in the LICENSE file.

package goauto

import (
	"errors"
	"testing"
)

type testTask struct {
	showError bool
}

func (t testTask) Run(i *TaskInfo) (err error) {
	if t.showError {
		return errors.New("This returns an error")
	}
	return
}

func TestTask(t *testing.T) {
	info := TaskInfo{}
	tsk := testTask{}
	err := tsk.Run(&info)
	if err != nil {
		t.Errorf("TestTask returned %v", err)
	}

	tsk.showError = true
	err = tsk.Run(&info)
	if err == nil {
		t.Errorf("TestTask should have returned an error")
	}
}
