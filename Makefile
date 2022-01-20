server-image:
	docker build -t ghcr.io/ilyazzz/aurer . -f images/server.Dockerfile
worker-image:
	docker build -t ghcr.io/ilyazzz/aurer-worker . -f images/worker.Dockerfile