package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"dagger/gitlocal/internal/dagger"
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

	// ── 1. Diff entre les refs (avec fallback si HEAD~1 invalide) ──
	diffRange := fmt.Sprintf("%s...%s", baseRef, headRef)

	stdout, stderr, exitCode, err := gitRaw(ctx, repo,
		[]string{"diff", "--name-only", "--diff-filter=AMR", diffRange},
	)
	if err != nil {
		return nil, err
	}
	if exitCode != 0 {
		// Fallback : si le range échoue (ex: shallow clone, repo à 1 commit),
		// on prend juste les fichiers du dernier commit ou le worktree
		stdout, stderr, exitCode, err = gitRaw(ctx, repo,
			[]string{"diff", "--name-only", "--diff-filter=AMR", "HEAD"},
		)
		if err != nil || exitCode != 0 {
			return nil, fmt.Errorf("git diff failed: %w (stderr: %s)", err, strings.TrimSpace(stderr))
		}
	}

	fileSet := make(map[string]struct{})
	for _, f := range getFiles(stdout) {
		if f != "" {
			fileSet[f] = struct{}{}
		}
	}

	// ── 2. Worktree + untracked (si local) ──
	if headRef == "" || headRef == "HEAD" {
		// Modifiés/staged
		stdout, _, _, _ := gitRaw(ctx, repo,
			[]string{"diff", "--name-only", "--diff-filter=AMR", "HEAD"},
		)
		for _, f := range getFiles(stdout) {
			if f != "" {
				fileSet[f] = struct{}{}
			}
		}

		// Untracked
		stdout, _, _, _ = gitRaw(ctx, repo,
			[]string{"ls-files", "--others", "--exclude-standard"},
		)
		for _, f := range getFiles(stdout) {
			if f != "" {
				fileSet[f] = struct{}{}
			}
		}
	}

	// Conversion en slice
	modifiedFiles := make([]string, 0, len(fileSet))
	for f := range fileSet {
		modifiedFiles = append(modifiedFiles, f)
	}
	sort.Strings(modifiedFiles) // déterministe pour les tests

	// ── 3. Groupage par extension (robuste) ──
	exts := make(map[string][]string)
	exts["all"] = modifiedFiles

	// Construit la map : extension -> nom du groupe
	// Ex: "go" -> "go", "yaml" -> "config", "yml" -> "config"
	groupByExt := make(map[string]string)
	for _, e := range ext {
		parts := strings.SplitN(e, ":", 2)
		name := strings.TrimSpace(parts[0])

		if len(parts) == 1 {
			groupByExt[name] = name
		} else {
			for _, rawExt := range strings.Split(parts[1], ",") {
				extClean := strings.TrimSpace(rawExt)
				if extClean != "" {
					// Détection de conflit silencieuse
					if existing, ok := groupByExt[extClean]; ok {
						fmt.Fprintf(os.Stderr, "warning: extension %q already mapped to %q, overwriting with %q\n",
							extClean, existing, name)
					}
					groupByExt[extClean] = name
				}
			}
		}
	}

	for _, f := range modifiedFiles {
		ext := filepath.Ext(f)             // ".go", ".yaml", ""
		ext = strings.TrimPrefix(ext, ".") // "go", "yaml", ""

		if ext == "" {
			continue // fichiers sans extension restent dans "all" seulement
		}

		if groupName, ok := groupByExt[ext]; ok {
			exts[groupName] = append(exts[groupName], f)
		}
	}

	// ── 4. S'assure que tous les groupes demandés existent (même vides) ──
	for _, e := range ext {
		parts := strings.SplitN(e, ":", 2)
		name := strings.TrimSpace(parts[0])
		if _, ok := exts[name]; !ok {
			exts[name] = []string{}
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
		b, _ := json.Marshal(files)

		result += fmt.Sprintf("%s=%s", name, string(b))
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
