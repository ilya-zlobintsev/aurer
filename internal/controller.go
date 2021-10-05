package internal

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Jguer/aur"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/ilyazzz/aurer/internal/repo"
)

type Controller struct {
	docker    *client.Client
	workers   []Worker
	aurClient *aur.Client
	RepoDir   string
}

func CreateController(docker *client.Client) Controller {
	workers := make([]Worker, 0)

	aurClient, err := aur.NewClient()

	if err != nil {
		panic(err)
	}

	return Controller{
		docker:    docker,
		workers:   workers,
		aurClient: aurClient,
		RepoDir:   "/tmp/out",
	}
}

func (c *Controller) CreateWorker(pkgname string, outdir string) (Worker, error) {
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
					TmpfsOptions: &mount.TmpfsOptions{
						Mode: 1777,
					},
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

	err := c.docker.ContainerRemove(ctx, worker.ContainerId, types.ContainerRemoveOptions{
		Force: true,
	})

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

func (c *Controller) BuildPackage(pkgname string, force ...bool) error {
	ctx := context.Background()

	for _, worker := range c.workers {
		if worker.Package == pkgname {
			return errors.New("worker already present for the given package")
		}
	}

	pkgInfo, err := c.aurClient.Info(ctx, []string{pkgname})

	if err != nil {
		return err
	}

	if len(pkgInfo) == 0 {
		return errors.New("package not found")
	}

	if !(len(force) > 0 && force[0]) {

		db, err := repo.ReadRepo(c.RepoDir)

		if err != nil {
			log.Print("Unable to read package database")
		} else {
			for _, pkg := range db {
				if pkg.Name == pkgname {
					if pkg.Version == pkgInfo[0].Version {
						log.Printf("Package %v already exists and is up-to-date, not building", pkgname)

						return nil
					}
				}
			}
		}

	}

	worker, err := c.CreateWorker(pkgname, c.RepoDir)

	if err != nil {
		return err
	}

	defer c.RemoveWorker(worker)

	err = c.RunWorker(worker)

	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) ListenToSignals() {
	sigs := make(chan os.Signal, 2)

	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	<-sigs

	go func() {
		<-sigs

		log.Print("caught second SIGTERM, terminating instantly! Warning: there may be leftover containers")
		os.Exit(1)
	}()

	wg := new(sync.WaitGroup)

	for _, worker := range c.workers {
		wg.Add(1)
		go func(w Worker) {
			log.Printf("stopping container %v", w.ContainerId)

			timeout, _ := time.ParseDuration("5s")

			err := c.docker.ContainerStop(context.Background(), w.ContainerId, &timeout)

			if err != nil {
				log.Printf("error removing worker: %v", err)
			}

			wg.Done()
		}(worker)
	}

	wg.Wait()

	os.Exit(0)
}
