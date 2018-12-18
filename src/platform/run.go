package platform

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// OSOptions supplies all possible OS options
var OSOptions = []string{"android", "darwin", "dragonfly", "freebsd", "linux", "nacl", "netbsd", "openbsd", "plan9", "solaris", "windows"}

// EnvVar is used to pass in variables to the runner
type EnvVar struct {
	Key   string
	Value string
}

// RunCommand lets you run a command in your platform's shell
func RunCommand(command string, environment []EnvVar, silent bool) {
	// var stdoutBuf, stderrBuf bytes.Buffer
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

	// stdoutIn, _ := cmd.StdoutPipe()
	// stderrIn, _ := cmd.StderrPipe()

	// var errStdout, errStderr error
	// stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	// stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	cmd.Start()
	// if err != nil {
	// 	log.Fatalf("cmd.Start() failed with '%s'\n", err)
	// }

	// go func() {
	// 	_, errStdout = io.Copy(stdout, stdoutIn)
	// }()

	// go func() {
	// 	_, errStderr = io.Copy(stderr, stderrIn)
	// }()

	_ = cmd.Wait()
}
