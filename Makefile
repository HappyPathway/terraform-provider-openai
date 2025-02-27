default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Run unit tests
.PHONY: test
test:
	go test ./... -v $(TESTARGS)

# Build provider
.PHONY: build
build:
	go build -o terraform-provider-openai

# Format code
.PHONY: fmt
fmt:
	go fmt ./...
	terraform fmt -recursive

# Install provider locally for testing
.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/HappyPathway/openai/0.1.0/darwin_amd64
	mv terraform-provider-openai ~/.terraform.d/plugins/HappyPathway/openai/0.1.0/darwin_amd64/

# Clean build artifacts
.PHONY: clean
clean:
	rm -f terraform-provider-openai