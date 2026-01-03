package main

import (
	"context"
	"strconv"

	"dagger/shfmt/internal/dagger"
)

type Shfmt struct{}

// Format formats shell scripts using shfmt
func (m *Shfmt) Fmt(
	ctx context.Context,
	source *dagger.Directory,

	// Generic flags
	// +optional
	list bool,
	// +optional
	listZero bool,
	// +optional
	write bool,
	// +optional
	diff bool,
	// +optional
	applyIgnore bool,
	// +optional
	filename string,

	// Parser flags
	// +optional
	languageDialect string,
	// +optional
	posix bool,
	// +optional
	simplify bool,

	// Printer flags
	// +optional
	indent int,
	// +optional
	binaryNextLine bool,
	// +optional
	caseIndent bool,
	// +optional
	spaceRedirects bool,
	// +optional
	keepPadding bool,
	// +optional
	funcNextLine bool,
	// +optional
	minify bool,
) (string, error) {
	execArgs := []string{"shfmt"}

	// Generic flags
	if listZero {
		execArgs = append(execArgs, "-l=0")
	} else if list {
		execArgs = append(execArgs, "-l")
	}

	if write {
		execArgs = append(execArgs, "-w")
	}

	if diff {
		execArgs = append(execArgs, "-d")
	}

	if applyIgnore {
		execArgs = append(execArgs, "--apply-ignore")
	}

	if len(filename) > 0 {
		execArgs = append(execArgs, "--filename", filename)
	}

	// Parser flags
	if posix {
		execArgs = append(execArgs, "-p")
	} else if len(languageDialect) > 0 {
		execArgs = append(execArgs, "-ln", languageDialect)
	}

	if simplify {
		execArgs = append(execArgs, "-s")
	}

	// Printer flags
	if indent > 0 {
		execArgs = append(execArgs, "-i", strconv.Itoa(indent))
	} else if indent == 0 {
		// 0 est la valeur par défaut (tabs), mais on peut l'expliciter si besoin
		// execArgs = append(execArgs, "-i", "0")
	}

	if binaryNextLine {
		execArgs = append(execArgs, "-bn")
	}

	if caseIndent {
		execArgs = append(execArgs, "-ci")
	}

	if spaceRedirects {
		execArgs = append(execArgs, "-sr")
	}

	if keepPadding {
		execArgs = append(execArgs, "-kp")
	}

	if funcNextLine {
		execArgs = append(execArgs, "-fn")
	}

	if minify {
		execArgs = append(execArgs, "-mn")
	}

	// Ajouter le répertoire à formater
	execArgs = append(execArgs, ".")

	return dag.Container().
		From("mvdan/shfmt:latest").
		WithMountedDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(execArgs).
		Stdout(ctx)
}
