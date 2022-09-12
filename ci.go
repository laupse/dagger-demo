package main

import (
	"context"
	"fmt"

	"github.com/dagger/cloak/engine"
	"github.com/dagger/cloak/examples/alpine/gen/core"
)

func main() {
	context := context.Background()
	err := engine.Start(context, &engine.Config{}, func(ctx engine.Context) error {
		var err error
		result, err := core.Dockerfile(ctx, ctx.Workdir, "Dockerfile")
		if err != nil {
			return err
		}

		_, err = core.PushImage(ctx, result.GetCore().Filesystem.ID, "kind-registry:5000/hello-app:latest")
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Println("error", err)
	}
}
