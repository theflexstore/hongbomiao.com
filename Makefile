# Docker
docker-build:
	docker build --file=web/Dockerfile --tag=hm-web .
	docker build --file=api/Dockerfile.development --tag=hm-api-dev .
	docker build --file=api/Dockerfile --tag=hm-api .
	docker build --file=api-go/Dockerfile.api.development --tag=hm-api-go-dev .
	docker build --file=api-go/Dockerfile.api --tag=hm-api-go .
	docker build --file=api-go/Dockerfile.grpc.development --tag=hm-api-go-grpc-dev .
	docker build --file=api-go/Dockerfile.grpc --tag=hm-api-go-grpc .

docker-run:
	docker run -p 80:80 web
	docker run -p 5000:5000 --name=hm_api_dev --rm --env-file=./api/.env.development.local.example.docker hm-api-dev
	docker run -p 5000:5000 --name=hm_api --rm --env-file=./api/.env.production.local.example hm-api
	docker run -p 5000:5000 --name=hm_api_go_dev --rm hm-api-go-dev
	docker run -p 5000:5000 --name=hm_api_go --rm --env=APP_ENV=production hm-api-go
	docker run -p 5000:5000 --name=hm_api_go_grpc_dev --rm hm-api-go-grpc-dev
	docker run -p 5000:5000 --name=hm_api_go_grpc --rm --env=APP_ENV=production --env=GRPC_HOST=0.0.0.0 hm-api-go-grpc

docker-sh:
	docker run  --rm -it hm-api-go-grpc sh

docker-ps:
	docker ps
	docker ps --all

docker-push:
	docker system prune

docker-prune:
	docker system prune

# Docker Compose
docker-compose-build:
	docker-compose --file=docker-compose.development.yml build
	docker-compose --file=docker-compose.cypress.yml build

docker-compose-up:
	docker-compose --file=docker-compose.development.yml up --detach
	docker-compose --file=docker-compose.cypress.yml up --detach

docker-compose-stop:
	docker-compose --file=docker-compose.development.yml stop
	docker-compose --file=docker-compose.cypress.yml stop

docker-compose-down:
	docker-compose --file=docker-compose.development.yml down --volumes
	docker-compose --file=docker-compose.cypress.yml down --volumes

# Kubernetes
minikube-start:
	minikube start

minikube-service-web:
	minikube service web-service

minikube-service-api-go:
	minikube service api-go-service

minikube-dashboard:
	minikube dashboard

minikube-delete:
	minikube delete

minikube-ip:
	minikube ip

kubectl-apply:
	kubectl apply -f kubernetes

kubectl-delete:
	kubectl delete -f kubernetes

kubectl-pods:
	kubectl get pods

kubectl-logs:
	kubectl logs pod_name --follow

kubectl-services:
	kubectl get services

kubectl-namespaces:
	kubectl get namespaces

skaffold:
	skaffold dev

# Prometheus
prom-curl:
	curl http://localhost:9464/metrics

prom-test:
	docker build --file=Dockerfile.prometheus.test .
