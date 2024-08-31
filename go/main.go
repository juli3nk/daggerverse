package main

import "dagger/go/internal/dagger"

type Go struct {
	Version  string
	Worktree *dagger.Directory
}

func New(
	// Go version
	version string,
	source *dagger.Directory,
) *Go {
	return &Go{Version: version, Worktree: source}
}
