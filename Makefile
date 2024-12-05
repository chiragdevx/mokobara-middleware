# Variables
BUILD_DIR = functions/productHandler
ZIP_FILE = productHandler.zip
LAMBDA_NAME = productHandler

# Cross-compiling environment variables
GOOS = linux
GOARCH = amd64

# Build the Go Lambda function
build:
	@echo "Building Go Lambda function..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -tags="lambda.norpc" -ldflags="-w -s" -o $(BUILD_DIR)/bootstrap $(BUILD_DIR)/mai	@echo "Build complete."

# Make the binary executable
chmod:
	@echo "Making binary executable..."
	chmod +x $(BUILD_DIR)/bootstrap
	@echo "Binary is executable."

# Zip the function code
zip: build chmod
	@echo "Zipping Lambda function..."
	# Remove any old zip if exists
	rm -f $(BUILD_DIR)/$(ZIP_FILE)
	# Zip the binary to the root of the zip archive
	cd $(BUILD_DIR) && zip -r ../../$(ZIP_FILE) bootstrap
	@echo "Zipped Lambda function."

# Deploy the function (using Terraform)
deploy: zip
	@echo "Deploying Lambda with Terraform..."
	terraform apply -auto-approve
	@echo "Deployment complete."

# Clean up old build files
clean:
	@echo "Cleaning up build files..."
	rm -f $(BUILD_DIR)/bootstrap $(BUILD_DIR)/$(ZIP_FILE)
	@echo "Cleaned up."
