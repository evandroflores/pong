
check-env:
ifndef GOPATH
	@echo "Couldn't find the GOPATH env"
	@exit 1
endif

run: check-env
	@go run main.go

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
	@dep status