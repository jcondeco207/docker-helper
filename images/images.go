package images

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/rodaine/table"
)

func DeleteImage(imageID string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	fmt.Print("Removing image ", imageID[:10], "... ")
	if _, err := cli.ImageRemove(ctx, imageID, types.ImageRemoveOptions{}); err != nil {
		panic(err)
	}
	fmt.Println("Image deleted")
}

func ShowAllImages() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	images, err := cli.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		panic(err)
	}

	imagesTable := table.New("Name", "ID", "Containers", "Size")

	for _, image := range images {
		imagesTable.AddRow(image.RepoTags[0], image.ID[:10], image.Containers, image.Size)
	}

	imagesTable.Print()
}

func GetAllImages() []types.ImageSummary {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	images, err := cli.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		panic(err)
	}

	return images
}
