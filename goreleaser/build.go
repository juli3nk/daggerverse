package main

import (
	"context"

	"dagger/goreleaser/internal/dagger"
)

func (m *Goreleaser) Build(
	ctx context.Context,
	// Automatically sets --snapshot if the repository is dirty
	// +optional
	// +default=false
	autoSnapshot bool,
	// Removes the 'dist' directory before building
	// +optional
	// +default=false
	clean bool,
	// Load configuration from file
	// +optional
	config string,
	// Builds only the specified build ids
	// +optional
	id string,
	// Copy the binary to the path after the build.
	// Only taken into account when using --single-target
	// and a single id (either with --id or if configuration only has one build)
	// +optional
	outputPath string,
	// Number of tasks to run concurrently (default: number of CPUs)
	// +optional
	parallelism string,
	// Builds only for current GOOS and GOARCH, regardless of what's set in the configuration file
	// +optional
	// +default=false
	singleTarget bool,
	// Skip the given options (valid options are: before, post-hooks, pre-hooks, validate)
	// +optional
	skip string,
	// Generate an unversioned snapshot build, skipping all validations
	// +optional
	// +default=false
	snapshot bool,
	// Timeout to the entire build process
	// +optional
	// +default="1h0m0s"
	timeout string,
	// Enable verbose mode
	// +optional
	// +default=false
	verboseOutput bool,
) *dagger.Container {
	execArgs := []string{"build"}

	if autoSnapshot {
		execArgs = append(execArgs, "--auto-snapshot")
	}

	if clean {
		execArgs = append(execArgs, "--clean")
	}

	if config != "" {
		execArgs = append(execArgs, "--config", config)
	}

	if id != "" {
		execArgs = append(execArgs, "--id", id)
	}

	if outputPath != "" {
		execArgs = append(execArgs, "--output", outputPath)
	}

	if parallelism != "" {
		execArgs = append(execArgs, "--parallelism", parallelism)
	}

	if singleTarget {
		execArgs = append(execArgs, "--single-target")
	}

	if skip != "" {
		execArgs = append(execArgs, "--skip", skip)
	}

	if snapshot {
		execArgs = append(execArgs, "--snapshot")
	}

	if timeout != "" {
		execArgs = append(execArgs, "--timeout", timeout)
	}

	if verboseOutput {
		execArgs = append(execArgs, "--verbose")
	}

	return m.run(execArgs)
}
