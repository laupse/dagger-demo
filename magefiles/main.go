package main

import (
	"context"
	"os"

	"dagger.io/dagger"
	"github.com/charmbracelet/log"
	platformFormat "github.com/containerd/containerd/platforms"
	"github.com/magefile/mage/mg"
	"github.com/sourcegraph/conc/pool"
)

// Run unit test in math folder
func Test() error {
	log.Info("Test")
	// Starting dagger engine && api session
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	// Reading dir including file
	dir := client.
		Host().
		Directory(".", dagger.HostDirectoryOpts{
			Include: []string{"./math", "go.mod", "go.sum"},
		})

	// Testing in a golang container
	_, err = client.
		Container().
		From("golang:alpine").
		WithWorkdir("/src").
		WithDirectory("/src", dir).
		WithExec([]string{"go", "test", "./math"}).
		Stdout(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Compile local go file into binary
func Build() error {
	log.Info("Build")
	// Starting dagger engine && api session
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	// Reading dir exluding file
	dir := client.
		Host().
		Directory(".", dagger.HostDirectoryOpts{
			Exclude: []string{"./magefiles", "go.work"},
		})

	// Building in a golang container
	_, err = client.
		Container().
		From("golang:alpine").
		WithWorkdir("/src").
		WithDirectory("/src", dir).
		WithExec([]string{"go", "build", "-o", "dagger-demo"}).
		Stdout(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Compile and run the binary
func Run() error {
	log.Info("Run")
	// Starting dagger engine && api session
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	// Reading dir exluding file
	dir := client.
		Host().
		Directory(".", dagger.HostDirectoryOpts{
			Exclude: []string{"./magefiles", "go.work"},
		})

	// Building in a golang container
	bin := client.
		Container().
		From("golang:alpine").
		WithWorkdir("/src").
		WithDirectory("/src", dir).
		WithExec([]string{"go", "build", "-o", "dagger-demo"}).
		File("dagger-demo")

	// Run binary from step above
	output, err := client.
		Container().
		From("alpine").
		WithFile("/bin/", bin).
		WithExec([]string{"dagger-demo", "2", "2"}).
		Stdout(ctx)
	if err != nil {
		return err
	}
	log.Info(output)

	return nil
}

// Test, Build, Run
func All() {
	mg.SerialDeps(Test, Build, Run)
}

// Build binary in multiple architecture concurrently
func BuildConcurrent() error {
	log.Info("BuildConcurrent")

	platforms := []dagger.Platform{
		"linux/amd64", // a.k.a. x86_64
		"linux/arm64", // a.k.a. aarch64
		"linux/s390x", // a.k.a. IBM S/390
	}

	// Starting dagger engine && api session
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	// Reading dir exluding file
	dir := client.
		Host().
		Directory(".", dagger.HostDirectoryOpts{
			Exclude: []string{"./magefiles", "go.work"},
		})

	p := pool.New().WithErrors()
	// Building in a golang container
	for _, platform := range platforms {
		p.Go(func() error {
			_, err = client.
				Container().
				From("golang:alpine").
				WithEnvVariable("GOOS", "linux").
				WithEnvVariable("GOARCH", architectureOf(platform)).
				WithWorkdir("/src").
				WithDirectory("/src", dir).
				WithExec([]string{"go", "build", "-o", "dagger-demo"}).
				Stdout(ctx)
			return err
		})
	}
	err = p.Wait()
	if err != nil {
		return err
	}

	return nil
}

// util that returns the architecture of the provided platform
func architectureOf(platform dagger.Platform) string {
	return platformFormat.MustParse(string(platform)).Architecture
}

// Run redis cmd against a redis server declared in Dagger
func Service() error {
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	redis := client.
		Container().
		From("redis:alpine").
		WithExec(nil)

	redisCLI := client.
		Container().
		From("redis:alpine").
		WithServiceBinding("redis-server", redis).
		WithEntrypoint([]string{"redis-cli", "-h", "redis-server"})

	_, err = redisCLI.
		WithExec([]string{"set", "foo", "abc"}).
		WithExec([]string{"save"}).
		Stdout(ctx)
	if err != nil {
		return err
	}

	val, err := redisCLI.
		WithExec([]string{"get", "foo"}).
		Stdout(ctx)
	if err != nil {
		return err
	}
	log.Infof("foo = %s", val)

	return nil
}

// Try to leak a secret in dagger
func Secret() error {
	ctx := context.Background()
	client, err := dagger.Connect(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	secret := client.SetSecret("secret", "this-is-not-going-to-leak")

	out, err := client.
		Container().
		From("alpine").
		WithSecretVariable("SECRET", secret).
		WithExec([]string{"sh", "-c", `echo -e "secret env data: $SECRET"`}).
		Stdout(ctx)
	if err != nil {
		panic(err)
	}

	log.Info(out)

	return nil
}

// Download image used in the engine
func Image() error {
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	images := []string{"alpine", "golang:alpine", "redis:alpine"}
	for _, image := range images {
		_, err = client.
			Container().
			From(image).
			WithExec([]string{"true"}).
			Stdout(ctx)
		if err != nil {
			return err
		}
	}

	// platforms := []dagger.Platform{
	// 	"linux/amd64", // a.k.a. x86_64
	// 	"linux/arm64", // a.k.a. aarch64
	// 	"linux/s390x", // a.k.a. IBM S/390
	// }

	// p := pool.New().WithErrors()
	// // Building in a golang container
	// for _, platform := range platforms {
	// 	p.Go(func() error {
	// 		_, err = client.Container(dagger.ContainerOpts{
	// 			Platform: platform,
	// 		}).
	// 			From("golang:alpine").
	// 			WithExec([]string{"true"}).
	// 			Stdout(ctx)
	// 		return err
	// 	})
	// }
	// err = p.Wait()
	// if err != nil {
	// 	return err
	// }

	return nil
}
