package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

const (
	licensesConfigImage = "ghcr.io/juli3nk/go-licenses-config:latest"
)

// Licenses scans the Go module and its dependencies using go-licenses.
// Fails if a forbidden, incompatible, or unrecognized license is detected.
// +check
func (m *Go) Licenses(
	ctx context.Context,
	// +optional
	// +default="./..."
	packages string,
) error {
	// Generate args from the config image
	cfgCtr := dag.Container().
		From(licensesConfigImage).
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithExec([]string{"go-licenses-config", "-packages", packages})

	stdout, err := cfgCtr.Stdout(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate license config: %w", err)
	}

	checkArgs, err := decodeLicensesConfigCheckArgs(stdout)
	if err != nil {
		return fmt.Errorf("failed to decode license args: %w", err)
	}

	execArgs := []string{"go-licenses"}

	if len(checkArgs) > 0 {
		// Ensure the subcommand is "check"
		if checkArgs[0] != "check" {
			return fmt.Errorf("expected subcommand 'check', got %q", checkArgs[0])
		}

		// Run go-licenses check
		execArgs = append(execArgs, checkArgs...)
	}

	_, err = dag.Container().
		From(fmt.Sprintf("golang:%s", m.Version)).
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume(fmt.Sprintf("go-mod-%s", m.Version))).
		WithMountedCache("/go/bin", dag.CacheVolume(fmt.Sprintf("go-bin-%s", m.Version))).
		WithExec([]string{"go", "install", "github.com/google/go-licenses/v2@latest"}).
		WithExec(execArgs).
		Sync(ctx)

	return err
}

func decodeLicensesConfigCheckArgs(b64 string) ([]string, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}

	var args []string
	if err := json.Unmarshal(data, &args); err != nil {
		return nil, err
	}
	return args, nil
}
