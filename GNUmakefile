default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	@echo "Starting Testing"
	@cd ./internal/provider
	@$env:TF_ACC="1"
	@go test -count=1 -v 

# Download any Dependencies
dependencies:
	@echo "Download go.mod Dependencies"
	@go mod download

# Get Updates and tidy mod and sum files
bump:
	@echo "Update and Tidy Dependencies"
	@go get -u ./...
	@go mod tidy

# Install all packages
install:
	@echo "Install All Packages"
	@go install

# Build the executable
build:
	@echo "Build the provider executables"
	@go build -o terraform-provider-provider.exe
	@go build -o terraform-provider-passwordstate.exe

# Generate Terraform Documentation
generate-docs:
	@echo "Generate Terraform Docs"
	@go generate ./...
