
check-env:
ifndef GOPATH
	@echo "[Makefile] GOPATH FAIL - Environment variable not set."
	@exit 1
else
	@echo "[Makefile] GOPATH OK"
endif
ifndef DATABASE_URL
	@echo "[Makefile] DATABASE_URL FAIL - Environment variable not set."
	@exit 1
else
	@echo "[Makefile] DATABASE_URL OK - $(DATABASE_URL)"
endif
ifndef PONG_TOKEN
	@echo "[Makefile] PONG_TOKEN FAIL - Environment variable not set."
else
	@echo "[Makefile] PONG_TOKEN OK"
endif

run: check-env
	@go run main.go

build: check-env
	@go build -o bin/pong

check-dep:
ifeq "$(shell command -v dep)" ""
	@echo "[Makefile] Dependency management FAIL - Please install https://golang.github.io/dep/"
	@exit 1
else
	@echo "[Makefile] Dependency management OK"
endif

vendor: check-dep
ifeq "$(wildcard Gopkg.toml)" ""
	@echo "[Makefile] Initializing dep..."
	@dep init
else
	@echo "[Makefile] Gopkg.toml OK"
endif
	@dep ensure
	@dep status

lint:
	@golangci-lint run

linter-install:
	@go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

test: check-env
	@go test -gcflags=-l ./... -coverprofile coverage.txt

cover: test
	@go tool cover -func coverage.txt

opencover:
	@go tool cover -html coverage.txt
