package main

import (
	"flag"
	"fmt"
	"foxy/platform"
	"io/ioutil"
	"os"
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
			platform.RunCommand(step, environment, false)
		}
	} else {
		var wg sync.WaitGroup
		wg.Add(len(chosenTask.Steps))
		for _, step := range chosenTask.Steps {
			go func(step string) {
				defer wg.Done()
				platform.RunCommand(step, environment, false)
			}(step)
		}
		wg.Wait()
	}
}
