package main

import (
	"context"

	"dagger/gitleaks/internal/dagger"
)

type Gitleaks struct{}

func (m *Gitleaks) Detect(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	exitCode string,
	// +optional
	reportFormat string,
	// +optional
	verbose bool,
) (string, error) {
	containerImage := "zricethezav/gitleaks:latest"

	args := []string{
		"gitleaks",
		"detect",
	}

	if len(exitCode) > 0 {
		args = append(args, "--exit-code", exitCode)
	}
	if len(reportFormat) > 0 {
		args = append(args, "--report-format", reportFormat)
	}
	if verbose {
		args = append(args, "--verbose")
	}

	return dag.Container().
		From(containerImage).
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(args).
		Stdout(ctx)
}
