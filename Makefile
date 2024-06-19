BUILD_ENVPARMS:=GOOS=linux GOARCH=amd64 CGO_ENABLED=0
PWD=$(shell pwd)
BUILD_TS:=$(shell date +%FT%T%z)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
GIT_COMMIT=$(shell git rev-parse --short HEAD)
HEAD_COMMIT=$(shell git rev-parse HEAD)
APP_VERSION?=$(shell git name-rev --tags --name-only ${HEAD_COMMIT})

LDFLAGS:=-X 'github.com/kshvakov/techleadconf/kinectl.ReleaseDate=$(BUILD_TS)'\
		 -X 'github.com/kshvakov/techleadconf/kinectl.GitCommit=$(GIT_COMMIT)'\
		 -X 'github.com/kshvakov/techleadconf/kinectl.GitBranch=$(GIT_BRANCH)'\
		 -X 'github.com/kshvakov/techleadconf/kinectl.AppVersion=$(APP_VERSION)'\

build:
	@[ -d .build ] || mkdir -p .build
	@$(BUILD_ENVPARMS) go build -ldflags "-s -w $(LDFLAGS)" -o .build/kinectl kinectl/cmd/kinectl/main.go
	@file  .build/kinectl
	@du -h .build/kinectl

deb: build # dogfooding

ifneq ($(APP_VERSION), undefined)
	@APP_VERSION=${APP_VERSION} go run kinectl/cmd/kinectl/main.go deb
	@dpkg -I .build/kinectl_*.deb
else
	echo "SKIP package version is ${APP_VERSION}"
endif

example-run-server:
	@go run examples/service/cmd/server/main.go

example-service-deb: build
	@$(BUILD_ENVPARMS) go build -ldflags "-s -w $(LDFLAGS)" -o .build/example examples/service/cmd/server/main.go
	@file  .build/example
	@du -h .build/example
	@APP_VERSION=v1.2.3 .build/kinectl deb --spec examples/service/spec.yml
	@dpkg -I .build/example_*.deb

example-service-play: example-service-deb
	@.build/kinectl play --spec examples/service/spec.yml --env dev --ops-dir  examples/ansible --debug

.PHONY: build