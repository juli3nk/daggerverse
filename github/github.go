package main

import (
	"context"
	"fmt"

	"dagger/publisher/internal/dagger"
)

// PublishGithubRelease creates a GitHub release with artifacts
func (m *Publisher) PublishGithubRelease(
	ctx context.Context,
	// Repository in format owner/repo
	repository string,
	// Release tag (e.g., v1.0.0)
	tag string,
	// Release title
	// +optional
	title string,
	// Release body/changelog
	// +optional
	body string,
	// Mark as prerelease
	// +optional
	// +default=false
	prerelease bool,
	// Mark as draft
	// +optional
	// +default=false
	draft bool,
	// Directory containing artifacts to upload
	// +optional
	artifacts *dagger.Directory,
) (string, error) {
	if m.GithubToken == nil {
		return "", fmt.Errorf("github token is required")
	}

	if title == "" {
		title = tag
	}

	ctr := m.githubCLI().
		WithSecretVariable("GITHUB_TOKEN", m.GithubToken)

	// Create release
	createArgs := []string{
		"gh", "release", "create", tag,
		"--repo", repository,
		"--title", title,
	}

	if body != "" {
		createArgs = append(createArgs, "--notes", body)
	}
	if prerelease {
		createArgs = append(createArgs, "--prerelease")
	}
	if draft {
		createArgs = append(createArgs, "--draft")
	}

	// If no artifacts, just create release
	if artifacts == nil {
		return ctr.WithExec(createArgs).Stdout(ctx)
	}

	// Mount artifacts and add them to release
	ctr = ctr.WithMountedDirectory("/artifacts", artifacts)

	// Get list of files
	files, err := artifacts.Entries(ctx)
	if err != nil {
		return "", err
	}

	// Add each file to the release command
	for _, file := range files {
		createArgs = append(createArgs, fmt.Sprintf("/artifacts/%s", file))
	}

	return ctr.WithExec(createArgs).Stdout(ctx)
}

// GetLatestRelease gets the latest release tag
func (m *Publisher) GetLatestRelease(
	ctx context.Context,
	// Repository in format owner/repo
	repository string,
) (string, error) {
	if m.GithubToken == nil {
		return "", fmt.Errorf("github token is required")
	}

	return m.githubCLI().
		WithSecretVariable("GITHUB_TOKEN", m.GithubToken).
		WithExec([]string{
			"gh", "release", "view",
			"--repo", repository,
			"--json", "tagName",
			"--jq", ".tagName",
		}).
		Stdout(ctx)
}

// DeleteRelease deletes a GitHub release
func (m *Publisher) DeleteRelease(
	ctx context.Context,
	// Repository in format owner/repo
	repository string,
	// Release tag to delete
	tag string,
	// Also delete the git tag
	// +optional
	// +default=false
	deleteTag bool,
) (string, error) {
	if m.GithubToken == nil {
		return "", fmt.Errorf("github token is required")
	}

	args := []string{
		"gh", "release", "delete", tag,
		"--repo", repository,
		"--yes",
	}

	if deleteTag {
		args = append(args, "--cleanup-tag")
	}

	return m.githubCLI().
		WithSecretVariable("GITHUB_TOKEN", m.GithubToken).
		WithExec(args).
		Stdout(ctx)
}

// UploadArtifact uploads additional artifacts to an existing release
func (m *Publisher) UploadArtifact(
	ctx context.Context,
	// Repository in format owner/repo
	repository string,
	// Release tag
	tag string,
	// Artifact file to upload
	artifact *dagger.File,
	// Artifact filename
	filename string,
) (string, error) {
	if m.GithubToken == nil {
		return "", fmt.Errorf("github token is required")
	}

	return m.githubCLI().
		WithSecretVariable("GITHUB_TOKEN", m.GithubToken).
		WithMountedFile(fmt.Sprintf("/tmp/%s", filename), artifact).
		WithExec([]string{
			"gh", "release", "upload", tag,
			fmt.Sprintf("/tmp/%s", filename),
			"--repo", repository,
			"--clobber",
		}).
		Stdout(ctx)
}

// githubCLI returns a container with GitHub CLI installed
func (m *Publisher) githubCLI() *dagger.Container {
	return dag.Container().
		From("ghcr.io/cli/cli:latest")
}
