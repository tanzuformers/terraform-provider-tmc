TEST?=$$(go list ./... |grep -v 'examples')
TESTTIMEOUT=180m

.EXPORT_ALL_VARIABLES:
  TF_SCHEMA_PANIC_ON_ERROR=1

default: build

build: fmt generate
	go install

fmt:
	@echo "==> Fixing source code with gofmt..."
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

generate:
	go generate ./...

testacc: fmt
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout $(TESTTIMEOUT) -ldflags="-X=github.com/tanzuformers/terraform-provider-tmc/version.ProviderVersion=acc"

.PHONY: build fmt generate testacc
