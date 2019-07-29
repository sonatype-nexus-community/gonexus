package nexusiq

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func startContainer(t *testing.T) (func() error, int) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		t.Fatal(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}

	// cli.Start container
	// id := "todo"

	return func() error {
		return nil //cli.ContainerStop(context.Background(), id, nil)
	}, 42
}

func TestFunctionalIQ(t *testing.T) {
	// t.Fatal("not implemented")
	stop, port := startContainer(t)
	defer stop()

	t.Log(port)
}
