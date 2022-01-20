server-image:
	docker build -t docker.io/ilyazzz/aurer . -f images/server.Dockerfile
worker-image:
	docker build -t docker.io/ilyazzz/aurer-worker . -f images/worker.Dockerfile