package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/gitlocal/internal/dagger"
)

type Status struct {
	IsRepo       bool
	IsClean      bool
	HasUntracked bool
	HasModified  bool
	Branch       string
	Head         string
}

// Status exécute `git status --porcelain=v2 --branch` et remplit la structure Status.
//
// Règles métier :
//   - Si ce n’est pas un repo git -> IsRepo=false, pas d’erreur.
//   - Si autre erreur git -> erreur.
func (m *Gitlocal) Status(
	ctx context.Context,
	repo *dagger.Directory,
) (*Status, error) {
	stdout, stderr, exitCode, err := m.gitRaw(ctx, repo,
		[]string{"status", "--porcelain=v2", "--branch"},
	)
	if err != nil {
		return nil, err
	}

	st := &Status{}

	if exitCode != 0 {
		// Cas “pas un dépôt git”
		if strings.Contains(stderr, "not a git repository") {
			st.IsRepo = false
			return st, nil
		}

		// Toute autre erreur est remontée
		return nil, fmt.Errorf("git status failed (exit %d): %s",
			exitCode, strings.TrimSpace(stderr))
	}

	// Ici, on sait que c’est bien un repo git
	st.IsRepo = true

	lines := strings.Split(stdout, "\n")

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}

		// Lignes de méta-données de branche : "# branch.* ..."
		if strings.HasPrefix(line, "# ") {
			fields := strings.Fields(line)
			if len(fields) < 3 {
				continue
			}

			key := fields[1]
			value := fields[2]

			switch key {
			case "branch.oid":
				// "(initial)" => branche encore sans commit
				if value != "(initial)" {
					st.Head = value
				}
			case "branch.head":
				// "(detached)" => pas de nom de branche
				if value != "(detached)" {
					st.Branch = value
				}
			}

			continue
		}

		// Lignes d’entrées (fichiers)
		// Untracked
		if strings.HasPrefix(line, "? ") {
			st.HasUntracked = true
			continue
		}

		// Porcelain v2 : "1 XY ..." ou "2 XY ..."
		parts := strings.SplitN(line, " ", 3)
		if len(parts) < 2 {
			continue
		}

		status := parts[1]
		if len(status) < 2 {
			continue
		}

		x := status[0] // index status
		y := status[1] // worktree status

		// Toute modification, ajout, suppression, renommage...
		if x != ' ' || y != ' ' {
			st.HasModified = true
		}
	}

	st.IsClean = !st.HasUntracked && !st.HasModified

	return st, nil
}

func (s *Status) Stdout(ctx context.Context) string {
	return fmt.Sprintf("%+v\n", s)
}
