package main

import (
	"context"

	"dagger/gitleaks/internal/dagger"
)

type Gitleaks struct{}

// Detect scans for secrets
// +check
func (m *Gitleaks) Detect(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
	// +optional
	exitCode string,
	// +optional
	reportFormat string,
	// +optional
	verboseOutput bool,
) error {
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
	if verboseOutput {
		execArgs = append(execArgs, "--verbose")
	}

	_, err := dag.Container().
		From("ghcr.io/gitleaks/gitleaks:latest").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(execArgs).
		Sync(ctx)

	return err
}
