package platform

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

// OSOptions supplies all possible OS options
var OSOptions = []string{"android", "darwin", "dragonfly", "freebsd", "linux", "nacl", "netbsd", "openbsd", "plan9", "solaris", "windows"}

// EnvVar is used to pass in variables to the runner
type EnvVar struct {
	Key   string
	Value string
}

func GenerateCommand(command string, environment []EnvVar) *exec.Cmd {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		command = strings.Replace(command, "&&", ";", -1)
		cmd = exec.Command("PowerShell.exe", "-Command", command)
	} else {
		cmd = exec.Command("/bin/bash", "-c", command)
	}
	var environments []string
	for _, env := range environment {
		environments = append(environments, fmt.Sprintf("%s=%s", env.Key, env.Value))
	}
	cmd.Env = append(os.Environ(), environments...)
	return cmd
}

// RunCommand lets you run a command in your platform's shell
func RunCommand(cmd *exec.Cmd, silent bool) {
	var stdoutBuf, stderrBuf bytes.Buffer
	var stdoutIn, stderrIn io.ReadCloser
	var stdout, stderr io.Writer
	var errStdout, errStderr error

	if silent == false {
		stdoutIn, _ = cmd.StdoutPipe()
		stderrIn, _ = cmd.StderrPipe()
		stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
		stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	}
	err := cmd.Start()
	if silent == false {
		if err != nil {
			log.Fatalf("cmd.Start() failed with '%s'\n", err)
		}

		go func() {
			_, errStdout = io.Copy(stdout, stdoutIn)
		}()

		go func() {
			_, errStderr = io.Copy(stderr, stderrIn)
		}()
	}

	cmd.Wait()
}

const defaultFailedCode = 1

// RunCommand lets you run a command in your platform's shell
func RunCommandd(command string, environment []EnvVar, silent bool) (stdout string, stderr string, exitCode int) {
	fmt.Println(command)
	var outbuf, errbuf bytes.Buffer
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		command = strings.Replace(command, "&&", ";", -1)
		cmd = exec.Command("PowerShell.exe", "-Command", command)
	} else {
		cmd = exec.Command("/bin/bash", "-c", command)
	}

	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	var environments []string
	for _, env := range environment {
		environments = append(environments, fmt.Sprintf("%s=%s", env.Key, env.Value))
	}
	cmd.Env = append(os.Environ(), environments...)

	err := cmd.Run()
	stdout = outbuf.String()
	stderr = errbuf.String()

	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			exitCode = defaultFailedCode
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	// log.Printf("command result, stdout: %v, stderr: %v, exitCode: %v", stdout, stderr, exitCode)
	return
}
