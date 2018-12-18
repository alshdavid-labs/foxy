package main

import (
	"flag"
	"fmt"
	"foxy/platform"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
)

var acceptedFileNames = []string{
	"foxy.yml",
	"FOXY.yml",
	"foxyfile.yml",
	"FOXYFILE.yml",
	"foxyfile",
	"FOXYFILE",
	"foxy",
	"FOXY",
	"taskfile",
	"TASKFILE",
	"taskfile.yml",
	"TASKFILE.yml",
}

func tryGetFoxyFile(defaults []string) string {
	var found string
	for _, pathname := range defaults {
		if _, err := os.Stat(pathname); !os.IsNotExist(err) {
			found = pathname
			break
		}
	}
	return found
}

func main() {
	var chosenTask platform.FoxyTask
	foxyFilePath := tryGetFoxyFile(acceptedFileNames)

	flag.StringVar(&foxyFilePath, "file", foxyFilePath, "a string")
	flag.StringVar(&foxyFilePath, "f", foxyFilePath, "a string")
	flag.Parse()

	foxyYamlBinary, _ := ioutil.ReadFile(foxyFilePath)
	taskData := platform.LoadYaml(string(foxyYamlBinary))

	if len(os.Args) != 1 {
		chosenTask = platform.FindTask(os.Args[len(os.Args)-1], taskData)
	}

	if chosenTask.Set == false {
		chosenTask = platform.FindDefault(taskData)
	}

	if chosenTask.Set == false {
		fmt.Println("No task selected")
		os.Exit(1)
	}

	var environment []platform.EnvVar
	environment = append(environment, platform.CastEnvFile(chosenTask.EnvFile)...)
	environment = append(environment, platform.CastEnvArguments(chosenTask.Env)...)

	// Run steps
	if chosenTask.Parallel == false {
		for _, step := range chosenTask.Steps {
			cmd := platform.GenerateCommand(step, environment)
			platform.RunCommand(cmd, chosenTask.Silent)
		}
	} else {
		var cmds []*exec.Cmd
		var wg sync.WaitGroup

		defer func() {
			fmt.Println("Process clean up")
			for _, cmd := range cmds {
				if cmd != nil {
					cmd.Process.Kill()
				}
			}
		}()

		if chosenTask.AbortOnExit == true {
			wg.Add(1)
		} else {
			wg.Add(len(chosenTask.Steps))
		}
		for _, step := range chosenTask.Steps {
			cmd := platform.GenerateCommand(step, environment)
			cmds = append(cmds, cmd)

			go func(cmd *exec.Cmd, step string) {
				defer wg.Done()
				platform.RunCommand(cmd, chosenTask.Silent)
				if chosenTask.AbortOnExit == true {
					fmt.Println("Exit triggered by", step)
				}
			}(cmd, step)
		}
		wg.Wait()
	}
}
