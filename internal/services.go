package internal

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/docker/docker/api/types"
)

func StartServices(c *Controller) {
	log.Println("Starting services")

	c.pullWorkerImage()

	go c.runPullCron()

	go c.runUpdatesChecker()
}

func (c *Controller) runPullCron() {
	interval := time.Hour * 24

	log.Printf("Worker image pull interval: %s", interval)

	for {
		time.Sleep(interval)

		c.pullWorkerImage()
	}
}

func (c *Controller) pullWorkerImage() {
	status := "Pulling worker image"

	c.addStatus(status)
	defer c.removeStatus(status)

	ctx := context.Background()

	log.Println("Pulling worker image")

	reader, err := c.docker.ImagePull(ctx, c.getWorkerImage(), types.ImagePullOptions{})

	if err != nil {
		log.Panicf("Error pulling worker image! %v", err)
	}

	defer reader.Close()

	io.Copy(os.Stderr, reader)

	log.Println("Worker image updated")
}

func (c *Controller) runUpdatesChecker() {
	intervalStr := os.Getenv("AURER_UPDATE_INTERVAL")

	var interval time.Duration

	if intervalStr != "" {
		interval, _ = time.ParseDuration(intervalStr)

	} else {
		interval, _ = time.ParseDuration("30m")
	}

	log.Printf("Checking for updates every %v", interval)

	for {
		log.Printf("Checking for updates...")

		err := c.Update()

		if err != nil {
			log.Printf("Error updating! %v", err)
		}

		time.Sleep(interval)
	}
}
