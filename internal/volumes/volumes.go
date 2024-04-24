package volumes

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/rodaine/table"
)

func DeleteVolume(volumeID string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	fmt.Print("Removing volume ", volumeID[:10], "... ")
	if err := cli.VolumeRemove(ctx, volumeID, false); err != nil {
		panic(err)
	}
	fmt.Println("Volume deleted")
}

func ShowAllVolumes() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	volumes, err := cli.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		panic(err)
	}

	vlist := volumes.Volumes

	volumesTable := table.New("Name", "ID", "Usage")

	for _, volume := range vlist {
		volumesTable.AddRow(volume.Name, volume.ClusterVolume.ID, volume.UsageData)
	}

	volumesTable.Print()
}
