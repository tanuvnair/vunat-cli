package projects

import "fmt"

type Project struct {
	Name     string
	Commands [][]string
}

var registry = map[string]Project{
	"example": {
		Name: "example",
		Commands: [][]string{
			{"echo", "hello from example project"},
		},
	},
}

func Get(name string) (Project, error) {
	project, ok := registry[name]
	if !ok {
		return Project{}, fmt.Errorf("unknown project: %s", name)
	}
	return project, nil
}
