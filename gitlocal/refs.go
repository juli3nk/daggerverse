package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"dagger/gitlocal/internal/dagger"
)

type ChangeRefs struct {
	Base string `json:"base"`
	Head string `json:"head"`
}

// ChangeRefs détecte la portée du changement.
//   - CI PR       : utilise la branche cible (origin/<base>)
//   - CI Push     : utilise le before/after du webhook
//   - Local       : détecte origin/main, l'upstream, ou fallback HEAD~1
//
// Tu peux forcer une base avec le paramètre `base` (ex: "origin/develop").
func (m *Gitlocal) ChangeRefs(
	ctx context.Context,
	// +defaultPath="."
	repo *dagger.Directory,
	// +optional
	base string,
) (*ChangeRefs, error) {
	// Manual override
	if base != "" {
		return &ChangeRefs{Base: base, Head: "HEAD"}, nil
	}

	// GitHub Pull Request
	if baseRef := os.Getenv("GITHUB_BASE_REF"); baseRef != "" {
		return &ChangeRefs{
			Base: fmt.Sprintf("origin/%s", baseRef),
			Head: "HEAD",
		}, nil
	}

	// GitLab Merge Request
	if targetBranch := os.Getenv("CI_MERGE_REQUEST_TARGET_BRANCH_NAME"); targetBranch != "" {
		head := os.Getenv("CI_COMMIT_SHA")
		if head == "" {
			head = "HEAD"
		}
		return &ChangeRefs{
			Base: fmt.Sprintf("origin/%s", targetBranch),
			Head: head,
		}, nil
	}

	// GitHub Push (before/after from payload)
	if eventPath := os.Getenv("GITHUB_EVENT_PATH"); eventPath != "" {
		payload, err := os.ReadFile(eventPath)
		if err == nil {
			var event struct {
				Before string `json:"before"`
				After  string `json:"after"`
			}
			if err := json.Unmarshal(payload, &event); err == nil {
				if event.Before != "" && event.Before != "0000000000000000000000000000000000000000" {
					return &ChangeRefs{Base: event.Before, Head: event.After}, nil
				}
			}
		}
	}

	// GitLab Push
	if before := os.Getenv("CI_COMMIT_BEFORE_SHA"); before != "" && before != "0000000000000000000000000000000000000000" {
		return &ChangeRefs{Base: before, Head: os.Getenv("CI_COMMIT_SHA")}, nil
	}

	// Local: git introspection in the container
	localBase, err := m.detectLocalBase(ctx, repo)
	if err != nil {
		return &ChangeRefs{Base: "HEAD~1", Head: "HEAD"}, nil // hard fallback
	}
	return &ChangeRefs{Base: localBase, Head: "HEAD"}, nil
}

func (r *ChangeRefs) GithubOutput(ctx context.Context) string {
	return fmt.Sprintf("base=%s\nhead=%s", r.Base, r.Head)
}

func (r *ChangeRefs) Json(ctx context.Context) string {
	b, err := json.Marshal(r)
	if err != nil {
		return ""
	}

	return string(b)
}

// detectLocalBase exécute git dans le conteneur pour trouver la branche de base.
func (m *Gitlocal) detectLocalBase(ctx context.Context, repo *dagger.Directory) (string, error) {
	script := `base=$(git rev-parse --abbrev-ref '@{u}' 2>/dev/null || true)
if [ -n "$base" ]; then echo "$base"; exit 0; fi

base=$(git symbolic-ref --short refs/remotes/origin/HEAD 2>/dev/null || true)
if [ -n "$base" ]; then echo "$base"; exit 0; fi

for b in origin/main origin/master origin/trunk; do
  if git rev-parse --verify "$b" >/dev/null 2>&1; then echo "$b"; exit 0; fi
done

echo "HEAD~1"`

	out, err := dag.Container().
		From(fmt.Sprintf("golang:%s", "latest")).
		WithMountedDirectory("/src", repo).
		WithWorkdir("/src").
		WithExec([]string{"sh", "-ce", script}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}
