package main

import (
	"context"

	"dagger/commitlint/internal/dagger"
)

type Commitlint struct{}

// Lint runs commitlint to lint commit messages.
// +check
func (m *Commitlint) Lint(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
	// +optional
	config string,
	// +optional
	// +default=false
	defaultConfig bool,
	// +optional
	from string,
	// +optional
	// +default=false
	fromLastTag bool,
	// +optional
	// +default=false
	last bool,
	// +optional
	// +default=false
	quiet bool,
	// +optional
	to string,
	// +optional
	// +default=false
	strict bool,
) error {
	var execArgs []string

	if config != "" {
		execArgs = append(execArgs, "--config", config)
	}

	if defaultConfig && config == "" {
		execArgs = append(execArgs, "--default-config")
	}

	if from != "" {
		execArgs = append(execArgs, "--from", from)
	}

	if fromLastTag && from == "" {
		execArgs = append(execArgs, "--from-last-tag")
	}

	if last && from == "" && !fromLastTag {
		execArgs = append(execArgs, "--last")
	}

	if quiet {
		execArgs = append(execArgs, "--quiet")
	}

	if to != "" {
		execArgs = append(execArgs, "--to", to)
	}

	if strict {
		execArgs = append(execArgs, "--strict")
	}

	_, err := dag.Container().
		From("commitlint/commitlint:latest").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(execArgs, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		Sync(ctx)

	return err
}
