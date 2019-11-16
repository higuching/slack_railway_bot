# Goパラメータ
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=slack_railway_bot
BINARY_ARM=${BINARY_NAME}_arm

# OS情報の取得
OS_NAME := $(shell uname -s | tr A-Z a-z)

all: build
arm:
	GOOS=linux GOARCH=arm ${GOBUILD} -o ./bin/${BINARY_ARM} -v

build:
	GOOS=darwin GOARCH=amd64 ${GOBUILD} -o ./bin/${BINARY_NAME} -v

clean:
	$(GOCLEAN)
	rm -f ./bin/$(BINARY_NAME)
	rm -f ./bin/$(BINARY_ARM)
