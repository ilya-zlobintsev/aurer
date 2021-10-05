package internal

import (
	"bufio"
	"context"
	"errors"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

type Worker struct {
	ContainerId string
	Package     string
}

func (c *Controller) RunWorker(w Worker) error {
	log.Printf("Running worker %v", w.ContainerId)

	ctx := context.Background()

	err := c.docker.ContainerStart(ctx, w.ContainerId, types.ContainerStartOptions{})

	if err != nil {
		return nil
	}

	attach, err := c.docker.ContainerAttach(ctx, w.ContainerId, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  false,
		Stdout: true,
		Stderr: true,
	})

	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(attach.Reader)

	for scanner.Scan() {
		log.Printf("%v: %v", w.ContainerId, scanner.Text())
	}

	statusCh, errCh := c.docker.ContainerWait(ctx, w.ContainerId, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case status := <-statusCh:
		if status.Error != nil {
			return errors.New(status.Error.Message)
		}

		log.Printf("Worker %v finished with status %v", w.ContainerId, status.StatusCode)

	}

	return nil
}
