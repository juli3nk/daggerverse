package main

import (
	"context"

	"dagger/goreleaser/internal/dagger"
)

func (m *Goreleaser) Release(
	ctx context.Context,
	// Automatically sets --snapshot if the repository is dirty
	// +optional
	// +default=false
	autoSnapshot bool,
	// Removes the 'dist' directory
	// +optional
	// +default=false
	clean bool,
	// Load configuration from file
	// +optional
	config string,
	// Whether to set the release to draft. Overrides release.draft in the configuration file
	// +optional
	// +default=false
	draft bool,
	// Whether to abort the release publishing on the first error
	// +optional
	// +default=false
	failFast bool,
	// Amount tasks to run concurrently (default: number of CPUs)
	// +optional
	parallelism string,
	// Load custom release notes footer from a markdown file
	// +optional
	// +default=false
	releaseFooter bool,
	// Load custom release notes footer from a templated markdown file (overrides --release-footer)
	// +optional
	// +default=false
	releaseFooterTmpl bool,
	// Load custom release notes header from a markdown file
	// +optional
	// +default=false
	releaseHeader bool,
	// Load custom release notes header from a templated markdown file (overrides --release-header)
	// +optional
	// +default=false
	releaseHeaderTmpl bool,
	// Load custom release notes from a markdown file (will skip GoReleaser changelog generation)
	// +optional
	// +default=false
	releaseNotes bool,
	// Load custom release notes from a templated markdown file (overrides --release-notes)
	// +optional
	// +default=false
	releaseNotesTmpl bool,
	// Skip the given options (valid options are: before, post-hooks, pre-hooks, validate)
	// +optional
	skip string,
	// Skip the given options (valid options are announce, archive, aur, aur-source, before, chocolatey, docker, flatpak, homebrew, ko, makeself, mcp, nfpm, nix, notarize, publish, sbom, scoop, sign, snapcraft, srpm, validate, winget)
	// +optional
	// +default=false
	snapshot bool,
	// Timeout to the entire release process
	// +optional
	// +default="1h0m0s"
	timeout string,
	// Enable verbose mode
	// +optional
	// +default=false
	verboseOutput bool,
) *dagger.Container {
	execArgs := []string{"release"}

	if autoSnapshot {
		execArgs = append(execArgs, "--auto-snapshot")
	}

	if clean {
		execArgs = append(execArgs, "--clean")
	}

	if config != "" {
		execArgs = append(execArgs, "--config", config)
	}

	if draft {
		execArgs = append(execArgs, "--draft")
	}

	if failFast {
		execArgs = append(execArgs, "--fail-fast")
	}

	if parallelism != "" {
		execArgs = append(execArgs, "--parallelism", parallelism)
	}

	if releaseFooter {
		execArgs = append(execArgs, "--release-footer")
	}

	if releaseFooterTmpl {
		execArgs = append(execArgs, "--release-footer-tmpl")
	}

	if releaseHeader {
		execArgs = append(execArgs, "--release-header")
	}

	if releaseHeaderTmpl {
		execArgs = append(execArgs, "--release-header-tmpl")
	}

	if releaseNotes {
		execArgs = append(execArgs, "--release-notes")
	}

	if releaseNotesTmpl {
		execArgs = append(execArgs, "--release-notes-tmpl")
	}

	if skip != "" {
		execArgs = append(execArgs, "--skip", skip)
	}

	if snapshot {
		execArgs = append(execArgs, "--snapshot")
	}

	if timeout != "" {
		execArgs = append(execArgs, "--timeout", timeout)
	}

	if verboseOutput {
		execArgs = append(execArgs, "--verbose")
	}

	return m.run(execArgs)
}
