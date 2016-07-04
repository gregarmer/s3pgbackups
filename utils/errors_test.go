package utils

import (
	"errors"
	"os"
	"os/exec"
	"testing"
)

func TestCheckErr(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		fakeError := errors.New("Fake Error")
		CheckErr(fakeError)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCheckErr")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestCheckErrWithoutError(t *testing.T) {
	CheckErr(nil)
}
