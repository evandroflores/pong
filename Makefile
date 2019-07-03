export GO111MODULE=on

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
	@echo "[Makefile] DATABASE_URL OK"
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
