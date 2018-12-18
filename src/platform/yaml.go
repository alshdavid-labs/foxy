package platform

import (
	yaml "gopkg.in/yaml.v2"
)

// LoadYaml will take a yaml string and parse it into a FoxyEvent
func LoadYaml(data string) map[string]FoxyTask {
	m := make(map[string]FoxyTask)
	yaml.Unmarshal([]byte(data), &m)

	for key, foxyEvent := range m {
		updated := foxyEvent
		updated.Set = true
		m[key] = updated
	}

	return m
}
