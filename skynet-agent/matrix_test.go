package main

import (
	"os/exec"
	"strings"
	"testing"
)

func Test_execCmd(t *testing.T) {
	rawCmd := strings.Split("/bin/ls -al", " ")
	cmd := exec.Command(rawCmd[0], rawCmd[1:]...)
	err := cmd.Run()
	if err != nil {
		t.Log(err)
	}
}
