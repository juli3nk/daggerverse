package main

import "context"

// +check
func (m *Goreleaser) Check(
	ctx context.Context,
	// +optional
	config string,
	// +optional
	// +default=false
	quiet bool,
	// +optional
	// +default=false
	verbose bool,
) (string, error) {
	var execArgs []string

	if config != "" {
		execArgs = append(execArgs, config)
	}

	if quiet {
		execArgs = append(execArgs, "--quiet")
	}

	if verbose {
		execArgs = append(execArgs, "--verbose")
	}

	return m.run(ctx, execArgs)
}
