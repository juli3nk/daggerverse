package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/gitlocal/internal/dagger"
)

// IsAncestor vérifie si ancestorRef est ancêtre de descendantRef.
//
// Règles :
//   - ancestorRef obligatoire.
//   - descendantRef vide → "HEAD".
//   - Si un des refs n’existe pas → erreur claire (pas false silencieux).
//   - Si tout existe → true/false selon le graphe de commits.
func (m *Gitlocal) IsAncestor(
	ctx context.Context,
	repo *dagger.Directory,
	ancestorRef string,
	descendantRef string,
) (bool, error) {
	if ancestorRef == "" {
		return false, fmt.Errorf("ancestorRef is required")
	}
	if descendantRef == "" {
		descendantRef = "HEAD"
	}

	// On valide d'abord que les refs existent, pour produire des erreurs lisibles.
	if err := m.ensureRefExists(ctx, repo, "ancestor", ancestorRef); err != nil {
		return false, err
	}
	if err := m.ensureRefExists(ctx, repo, "descendant", descendantRef); err != nil {
		return false, err
	}

	// Ensuite on délègue à git merge-base --is-ancestor.
	_, stderr, exitCode, err := m.gitRaw(ctx, repo,
		[]string{"merge-base", "--is-ancestor", ancestorRef, descendantRef},
	)
	if err != nil {
		return false, err
	}

	switch exitCode {
	case 0:
		// ancestorRef est ancêtre de descendantRef
		return true, nil
	case 1:
		// ancestorRef n'est PAS ancêtre de descendantRef
		return false, nil
	default:
		// Erreur inattendue (problème de repo, etc.)
		return false, fmt.Errorf(
			"git merge-base --is-ancestor %s %s failed (exit %d): %s",
			ancestorRef, descendantRef, exitCode, strings.TrimSpace(stderr),
		)
	}
}

// ensureRefExists vérifie qu'un ref est résolvable par git rev-parse.
// label est utilisé pour un message d’erreur plus clair ("ancestor" / "descendant").
func (m *Gitlocal) ensureRefExists(
	ctx context.Context,
	repo *dagger.Directory,
	label string,
	ref string,
) error {
	_, stderr, exitCode, err := m.gitRaw(ctx, repo,
		[]string{"rev-parse", "--verify", ref},
	)
	if err != nil {
		return err
	}

	if exitCode == 0 {
		return nil
	}

	stderrTrim := strings.TrimSpace(stderr)

	// Cas "pas un dépôt git" → on renvoie une erreur explicite.
	if strings.Contains(stderrTrim, "not a git repository") {
		return fmt.Errorf("not a git repository: %s", stderrTrim)
	}

	// Cas ref introuvable.
	return fmt.Errorf("%s ref %q not found: %s", label, ref, stderrTrim)
}
