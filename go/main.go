package main

import "dagger/go/internal/dagger"

type Go struct {
	Version  string
	Worktree *dagger.Directory
}

func New(
	// Go version
	// +default="1.26.1"
	version string,
	// +defaultPath="."
	source *dagger.Directory,
) *Go {
	return &Go{Version: version, Worktree: source}
}
