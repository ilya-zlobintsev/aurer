package internal

import (
	"context"
	"log"
	"time"
	"github.com/docker/docker/api/types"
)

func StartServices(c *Controller) {
	log.Println("Starting services")

	pullWorkerImage(c)

	go runPullCron(c)
}

func runPullCron(c *Controller) {
	interval := time.Hour * 24

	log.Printf("Worker image pull interval: %s", interval)

	for {
		time.Sleep(interval)

		pullWorkerImage(c)
	}
}

func pullWorkerImage(c *Controller) {
	ctx := context.Background()

	log.Println("Pulling worker image")

	_, err := c.docker.ImagePull(ctx, "docker.io/ilyazzz/aur-cd-worker", types.ImagePullOptions{})

	if err != nil {
		log.Panicf("Error pulling worker image! %v", err)
	}

	log.Println("Worker image updated")
}
