package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

type FoxyEvent struct {
	Context string   `yaml:"context"`
	Watch   string   `yaml:"watch"`
	Steps   []string `yaml:"steps,flow"`
	Env     []string `yaml:"env,flow"`
	EnvFile string   `yaml:"string"`
}

func loadYaml(data string) map[string]FoxyEvent {
	m := make(map[string]FoxyEvent)
	yaml.Unmarshal([]byte(data), &m)
	return m
}

func runCommand(command string) {
	var stdoutBuf, stderrBuf bytes.Buffer

	cmd := exec.Command("/bin/bash", "-c", command)

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	_ = cmd.Wait()
}

func main() {
	var foxyFilePath string
	flag.StringVar(&foxyFilePath, "file", "./foxyfile.yml", "a string")
	flag.Parse()
	foxyYamlBinary, _ := ioutil.ReadFile(foxyFilePath)
	foxyYaml := string(foxyYamlBinary)
	foxyData := loadYaml(foxyYaml)
	foxyChosenEvent := os.Args[len(os.Args)-1]

	// fmt.Println(foxyFilePath)
	// fmt.Println(foxyYaml)
	// fmt.Println(foxyChosenEvent)

	// fmt.Println(foxyData[foxyChosenEvent])

	for _, step := range foxyData[foxyChosenEvent].Steps {
		runCommand(step)
	}
}
