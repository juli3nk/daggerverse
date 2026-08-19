package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"dagger/gitlocal/internal/dagger"

	"github.com/juli3nk/go-utils/filedir"
)

type ChangedFiles struct {
	Files string
}

// ChangedFiles retourne la liste des fichiers modifiés entre baseRef et headRef.
func (m *Gitlocal) ChangedFiles(
	ctx context.Context,
	// +defaultPath="."
	repo *dagger.Directory,
	// +default="origin/main"
	baseRef string,
	// +default="HEAD"
	headRef string,
	// +optional
	ext []string,
) (*ChangedFiles, error) {
	diffRange := fmt.Sprintf("%s...%s", baseRef, headRef)

	// ── 1. Diff entre les deux refs ──
	stdout, stderr, exitCode, err := m.gitRaw(ctx, repo,
		[]string{"diff", "--name-only", diffRange},
	)
	if err != nil {
		return nil, err
	}
	if exitCode != 0 {
		return nil, fmt.Errorf("git diff --name-only %s failed (exit %d): %s",
			diffRange, exitCode, strings.TrimSpace(stderr))
	}

	fileSet := make(map[string]struct{})
	for _, f := range getFiles(stdout) {
		fileSet[f] = struct{}{}
	}

	// ── 2. Changements non commités (worktree + index) ──
	if headRef == "" || headRef == "HEAD" {
		stdout, stderr, exitCode, err = m.gitRaw(ctx, repo,
			[]string{"diff", "--name-only", "HEAD"},
		)
		if err != nil {
			return nil, err
		}
		if exitCode != 0 {
			return nil, fmt.Errorf("git diff --name-only HEAD failed (exit %d): %s",
				exitCode, strings.TrimSpace(stderr))
		}
		for _, f := range getFiles(stdout) {
			fileSet[f] = struct{}{}
		}

		// Bonus : fichiers untracked (optionnel mais utile en local)
		stdout, _, exitCode, err = m.gitRaw(ctx, repo,
			[]string{"ls-files", "--others", "--exclude-standard"},
		)
		if err == nil && exitCode == 0 {
			for _, f := range getFiles(stdout) {
				fileSet[f] = struct{}{}
			}
		}
	}

	// Conversion en slice
	modifiedFiles := make([]string, 0, len(fileSet))
	for f := range fileSet {
		modifiedFiles = append(modifiedFiles, f)
	}

	// ── 3. Groupage par extension ──

	exts := make(map[string][]string)
	exts["all"] = modifiedFiles

	for _, e := range ext {
		parts := strings.SplitN(e, ":", 2)
		name := parts[0]
		var patterns []string

		if len(parts) == 1 {
			// --ext "go"  →  name="go", patterns=["go"]
			patterns = []string{parts[0]}
		} else {
			// --ext "code:go,mod"  →  name="code", patterns=["go","mod"]
			for _, p := range strings.Split(parts[1], ",") {
				patterns = append(patterns, strings.TrimSpace(p))
			}
		}

		var matched []string
		for _, p := range patterns {
			p = strings.TrimSpace(p)
			files := filedir.FilterFileByExtensionOrName(modifiedFiles, "ext", p)
			matched = append(matched, files...)
		}

		if len(matched) > 0 {
			exts[name] = matched
		}
	}

	out, err := json.Marshal(exts)
	if err != nil {
		return nil, err
	}

	return &ChangedFiles{Files: string(out)}, nil
}

func (f *ChangedFiles) GithubOutput(ctx context.Context) string {
	var data map[string][]string

	json.Unmarshal([]byte(f.Files), &data)

	result := ""
	count := len(data) - 1
	i := 0
	for name, files := range data {
		result += fmt.Sprintf("%s=%s", name, strings.Join(files, ","))
		if i < count {
			result += "\n"
		}
		i++
	}

	return result
}

func (f *ChangedFiles) Json(ctx context.Context) string {
	return f.Files
}

func getFiles(input string) []string {
	clean := strings.TrimSpace(input)
	if clean == "" {
		return []string{}
	}

	return strings.Split(clean, "\n")
}
