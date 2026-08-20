package main

import "context"

// +check
func (m *Goreleaser) Check(
	ctx context.Context,
	// +optional
	config string,
	// +optional
	// +default=false
	quietly bool,
	// +optional
	// +default=false
	verboseOutput bool,
) (string, error) {
	var execArgs []string

	if config != "" {
		execArgs = append(execArgs, config)
	}

	if quietly {
		execArgs = append(execArgs, "--quiet")
	}

	if verboseOutput {
		execArgs = append(execArgs, "--verbose")
	}

	return m.run(ctx, execArgs)
}
