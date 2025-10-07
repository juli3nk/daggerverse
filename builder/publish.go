package main

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"dagger/builder/internal/dagger"
)

type ImagePublishConfig struct {
	name       string
	containers []*dagger.Container
	latest     bool
}

// PublishContainer publishes a container image to a registry
func (m *Builder) PublishImage(
	ctx context.Context,
	imageConfigs []ImagePublishConfig,
	version string,
	registryAddress string,
	registryNamespace string,
	registryUsername string,
	registrySecret *dagger.Secret,
) error {
	if err := validateRegistryParams(registryAddress, registryNamespace, registryUsername, registrySecret); err != nil {
		return err
	}

	ctr := dag.Container().WithRegistryAuth(registryAddress, registryUsername, registrySecret)

	// Parallel publication
	var wg sync.WaitGroup
	errorsChan := make(chan error, len(imageConfigs))

	for _, config := range imageConfigs {
		if len(config.containers) == 0 {
			continue // Skip if no containers for this type
		}

		wg.Add(1)
		go func(cfg ImagePublishConfig) {
			defer wg.Done()
			err := publishImageConfig(ctx, ctr, cfg, version, registryAddress, registryNamespace)
			errorsChan <- err
		}(config)
	}

	wg.Wait()
	close(errorsChan)

	// Collect the errors
	var publishErrors []error
	for err := range errorsChan {
		if err != nil {
			publishErrors = append(publishErrors, err)
		}
	}

	if len(publishErrors) > 0 {
		return fmt.Errorf("publish failed: %w", errors.Join(publishErrors...))
	}

	return nil
}

// Helper method to publish an image config
func publishImageConfig(
	ctx context.Context,
	ctr *dagger.Container,
	config ImagePublishConfig,
	version string,
	registryAddress string,
	registryNamespace string,
) error {
	imageName := fmt.Sprintf("%s/%s/%s", registryAddress, registryNamespace, config.name)

	imageNameTags := []string{
		fmt.Sprintf("%s:%s", imageName, version),
	}

	if config.latest {
		imageNameTags = append(imageNameTags, fmt.Sprintf("%s:latest", imageName))
	}

	// Publish the images (sequential per image, but images in parallel)
	for _, image := range imageNameTags {
		fmt.Printf("Publishing %s...\n", image)
		_, err := ctr.Publish(ctx, image, dagger.ContainerPublishOpts{
			PlatformVariants: config.containers,
		})
		if err != nil {
			return fmt.Errorf("failed to publish %s: %w", image, err)
		}
		fmt.Printf("Published: %s\n", image)
	}

	return nil
}

func validateRegistryParams(address, namespace, username string, secret *dagger.Secret) error {
	if address == "" {
		return nil // No registry = OK
	}

	if namespace == "" {
		return fmt.Errorf("registry namespace required when address is specified")
	}

	if username == "" {
		return fmt.Errorf("registry username required when address is specified")
	}

	if secret == nil {
		return fmt.Errorf("registry secret required when address is specified")
	}

	return nil
}
