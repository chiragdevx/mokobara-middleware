# Variables
BUILD_DIR = functions/productHandler
ZIP_FILE = productHandler.zip
LAMBDA_NAME = productHandler

# Cross-compiling environment variables
GOOS = linux
GOARCH = amd64

# List of directories to process
DIRECTORIES := $(shell find . -max-depth 1 -type d)

# Build the Go Lambda function
build:
	@echo "Building Go Lambda function..."
	for dir in $(DIRECTORIES); do \
		if [ -f "$$dir/main.go" ]; then \
			echo "Processing $$dir"; \
			cd $$dir && \
			GOOS=$(GOOS) GOARCH=$(GOARCH) go build -tags="lambda.norpc" -ldflags="-w -s" -o ../bootstrap main.go && \
			chmod +x ../bootstrap && \
			cd ..; \
		fi; \
	done
	@echo "Build complete."

# Zip the function code
zip: build
	@echo "Zipping Lambda function..."
	rm -f $(BUILD_DIR)/$(ZIP_FILE)
	for dir in $(DIRECTORIES); do \
		if [ -f "$$dir/main.go" ]; then \
			cd $$dir && zip -r ../../$(ZIP_FILE) bootstrap && cd ..; \
		fi; \
	done
	@echo "Zipped Lambda function."

# Deploy the function (using Terraform)
deploy: zip
	@echo "Deploying Lambda with Terraform..."
	terraform apply -auto-approve
	@echo "Deployment complete."

# Clean up old build files
clean:
	@echo "Cleaning up build files..."
	for dir in $(DIRECTORIES); do \
		if [ -f "$$dir/bootstrap" ]; then \
			rm $$dir/bootstrap; \
		fi; \
	done
	@echo "Cleaned up."
