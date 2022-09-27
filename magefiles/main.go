package main

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"go.dagger.io/dagger/engine"
	"go.dagger.io/dagger/sdk/go/dagger"
)

func Push() {
	context := context.Background()
	err := engine.Start(context, &engine.Config{}, func(ctx engine.Context) error {
		client, err := dagger.Client(ctx)
		if err != nil {
			return err
		}

		srcId, err := source(ctx, client)
		if err != nil {
			return err
		}
		fmt.Println(srcId)

		id, err := dockerfile(ctx, client, srcId)
		if err != nil {
			return err
		}
		fmt.Println(id)

		err = push(ctx, client, id, "ghcr.io/laupse/hello-world")
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func source(ctx context.Context, client graphql.Client) (string, error) {
	req := &graphql.Request{
		Query: `
			{
				host {
					workdir {
						read {
							id
						}
					}
				}
			}
		`,
		Variables: map[string]any{},
	}
	resp := struct {
		Host struct {
			Workdir struct {
				Read struct {
					Id string
				}
			}
		}
	}{}
	err := client.MakeRequest(ctx, req, &graphql.Response{Data: &resp})
	if err != nil {
		return "", err
	}

	return resp.Host.Workdir.Read.Id, nil

}

func dockerfile(ctx context.Context, client graphql.Client, srcId string) (string, error) {
	req := &graphql.Request{
		Query: `
			query ($fsid: FSID!) {
				core {
					filesystem(id: $fsid) {
						dockerbuild(dockerfile: "Dockerfile") {
							id
						}
					}
				}
			}
		`,
		Variables: map[string]any{
			"fsid": srcId,
		},
	}
	resp := struct {
		Host struct {
			Workdir struct {
				Read struct {
					Dockerbuild struct {
						Id string
					}
				}
			}
		}
	}{}
	err := client.MakeRequest(ctx, req, &graphql.Response{Data: &resp})
	if err != nil {
		return "", err
	}

	return resp.Host.Workdir.Read.Dockerbuild.Id, nil
}

func push(ctx context.Context, client graphql.Client, id, ref string) error {

	req := &graphql.Request{
		Query: `
		query ($fsid: FSID!, $ref: String!) {
			core {
				filesystem(id: $fsid)  {
					pushImage(ref: $ref) 
				}
			}
		}`,
		Variables: map[string]any{
			"fsid": id,
			"ref":  ref,
		},
	}

	resp := struct {
		push bool
	}{}

	err := client.MakeRequest(ctx, req, &graphql.Response{Data: &resp})
	if err != nil {
		return err
	}

	return nil

}
