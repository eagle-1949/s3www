PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)


all: build

# Builds minio locally.
build: checks
	@echo "Building s3www binary to './s3www'"
	@CGO_ENABLED=0 go build  -o $(PWD)/s3www

# Builds minio and installs it to $GOPATH/bin.
install: build
	@echo "Installing s3www binary to '$(GOPATH)/bin/s3www'"
	@mkdir -p $(GOPATH)/bin && cp $(PWD)/s3www $(GOPATH)/bin/s3www
	@echo "Installation successful. To learn more, try \"s3www --help\"."

clean:
	@echo "Cleaning up all the generated files"
	@find . -name '*.test' | xargs rm -fv
	@rm -rvf s3www
	@rm -rvf build
	@rm -rvf release
