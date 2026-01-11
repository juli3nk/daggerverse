package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/gitlocal/internal/dagger"
)

// HeadInfo décrit l’état de HEAD pour un dépôt.
type HeadInfo struct {
	IsRepo     bool   // false si le répertoire n'est pas un repo git
	Commit     string // SHA de HEAD (vide si repo sans commit)
	Ref        string // ex: "refs/heads/main" (ou vide si détaché)
	Branch     string // ex: "main" (vide si détaché)
	IsDetached bool   // true si HEAD détaché
}

// HeadInfo retourne des infos sur HEAD.
//
// Règles métier :
//   - Si ce n’est pas un repo git -> IsRepo=false, pas d’erreur.
//   - Si repo sans commit -> IsRepo=true, Commit vide, pas d’erreur.
//   - Autre erreur git -> erreur.
func (m *Gitlocal) HeadInfo(
	ctx context.Context,
	repo *dagger.Directory,
) (*HeadInfo, error) {
	info := &HeadInfo{}

	// 1) Vérifier qu’on est bien dans un repo + récupérer la ref symbolique de HEAD
	stdout, stderr, exitCode, err := m.gitRaw(ctx, repo,
		[]string{"rev-parse", "--symbolic-full-name", "HEAD"},
	)
	if err != nil {
		return nil, err
	}

	if exitCode != 0 {
		// Pas un dépôt git
		if strings.Contains(stderr, "not a git repository") {
			info.IsRepo = false
			return info, nil
		}

		return nil, fmt.Errorf("git rev-parse --symbolic-full-name HEAD failed (exit %d): %s",
			exitCode, strings.TrimSpace(stderr))
	}

	info.IsRepo = true
	ref := strings.TrimSpace(stdout)
	if strings.HasPrefix(ref, "refs/") {
		info.Ref = ref
	}

	// 2) Récupérer le nom de branche “humain” ou détecter HEAD détaché
	stdout, stderr, exitCode, err = m.gitRaw(ctx, repo,
		[]string{"rev-parse", "--abbrev-ref", "HEAD"},
	)
	if err != nil {
		return nil, err
	}

	if exitCode != 0 {
		return nil, fmt.Errorf("git rev-parse --abbrev-ref HEAD failed (exit %d): %s",
			exitCode, strings.TrimSpace(stderr))
	}

	abbrev := strings.TrimSpace(stdout)
	if abbrev == "HEAD" {
		// HEAD détaché
		info.IsDetached = true
		info.Branch = ""
	} else {
		info.IsDetached = false
		info.Branch = abbrev
	}

	// 3) SHA de HEAD (peut échouer si repo sans commit)
	stdout, stderr, exitCode, err = m.gitRaw(ctx, repo,
		[]string{"rev-parse", "HEAD"},
	)
	if err != nil {
		return nil, err
	}

	if exitCode != 0 {
		// Cas repo sans commit : git renvoie typiquement un message du genre
		// "does not have any commits yet" ou "needed a single revision"
		if strings.Contains(stderr, "does not have any commits yet") ||
			strings.Contains(stderr, "needed a single revision") {
			// On considère ça comme un état valide :
			// IsRepo = true, Commit vide.
			info.Commit = ""
			return info, nil
		}

		return nil, fmt.Errorf("git rev-parse HEAD failed (exit %d): %s",
			exitCode, strings.TrimSpace(stderr))
	}

	info.Commit = strings.TrimSpace(stdout)

	return info, nil
}
