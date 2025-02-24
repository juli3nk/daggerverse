package main

import (
	"context"
	"fmt"

	"dagger/platformio/internal/dagger"
)

type Platformio struct{}

func (m *Platformio) Check(
	ctx context.Context,
	source *dagger.Directory,
) (string, error) {
	execArgs := []string{
		"check",
	}
	opts := dagger.ContainerWithExecOpts{
		UseEntrypoint: true,
	}

	return m.pio().
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(execArgs, opts).
		Stdout(ctx)
}

func (m *Platformio) Run(
	ctx context.Context,
	source *dagger.Directory,
	env string,
) *dagger.File {
	execArgs := []string{
		"run",
	}
	opts := dagger.ContainerWithExecOpts{
		UseEntrypoint: true,
	}

	return m.pio().
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(execArgs, opts).
		File(fmt.Sprintf(".pio/build/%s/firmware.bin", env))
}

func (m *Platformio) pio() *dagger.Container {
	return dag.Apko().Wolfi().
		WithPackages([]string{"python-3.13"}).
		Container().
		WithExec([]string{"python", "-m", "ensurepip", "--default-pip"}).
		WithExec([]string{"python", "-m", "pip", "install", "-U", "platformio"}).
		WithEntrypoint([]string{"pio"})
}
