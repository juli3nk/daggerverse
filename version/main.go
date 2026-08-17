package main

import (
	"context"
	"strings"

	"dagger/version/internal/dagger"
)

type Version struct{}

func New() *Version {
	return &Version{}
}

// Resolve computes a version string from explicit inputs.
// Priority: explicit version > semantic version tag > short commit > "unknown"
// Appends "-dirty" if uncommitted changes are detected.
func (m *Version) Resolve(
	// Explicit version provided by the user (e.g., "1.2.3"). Empty string skips this.
	// +optional
	version string,
	// Latest git tag (e.g., "v1.2.3" or "1.2.3"). Empty string skips this.
	// +optional
	tag string,
	// Latest git commit SHA. Used as fallback when no tag or version is provided.
	// +optional
	commit string,
	// Whether there are uncommitted changes in the working tree.
	// +optional
	// +default=false
	dirty bool,
) string {
	var result string

	switch {
	case version != "":
		result = version
	case tag != "":
		result = tag
	case commit != "":
		result = shortSHA(commit)
	default:
		result = "unknown"
	}

	if dirty && !strings.HasSuffix(result, "-dirty") {
		result += "-dirty"
	}

	return result
}

// ResolveFromSource extracts git metadata from a source directory using the gitlocal module
// and computes a version string automatically.
func (m *Version) ResolveFromSource(
	ctx context.Context,
	// Source directory containing a git repository.
	source *dagger.Directory,
	// Explicit version override. If empty, version is auto-detected from git metadata.
	// +optional
	version string,
) (string, error) {
	git := dag.Gitlocal()

	// Get latest commit
	headInfo := git.HeadInfo(source)
	commit, err := headInfo.Commit(ctx)
	if err != nil {
		return "", err
	}

	// Get latest tag
	tag, err := git.LatestTag(ctx, source)
	if err != nil {
		return "", err
	}

	// Check if dirty
	status := git.Status(source)

	isRepo, err := status.IsRepo(ctx)
	if err != nil {
		return "", err
	}
	isClean, err := status.IsClean(ctx)
	if err != nil {
		return "", err
	}

	dirty := isRepo && !isClean

	return m.Resolve(version, tag, commit, dirty), nil
}

// shortSHA returns the first 7 characters of a commit SHA.
func shortSHA(commit string) string {
	if len(commit) < 7 {
		return commit
	}
	return commit[:7]
}
