package containers

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"golang.org/x/crypto/ssh/terminal"
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

	fmt.Printf("\nAll containers:\n")

	for _, container := range containers {
		//status := "running^stopped"

		fmt.Printf("\nContainer ID: %s\n", container.ID[:10])
		fmt.Printf("Image: %s\n", container.Image)
		fmt.Printf("Status: %s\n", container.State)
		fmt.Printf("\n----------------\n")
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

	// Start the exec instance
	err = cli.ContainerExecStart(ctx, execCreateResp.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}

	// Wait for the exec instance to complete (optional)
	/*_, err = cli.ContainerExecInspect(ctx, execCreateResp.ID)
	if err != nil {
		return err
	}*/

	return nil
}

func AttachToContainer(containerID string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(fd, oldState)

	execConfig := types.ExecConfig{
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"bash"},
	}

	exec, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		panic(err)
	}

	resp, err := cli.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		panic(err)
	}

	go io.Copy(os.Stdout, resp.Reader)
	go io.Copy(resp.Conn, os.Stdin)

	resp.Close()
	resp.Conn.Close()
}
