export GOPATH := $(PWD)
export GOBIN := $(PWD)/bin

PACKAGES := $(shell env GOPATH=$(GOPATH) go list ./... | grep -v "home")

get:
	go get -v $(PACKAGES)

install:
	go install -v $(PACKAGES)

build:
	go build -v -o bin/bot src/bot/bot.go
