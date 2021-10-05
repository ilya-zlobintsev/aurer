package main

import (
	"github.com/docker/docker/client"
	"github.com/ilyazzz/aurer/internal"
	"github.com/ilyazzz/aurer/web"
)

func main() {
	docker, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		panic(err)
	}

	c := internal.CreateController(docker)

	web := web.InitWeb(&c)

	go c.ListenToSignals()

	// go internal.StartServices(&c)

	web.Run()
}
