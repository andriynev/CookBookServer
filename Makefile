ifneq ($(wildcard /etc/ssl/certs/ca-certificates.crt),)
	export CURL_CA_BUNDLE=/etc/ssl/certs/ca-certificates.crt
endif

ifeq ($(shell which go),)
	export PATH := /usr/local/go/bin:$(PATH)
endif

ARCH=amd64

dependencies:
	go mod download
	go mod tidy

build-app:
	go build -o $(PWD)/bin/api src/api

run:
	go run $(PWD)/src/api/main.go
