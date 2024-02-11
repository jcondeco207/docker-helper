package containers

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/rodaine/table"
)

func ShowRunning() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Running containers (%d): \n", len(containers))
	for _, container := range containers {
		fmt.Println("container ", container.Names, " ", container.ID[:10])
	}
}

func GetRunningContainers() []types.Container {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	return containers
}

func GetStoppedContainers() []types.Container {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}

	var stoppedContainers []types.Container
	for _, container := range containers {
		if container.State == "exited" {
			stoppedContainers = append(stoppedContainers, container)
		}
	}

	return stoppedContainers
}

func StopAllContainers() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Print("Stopping container ", container.ID[:10], "... ")
		noWaitTimeout := 0 // to not wait for the container to exit gracefully
		if err := cli.ContainerStop(ctx, container.ID, containertypes.StopOptions{Timeout: &noWaitTimeout}); err != nil {
			panic(err)
		}
		fmt.Println("Success")
	}
}

func StartContainer(containerID string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	fmt.Print("Starting container ", containerID[:10], "... ")

	if err := cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println("container started")
}

func StopContainer(containerID string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	fmt.Print("Stopping container ", containerID[:10], "... ")
	noWaitTimeout := 0 // to not wait for the container to exit gracefully
	if err := cli.ContainerStop(ctx, containerID, containertypes.StopOptions{Timeout: &noWaitTimeout}); err != nil {
		panic(err)
	}
	fmt.Println("container stopped")
}

func DeleteContainer(containerID string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	fmt.Print("Removing container ", containerID[:10], "... ")
	if err := cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{}); err != nil {
		panic(err)
	}
	fmt.Println("container deleted")
}

func ShowAllContainers() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}

	// Output in a table
	containersTable := table.New("ID", "Image", "State")
	for _, container := range containers {
		containersTable.AddRow(container.ID[:10], container.Image, container.State)
	}
}

func ExecFunction(containerID string, command []string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	// Create an exec instance
	execCreateResp, err := cli.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		Cmd:          command,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return err
	}

	// Attach to the exec instance
	hijackedResponse, err := cli.ContainerExecAttach(ctx, execCreateResp.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}
	defer hijackedResponse.Close()

	// Start the exec instance
	err = cli.ContainerExecStart(ctx, execCreateResp.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}

	// Read the output
	stdout, _ := io.ReadAll(hijackedResponse.Reader)
	fmt.Println(string(stdout))

	return nil
}

func AttachToContainer(containerID string) error {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	// Attach to the container
	resp, err := cli.ContainerAttach(ctx, containerID, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return err
	}
	defer resp.Close()

	// Redirect the container's output to the application's stdout and stderr
	go func() {
		_, _ = io.Copy(os.Stdout, resp.Reader)
	}()

	// Wait for the user to exit the terminal
	waitCh, errCh := cli.ContainerWait(ctx, containerID, containertypes.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-waitCh:
	}

	fmt.Println("You have exited the container terminal.")

	return nil
}
