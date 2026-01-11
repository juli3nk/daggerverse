package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/gitlocal/internal/dagger"
)

type ChangedFilesOptions struct {
	// Inclure les fichiers non suivis (git ls-files --others --exclude-standard)
	IncludeUntracked bool

	// Si non vide, ne garder que les fichiers dont le chemin commence par
	// au moins un des préfixes fournis (ex: "apps/api", "libs/shared").
	FilterPaths []string
}

// ChangedFiles retourne la liste des fichiers modifiés entre baseRef et headRef.
//
// Comportement :
//   - baseRef est obligatoire (ex : "origin/main").
//   - si headRef est vide ou "HEAD" :
//   - fichiers modifiés par des commits entre baseRef et HEAD
//     (git diff --name-only baseRef...HEAD)
//   - + fichiers modifiés dans le worktree/stage par rapport à HEAD
//     (git diff --name-only HEAD)
//   - + éventuellement fichiers non suivis (IncludeUntracked).
//   - si headRef est non vide et != "HEAD" :
//   - uniquement les fichiers modifiés entre baseRef et headRef
//     (git diff --name-only baseRef...headRef) — commits uniquement.
//   - Les résultats sont dédupliqués et filtrés par FilterPaths si fourni.
func (m *Gitlocal) ChangedFiles(
	ctx context.Context,
	repo *dagger.Directory,
	baseRef string,
	headRef string,
	opts *ChangedFilesOptions,
) ([]string, error) {
	if baseRef == "" {
		return nil, fmt.Errorf("baseRef is required")
	}
	if opts == nil {
		opts = &ChangedFilesOptions{}
	}

	actualHead := headRef
	if actualHead == "" {
		actualHead = "HEAD"
	}

	seen := make(map[string]struct{})
	var result []string

	// Helper pour ajouter des chemins, avec trim + filtre + dédup.
	addPaths := func(stdout string) {
		for _, raw := range strings.Split(stdout, "\n") {
			path := strings.TrimSpace(raw)
			if path == "" {
				continue
			}
			if !matchesFilter(path, opts.FilterPaths) {
				continue
			}
			if _, ok := seen[path]; ok {
				continue
			}
			seen[path] = struct{}{}
			result = append(result, path)
		}
	}

	// 1) Fichiers modifiés par des commits entre baseRef et actualHead.
	//
	// On utilise baseRef...actualHead (triple dot) pour le cas typique :
	// "tous les changements portés par ma branche par rapport à origin/main".
	// Si tu préfères une diff stricte entre deux points, remplace par baseRef..actualHead.
	stdout, stderr, exitCode, err := m.gitRaw(ctx, repo,
		[]string{"diff", "--name-only", fmt.Sprintf("%s...%s", baseRef, actualHead)},
	)
	if err != nil {
		return nil, err
	}
	if exitCode != 0 {
		return nil, fmt.Errorf(
			"git diff --name-only %s...%s failed (exit %d): %s",
			baseRef, actualHead, exitCode, strings.TrimSpace(stderr),
		)
	}
	addPaths(stdout)

	// 2) Si on est sur HEAD (ou équivalent) : ajouter les changements non commités
	//    (worktree + index) par rapport à HEAD.
	if headRef == "" || headRef == "HEAD" {
		stdout, stderr, exitCode, err = m.gitRaw(ctx, repo,
			[]string{"diff", "--name-only", "HEAD"},
		)
		if err != nil {
			return nil, err
		}
		if exitCode != 0 {
			return nil, fmt.Errorf(
				"git diff --name-only HEAD failed (exit %d): %s",
				exitCode, strings.TrimSpace(stderr),
			)
		}
		addPaths(stdout)
	}

	// 3) Optionnel : fichiers non suivis.
	if opts.IncludeUntracked {
		stdout, stderr, exitCode, err = m.gitRaw(ctx, repo,
			[]string{"ls-files", "--others", "--exclude-standard"},
		)
		if err != nil {
			return nil, err
		}
		if exitCode != 0 {
			return nil, fmt.Errorf(
				"git ls-files --others --exclude-standard failed (exit %d): %s",
				exitCode, strings.TrimSpace(stderr),
			)
		}
		addPaths(stdout)
	}

	return result, nil
}

// matchesFilter retourne true si path matche au moins un des préfixes.
// Si filters est vide, tout passe.
func matchesFilter(path string, filters []string) bool {
	if len(filters) == 0 {
		return true
	}
	for _, prefix := range filters {
		p := strings.TrimSpace(prefix)
		if p == "" {
			continue
		}
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}
