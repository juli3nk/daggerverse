package main

import (
	"fmt"

	"dagger/python/internal/dagger"
)

type Python struct {
	Version  string
	Worktree *dagger.Directory
}

func New(
	// Python version
	version string,
	source *dagger.Directory,
) *Python {
	return &Python{Version: version, Worktree: source}
}

func (m *Python) container(tool string) *dagger.Container {
	pipInstallCommand := []string{"pip", "install", tool}

	return dag.Container().
		From(fmt.Sprintf("python:%s", m.Version)).
		WithExec([]string{"pip", "install", "--upgrade", "pip"}).
		WithExec(pipInstallCommand).
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src")
}

func validateTool(tools []string, name string) error {
	for _, tool := range tools {
		if tool == name {
			return nil
		}
	}

	return fmt.Errorf("the tool '%s' is not supported", name)
}
