
include .env
PROJECTNAME=$(shell basename "$(pwd)")
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
#GOPATH="$(GOBASE)/vendor:$(GOBASE)"

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## build: build default application
build:
	go build -o $(GOBIN)/main cmd/shortener/main.go

## run-config: run application with config.json
run-config:
	go run cmd/shortener/main.go -c config.json

## run-with-memory: run application with memory storage
run-with-memory:
	go run cmd/shortener/main.go

## run-with-db: run application with db postgres storage
run-with-db:
	go run cmd/shortener/main.go -d postgres://root:12345@localhost:5432/shorten

## run-with-file: run application with file storage
run-with-file:
	go run cmd/shortener/main.go -f storage.txt

## ldflags: ldflags test incremenent19 - buildVersion, buildDate, buildCommit,
ldflags:
	go run -ldflags "-X 'main.buildVersion=v1.0.1' -X 'main.buildDate=07.12.2022' -X 'main.buildCommit=test commit'" cmd/shortener/main.go

## pprof: run pprof heap profile on port 9090
pprof:
	PPROF_TMPDIR=$(GOBASE)/profiles go tool pprof -http=":9090" -seconds=30  http://localhost:8080/debug/pprof/heap

## pprof-load: add benchmark load to project
pprof-load:
	go run pprof_load.go

## pprof-diff: show diff pprof profiles
pprof-diff:
	go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof

## compile: Compiling for every OS and Platform
compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o $(GOBIN)/main-linux-arm cmd/shortener/main.go
	GOOS=linux GOARCH=arm64 go build -o $(GOBIN)/main-linux-arm64 cmd/shortener/main.go

## certificate: Generate ssl certificate
certificate:
	go run cmd/shortener/certificate.go

## proto: generate proto files
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/shortener.proto

help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
