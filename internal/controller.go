package internal

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type Controller struct {
	docker  *client.Client
	workers []Worker
}

func CreateController(docker *client.Client) Controller {
	workers := make([]Worker, 0)

	return Controller{
		docker:  docker,
		workers: workers,
	}
}

func (c *Controller) createWorker(pkgname string, outdir string) (Worker, error) {
	ctx := context.Background()

	container, err := c.docker.ContainerCreate(
		ctx,
		&container.Config{
			Image: "docker.io/ilyazzz/aur-cd-worker",
			Labels: map[string]string{
				"aurer.worker": "1",
			},
			Env: []string{"PACKAGE=" + pkgname},
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: "/etc/aurer/mirrorlist",
					Target: "/etc/pacman.d/mirrorlist",
				},
				{
					Type:   mount.TypeBind,
					Source: "/etc/aurer/makepkg.conf",
					Target: "/etc/makepkg.conf",
				},
				{
					Type:   mount.TypeBind,
					Source: outdir,
					Target: "/output",
				},
				{
					Type:   mount.TypeTmpfs,
					Target: "/work",
				},
			},
		},
		nil,
		nil,
		"",
	)

	if err != nil {
		return Worker{}, err
	}

	worker := Worker{
		ContainerId: container.ID,
		Package:     pkgname,
	}

	c.workers = append(c.workers, worker)

	log.Printf("Created worker with container ID %v", container.ID)

	return worker, nil
}

func (c *Controller) RemoveWorker(worker Worker) error {
	log.Printf("Removing worker with container ID %v", worker.ContainerId)

	ctx := context.Background()

	err := c.docker.ContainerRemove(ctx, worker.ContainerId, types.ContainerRemoveOptions{})

	for i, w := range c.workers {
		if w.ContainerId == worker.ContainerId {
			c.workers[i] = c.workers[len(c.workers)-1]

			c.workers = c.workers[:len(c.workers)-1]
		}
	}

	return err
}
func (c *Controller) GetWorkers() []Worker {
	return c.workers
}
