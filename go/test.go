package main

import (
	"context"
	"fmt"
)

// Test runs Go tests
func (m *Go) Test(
	ctx context.Context,
	// +optional
	// +default="./..."
	packages string,
	// +optional
	race bool,
	// +optional
	cover bool,
) (string, error) {
	execArgs := []string{"go", "test"}
	if race {
		execArgs = append(execArgs, "-race")
	}
	if cover {
		execArgs = append(execArgs, "-cover")
	}

	execArgs = append(execArgs, "-v", packages)

	return dag.Container().
		From(fmt.Sprintf("golang:%s", m.Version)).
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume(fmt.Sprintf("go-mod-%s", m.Version))).
		WithMountedCache("/go/build-cache", dag.CacheVolume(fmt.Sprintf("go-build-%s", m.Version))).
		WithEnvVariable("GOCACHE", "/go/build-cache").
		WithExec(execArgs).
		Stdout(ctx)
}

// TestCoverage runs tests with coverage report
func (m *Go) TestCoverage(
	ctx context.Context,
	// Packages to test
	// +optional
	// +default="./..."
	packages string,
	// Minimum coverage percentage required
	// +optional
	// +default=80
	minCoverage int,
) (string, error) {
	ctr := dag.Container().
		From(fmt.Sprintf("golang:%s", m.Version)).
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+m.Version)).
		WithExec([]string{
			"go", "test",
			"-coverprofile=coverage.out",
			"-covermode=atomic",
			packages,
		})

	// Check coverage threshold
	return ctr.
		WithExec([]string{
			"go", "tool", "cover",
			"-func=coverage.out",
		}).
		Stdout(ctx)
}

// Benchmark runs benchmarks
func (m *Go) Benchmark(
	ctx context.Context,
	// Packages to benchmark
	// +optional
	// +default="./..."
	packages string,
	// Number of times to run each benchmark
	// +optional
	// +default=10
	count int,
) (string, error) {
	return dag.Container().
		From(fmt.Sprintf("golang:%s", m.Version)).
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+m.Version)).
		WithExec([]string{
			"go", "test",
			"-bench=.",
			"-benchmem",
			fmt.Sprintf("-count=%d", count),
			packages,
		}).
		Stdout(ctx)
}

// Race runs tests with race detector
func (m *Go) Race(
	ctx context.Context,
	// +optional
	// +default="./..."
	packages string,
) (string, error) {
	return dag.Container().
		From(fmt.Sprintf("golang:%s", m.Version)).
		WithMountedDirectory("/src", m.Worktree).
		WithWorkdir("/src").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-"+m.Version)).
		WithExec([]string{"go", "test", "-race", packages}).
		Stdout(ctx)
}
