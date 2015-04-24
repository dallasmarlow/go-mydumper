package mydumper

import (
	"bytes"
	"errors"
	"log"
	"os/exec"
	"time"
)

var (
	ErrExecTimeout            = errors.New(`Bash execution timeout, child terminated`)
	ErrExecTimeoutTermFailure = errors.New(`Bash execution timeout, child termination failed`)
)

func bashExec(bashCmd string, stdout, stderr *bytes.Buffer) (*exec.Cmd, error) {
	cmd := exec.Command(`/bin/bash`, `-c`, bashCmd)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return cmd, nil
}

func ExecWithTimeout(bashCmd string, timeout time.Duration) (string, string, error) {
	var stdout, stderr bytes.Buffer

	cmd, err := bashExec(bashCmd, &stdout, &stderr)
	if err != nil {
		return stdout.String(), stderr.String(), err
	}

	errs := make(chan error)
	defer close(errs)
	go func() {
		errs <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		if err := cmd.Process.Kill(); err != nil {
			log.Println(`Bash exec timeout, child term failure: ` + err.Error())

			return &stdout, &stderr, ErrExecTimeoutTermFailure
		}

		return &stdout, &stderr, ErrExecTimeout
	case err := <-errs:
		return &stdout, &stderr, err
	}
}
