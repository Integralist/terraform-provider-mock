dependencies:
	@echo "Download go.mod dependencies"
	@go mod download

install-tools: dependencies
	@echo "Installing tools from tools/tools.go"
	@cat tools/tools.go | grep _ | awk -F '"' '{print $$2}' | xargs -tI {} go install {}

generate-docs: install-tools
	tfplugindocs generate

validate-docs: install-tools
	tfplugindocs validate

build:
	go build -o terraform-provider-mock
