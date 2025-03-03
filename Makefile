.PHONY: build install test testacc clean examples

default: build

# Binary name
BINARY=terraform-provider-openai
VERSION=5.0.0
OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)
PROVIDER_PATH=~/.terraform.d/plugins/registry.terraform.io/happypathway/openai/${VERSION}/${OS_ARCH}

build:
	go build -o ${BINARY}

install: build
	mkdir -p ${PROVIDER_PATH}
	cp ${BINARY} ${PROVIDER_PATH}/

uninstall:
	rm -rf ${PROVIDER_PATH}

test:
	go test ./...

testacc:
	TF_ACC=1 go test ./internal/... -v

clean:
	go clean
	rm -f ${BINARY}

examples: install
	@echo "Running examples..."
	@for dir in $(shell ls -d examples/*/); do \
		echo "Applying examples in $$dir"; \
		cd $$dir && terraform init && terraform apply -auto-approve; \
		if [ $$? -ne 0 ]; then \
			echo "Error applying example in $$dir"; \
			exit 1; \
		fi; \
		cd -; \
	done

destroy-examples:
	@echo "Destroying examples..."
	@for dir in $(shell ls -d examples/*/); do \
		echo "Destroying resources in $$dir"; \
		cd $$dir && terraform destroy -auto-approve; \
		if [ $$? -ne 0 ]; then \
			echo "Error destroying example in $$dir"; \
			exit 1; \
		fi; \
		cd -; \
	done

lint:
	golangci-lint run --fix

fmt:
	go fmt ./...