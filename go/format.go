package main

import (
	"bufio"
	"context"
	"strings"
)

func (m *Go) Fmt() ([]string, error) {
	var result []string

	out, err := dag.Container().
		From("golang:latest").
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec([]string{"gofmt", "-l"}).
		Stdout(context.TODO())
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	return result, nil
}
