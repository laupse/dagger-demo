package main

import (
	"context"
	"fmt"
	"os"

	"go.dagger.io/dagger/engine"
	"go.dagger.io/dagger/router"
	"go.dagger.io/dagger/sdk/go/dagger"
	"go.dagger.io/dagger/sdk/go/dagger/api"
)

func Test() {
	ctx := context.Background()
	if err := engine.Start(ctx, &engine.Config{}, func(ctx context.Context, r *router.Router) error {
		client, err := dagger.Connect(ctx)
		if err != nil {
			return err
		}
		core := client.Core()

		workdir, err := core.Host().Workdir().Read().Directory(".").ID(ctx)
		if err != nil {
			return err
		}
		execOpts := api.ContainerExecOpts{
			Args: []string{"/usr/local/go/bin/go", "test", "./math", "-v"},
		}
		_, err = core.Container().From("golang:1.19.1-alpine3.16").WithMountedDirectory("/go/src/hello", workdir).WithWorkdir("/go/src/hello").Exec(execOpts).Stdout().Contents(ctx)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Push() {
	ctx := context.Background()
	if err := engine.Start(ctx, &engine.Config{}, func(ctx context.Context, r *router.Router) error {
		client, err := dagger.Connect(ctx)
		if err != nil {
			return err
		}
		core := client.Core()

		workdir, err := core.Host().Workdir().Read().ID(ctx)
		if err != nil {
			return err
		}

		_, err = core.Container().Build(workdir, api.ContainerBuildOpts{
			Dockerfile: "Dockerfile",
		}).Publish(ctx, "ghcr.io/laupse/hello-world:dagger")
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
