package main

import (
	"github.com/docker/docker/client"
	"github.com/ilyazzz/aurer/internal"
	"github.com/ilyazzz/aurer/web"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	docker, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		panic(err)
	}

	c := internal.CreateController(docker)

	web := web.InitWeb(&c)

	go internal.StartServices(&c)

	web.Run()
}
