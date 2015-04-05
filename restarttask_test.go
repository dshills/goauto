package goauto

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRestartBlocking(t *testing.T) {
	tsk := NewRestartTask("cat") // blocking command
	ti := TaskInfo{Tout: os.Stdout, Terr: os.Stderr, Verbose: Verbose}

	err := tsk.Restart(&ti)
	assert.Nil(t, err)

	err = tsk.Kill(&ti)
	assert.Nil(t, err)

	err = tsk.Restart(&ti)
	assert.Nil(t, err)

	time.Sleep(3 * time.Second)

	err = tsk.Restart(&ti)
	assert.Nil(t, err)

	err = tsk.Restart(&ti)
	assert.Nil(t, err)
}

func TestRestartExited(t *testing.T) {
	tsk := NewRestartTask("echo", "GoAuto!!!") // non blocking command
	ti := TaskInfo{Tout: os.Stdout, Terr: os.Stderr, Verbose: Verbose}

	err := tsk.Restart(&ti)
	assert.Nil(t, err)

	time.Sleep(3 * time.Second)

	err = tsk.Kill(&ti)
	assert.Nil(t, err)
}
