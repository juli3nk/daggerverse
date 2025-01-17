package main

import (
	"context"
)

// Black for code formatting.
// isort for import sorting.
var toolsFormatter = []string{
	"black",
	"isort",
}

func (m *Python) Fmt(
	ctx context.Context,
	tool string,
	// +optional
	// +default="."
	source string,
) (string, error) {
	var execArgs []string

	if err := validateTool(toolsFormatter, tool); err != nil {
		return "", err
	}

	checkParam := "--check"
	if tool == "isort" {
		checkParam = "--check-only"
	}

	execArgs = append(execArgs, checkParam, source)

	return m.container().
		WithExec(append([]string{tool}, execArgs...)).
		Stdout(ctx)
}
