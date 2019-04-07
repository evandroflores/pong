
check-env:
ifndef GOPATH
	@echo "Couldn't find the GOPATH env"
	@exit 1
endif
ifndef DATABASE_URL
	@echo "Couldn't find the DATABASE_URL env"
	@exit 1
endif
ifndef PONG_TOKEN
	@echo "Couldn't find the PONG_TOKEN env"
endif

run: check-env
	@go run main.go

build: check-env
	@go build -o bin/pong

check-dep:
ifeq "$(shell command -v dep)" ""
	@echo "dep is not available please install https://golang.github.io/dep/"
	@exit 1
endif

vendor: check-dep
ifeq "$(wildcard Gopkg.toml)" ""
	@echo "Initializing dep..."
	@dep init
endif
	@dep ensure
	@dep status

lint:
	@golangci-lint run

linter-install:
	@go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

test:
	@go test -gcflags=-l ./... -coverprofile coverage.txt

cover: test
	@go tool cover -func coverage.txt

opencover:
	@go tool cover -html coverage.txt
