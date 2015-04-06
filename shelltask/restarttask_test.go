package shelltask

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/dshills/goauto"
)

func TestRestartBlocking(t *testing.T) {
	tsk := NewRestartTask("cat") // blocking command
	ti := goauto.TaskInfo{Tout: ioutil.Discard, Terr: ioutil.Discard, Verbose: goauto.Silent}

	err := tsk.Restart(&ti)
	if err != nil {
		t.Error(err)
	}

	err = tsk.Kill(&ti)
	if err != nil {
		t.Error(err)
	}

	err = tsk.Restart(&ti)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(3 * time.Second)

	err = tsk.Restart(&ti)
	if err != nil {
		t.Error(err)
	}

	err = tsk.Restart(&ti)
	if err != nil {
		t.Error(err)
	}
}

func TestRestartExited(t *testing.T) {
	tsk := NewRestartTask("echo", "GoAuto!!!") // non blocking command
	ti := goauto.TaskInfo{Tout: ioutil.Discard, Terr: ioutil.Discard, Verbose: goauto.Silent}

	err := tsk.Restart(&ti)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(3 * time.Second)

	err = tsk.Kill(&ti)
	if err != nil {
		t.Error(err)
	}
}

func TestRestartWorkflow(t *testing.T) {
	tsk := NewRestartTask("echo", "GoAuto!!!") // non blocking command
	ti := goauto.TaskInfo{Tout: ioutil.Discard, Terr: ioutil.Discard, Verbose: goauto.Verbose}
	wf := goauto.NewWorkflow(tsk)
	wf.Run(&ti)
}
