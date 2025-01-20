package main

import "context"

// Flake8 for PEP 8 checks.
// Pylint for deeper analysis.
// Mypy for static typing.
// Bandit for security analysis.
var toolsLinter = []string{
	"bandit",
	"flake8",
	"mypy",
	"pylint",
}

func (m *Python) Lint(
	ctx context.Context,
	tool string,
	// +optional
	// +default="."
	source string,
) (string, error) {
	var execArgs []string

	if err := validateTool(toolsLinter, tool); err != nil {
		return "", err
	}

	if tool == "bandit" {
		execArgs = append(execArgs, "-r")
	}

	execArgs = append(execArgs, source)

	return m.container(tool).
		WithExec(append([]string{tool}, execArgs...)).
		Stdout(ctx)
}
