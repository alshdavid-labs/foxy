package platform

import (
	"runtime"

	"github.com/joho/godotenv"
)

// FoxyConfig is the loaded YAML file
type FoxyConfig map[string]FoxyTask

// FoxyTask is what makes up the yaml files
type FoxyTask struct {
	Set      bool
	Parallel bool              `yaml:"parallel"`
	Default  bool              `yaml:"default"`
	Platform string            `yaml:"platform"`
	Steps    []string          `yaml:"steps,flow"`
	Env      map[string]string `yaml:"env,flow"`
	EnvFile  string            `yaml:"env_file"`
}

// CastEnvFile takes a location of an env file, loads it
// then casts the resulting map into application format
func CastEnvFile(envFileLocation string) []EnvVar {
	var environment []EnvVar
	if envFileLocation == "" {
		return environment
	}
	env, _ := godotenv.Read(envFileLocation)
	environment = append(CastEnvArguments(env))
	return environment
}

// CastEnvArguments takes a map of env variables and casts them to
// application format
func CastEnvArguments(envArguments map[string]string) []EnvVar {
	var environment []EnvVar
	for key, value := range envArguments {
		environment = append(environment, EnvVar{
			Key:   key,
			Value: value,
		})
	}
	return environment
}

// FindDefault takes a config
func FindDefault(allTasks FoxyConfig) FoxyTask {
	var defaultTask FoxyTask
	for key, foxyEvent := range allTasks {
		if foxyEvent.Default == true {
			defaultTask = FindTask(key, allTasks)
			break
		}
	}
	return defaultTask
}

// FindTask will return the selected task for the platform
// and taskname
func FindTask(taskName string, allTasks FoxyConfig) FoxyTask {
	for taskKey, task := range allTasks {
		if taskKey == taskName+"("+runtime.GOOS+")" {
			return task
		}
	}
	return allTasks[taskName]
}
