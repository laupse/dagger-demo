package main

import (
	"context"
	"os"

	"dagger.io/dagger"
)

func Test() error {
	ctx := context.Background()

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}

	workdir, err := client.Host().Workdir().Read().Directory(".").ID(ctx)
	if err != nil {
		return err
	}

	_, err = client.
		Container().
		From("golang:1.19.1-alpine3.16").
		WithMountedDirectory("/go/src/hello", workdir).
		WithWorkdir("/go/src/hello").
		Exec(dagger.ContainerExecOpts{
			Args: []string{"/usr/local/go/bin/go", "test", "./math", "-v"},
		}).
		Stdout().
		Contents(ctx)
	if err != nil {
		return err
	}

	return nil
}

func Push() error {
	ctx := context.Background()

	client, err := dagger.Connect(ctx)
	if err != nil {
		return err
	}

	workdir, err := client.Host().Workdir().Read().ID(ctx)
	if err != nil {
		return err
	}

	_, err = client.Container().Build(workdir, dagger.ContainerBuildOpts{
		Dockerfile: "Dockerfile",
	}).Publish(ctx, "ghcr.io/laupse/hello-world:dagger")
	if err != nil {
		return err
	}

	return nil
}
