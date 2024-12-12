.PHONY: build_all clean_all

LAMBDA_FUNCTIONS := productHandler

build:
	@for function in $(LAMBDA_FUNCTIONS); do \
		echo "Building $$function..."; \
		cd functions/$$function && GOOS=linux GOARCH=amd64 go build -o bootstrap . && cd - > /dev/null; \
	done

clean:
	@echo "Cleaning up all build artifacts..."
	@for function in $(LAMBDA_FUNCTIONS); do \
		echo "Cleaning $$function..."; \
		cd functions/$$function && rm -f bootstrap && cd - > /dev/null; \
	done
	rm -f functions/*.zip

deploy: build
	terraform apply -auto-approve