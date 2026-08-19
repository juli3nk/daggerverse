package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type ModuleUpdate struct {
	Path    string `json:"Path"`
	Version string `json:"Version"`
	Update  struct {
		Version string `json:"Version"`
	} `json:"Update,omitempty"`
}

// ScanDependencies checks for outdated dependencies
// +check
func (m *Go) AuditDependencies(ctx context.Context) (string, error) {
	// 1. Exécution avec cache
	output, err := dag.Container().
		From(fmt.Sprintf("golang:%s", m.Version)).
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume(fmt.Sprintf("go-mod-%s", m.Version))).
		WithExec([]string{"go", "list", "-u", "-m", "-json", "all"}).
		Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("go list failed: %w", err)
	}

	// 2. Parsing line by line (go list -m -json sorting JSON objects separeted by par \n)
	var outdated []string
	decoder := json.NewDecoder(strings.NewReader(output))
	for {
		var mod ModuleUpdate
		if err := decoder.Decode(&mod); err == io.EOF {
			break
		} else if err != nil {
			return "", fmt.Errorf("parsing module json: %w", err)
		}
		if mod.Update.Version != "" {
			outdated = append(outdated, fmt.Sprintf("%s: %s → %s", mod.Path, mod.Version, mod.Update.Version))
		}
	}

	if len(outdated) == 0 {
		return "All dependencies are up to date", nil
	}

	return strings.Join(outdated, "\n"), nil
}
