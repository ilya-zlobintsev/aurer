package internal

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Jguer/aur"
	"github.com/Jguer/go-alpm/v2"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/ilyazzz/aurer/internal/repo"
)

type StatusMsg struct {
	Status   []string
	Workers  []Worker
	Packages []repo.PkgInfo
}

type Controller struct {
	docker      *client.Client
	workers     []Worker
	aurClient   *aur.Client
	RepoDir     string
	status      []string
	StatusChans []chan StatusMsg
}

func isDocker() bool {
	_, err := os.Stat("/.dockerenv")

	if err != nil {
		_, err = os.Stat("/run/.containerenv") // For Podman

		return err == nil
	} else {
		return true
	}
}

func CreateController(docker *client.Client) Controller {
	aurClient, err := aur.NewClient()

	if err != nil {
		panic(err)
	}

	var repoDir string

	if isDocker() {
		repoDir = "/repo"
	} else {
		repoDir = os.Getenv("AURER_REPO_DIR")

		if repoDir == "" {
			log.Println("Not running in Docker and AURER_REPO_DIR is not set! Aborting.")
			os.Exit(1)
		}
	}

	log.Printf("Serving repo from %v", repoDir)

	return Controller{
		docker:      docker,
		workers:     make([]Worker, 0),
		aurClient:   aurClient,
		RepoDir:     repoDir,
		status:      make([]string, 0),
		StatusChans: make([]chan StatusMsg, 0),
	}
}

func (c *Controller) addStatus(status string) {
	c.status = append(c.status, status)
	go c.updateStatus()
}

func (c *Controller) removeStatus(status string) {
	for i, w := range c.status {
		if w == status {
			c.status[i] = c.status[len(c.status)-1]

			c.status = c.status[:len(c.status)-1]

			break
		}
	}

	go c.updateStatus()
}

func (c *Controller) updateStatus() {
	for _, ch := range c.StatusChans {
		ch <- c.GetStatus()
	}
}

func (c *Controller) GetStatus() StatusMsg {
	repo, err := repo.ReadRepo(c.RepoDir)

	if err != nil {
		log.Printf("Failed to read repo: %v", repo)
	}
	return StatusMsg{
		Status:   c.status,
		Workers:  c.workers,
		Packages: repo,
	}
}

const DEFAULT_WORKER_IMAGE = "ghcr.io/ilyazzz/aurer-worker"

func (c *Controller) getWorkerImage() string {

	image := os.Getenv("AURER_WORKER_IMAGE")

	if image == "" {
		image = DEFAULT_WORKER_IMAGE
	}

	return image
}

func (c *Controller) CreateWorker(pkgname string, outdir string) (Worker, error) {
	ctx := context.Background()

	container, err := c.docker.ContainerCreate(
		ctx,
		&container.Config{
			Image: c.getWorkerImage(),
			Labels: map[string]string{
				"aurer.worker": "1",
			},
			Env: []string{"PACKAGE=" + pkgname},
		},
		&container.HostConfig{
			Tmpfs: map[string]string{
				"/work": "rw,exec,nosuid,nodev",
			},
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
					Target: "/repo",
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

	go c.updateStatus()

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
	go c.updateStatus()

	return err
}

func (c *Controller) GetWorkers() []Worker {
	return c.workers
}

func (c *Controller) BuildPackage(pkgname string, force ...bool) error {
	status := fmt.Sprintf("Building %v", pkgname)

	c.addStatus(status)
	defer c.removeStatus(status)

	for _, worker := range c.workers {
		if worker.Package == pkgname {
			return errors.New("worker already present for the given package")
		}
	}

	outdir := c.getRepoPath()

	log.Printf("Using %v as output folder", outdir)

	worker, err := c.CreateWorker(pkgname, outdir)

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

func (c *Controller) Update() error {
	status := "Updating packages"

	c.addStatus(status)
	defer c.removeStatus(status)

	ctx := context.Background()

	db, err := repo.ReadRepo(c.RepoDir)

	if err != nil {
		return err
	}

	pkgs := make(map[string]string) // Package name and version
	var pkgNames []string

	for _, pkg := range db {
		pkgs[pkg.Name] = pkg.Version
		pkgNames = append(pkgNames, pkg.Name)
	}

	info, err := c.aurClient.Info(ctx, pkgNames)

	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, pkgInfo := range info {
		oldVer := pkgs[pkgInfo.Name]
		newVer := pkgInfo.Version

		if alpm.VerCmp(newVer, oldVer) > 0 {

			log.Printf("Detected package update for %v: (%v -> %v)", pkgInfo.Name, oldVer, newVer)
			wg.Add(1)
			go func(name string) {
				defer wg.Done()

				c.BuildPackage(name)
			}(pkgInfo.Name)
		}

	}

	wg.Wait()

	go c.updateStatus()

	return nil
}

func (c *Controller) getRepoPath() string {
	if isDocker() {
		path := c.getMountHostPath("/repo")

		if path == "" {
			log.Printf("Running in Docker and /repo is not mounted! Cannot proceed.")
			os.Exit(1)

			return ""
		} else {
			return path
		}
	} else {
		return c.RepoDir
	}
}

func (c *Controller) getMountHostPath(dest string) string {
	listFilters := filters.NewArgs(filters.Arg("label", "aurer.server=1"))

	containers, err := c.docker.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: listFilters,
	})

	if err != nil {
		log.Panicf("Failed to list containers! %v", err)
	}

	if len(containers) == 0 {
		log.Panicf("Running in Docker and the server container cannot be found!")
	}

	container := containers[0]

	for _, mount := range container.Mounts {
		if mount.Destination == dest {
			return mount.Source
		}
	}

	return ""
}
