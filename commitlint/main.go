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
	// A list of arguments to pass to commitlint.
	args []string,
) error {
	if len(args) == 0 {
		args = append(args, "--last")
	}

	_, err := dag.Container().
		From("commitlint/commitlint:latest").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(args, dagger.ContainerWithExecOpts{UseEntrypoint: true}).
		Sync(ctx)

	return err
}
