SHELL = /bin/zsh

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

clean: 
	@docker-compose down
	@docker system prune -a;
	@docker volume prune -a;

up: 
	@docker-compose -f docker-compose.yml up -d --build

health.check: # return 42
	@curl -v http://127.0.0.1:9999

stress:
	@./stress-test/executar-teste-local.sh

docker.build: 
	@docker build -t raissageek/rinha_backend:v2 .

docker.push:
	@docker push raissageek/rinha_backend:v2
