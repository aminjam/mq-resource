GOTOOLS = github.com/FiloSottile/gvt \
	github.com/onsi/ginkgo/ginkgo
PACKAGES=$(shell go list ./... | grep -v vendor | sort | uniq)
DOCKER_IMAGE="aminjam/mq-resource"
PLUGIN=""

all: format check-env update-deps test build docker

build:
		@echo "--> Building mq-resource"
		@GOARCH=amd64 GOOS=linux go build -o built-check check/cmd/main.go
		@GOARCH=amd64 GOOS=linux go build -o built-in in/cmd/main.go
		@GOARCH=amd64 GOOS=linux go build -o built-out out/cmd/main.go

check-env:
ifeq (${GO15VENDOREXPERIMENT},)
	@echo "ERR: Use Go >1.5 and set GO15VENDOREXPERIMENT flag (source .envrc)"
	@exit 1
endif

docker:
	@echo "--> Docker build and push"
	@docker build -t ${DOCKER_IMAGE} ${PWD}
	@docker push ${DOCKER_IMAGE}

format:
		@echo "--> Running go fmt"
		@go fmt $(PACKAGES)

test:
		@echo "--> Running Tests"
		@${PWD}/scripts/integration-tests.sh ${PLUGIN}

update-deps:
		@echo "--> Updating dependencies"
		@go get -v $(GOTOOLS)
		@gvt update --all

.PHONY: all build check-env docker format test update-deps
