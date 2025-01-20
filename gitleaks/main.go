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
	execArgs := []string{
		"gitleaks",
		"detect",
	}

	if len(exitCode) > 0 {
		execArgs = append(execArgs, "--exit-code", exitCode)
	}
	if len(reportFormat) > 0 {
		execArgs = append(execArgs, "--report-format", reportFormat)
	}
	if verbose {
		execArgs = append(execArgs, "--verbose")
	}

	return dag.Container().
		From("zricethezav/gitleaks:latest").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(execArgs).
		Stdout(ctx)
}
