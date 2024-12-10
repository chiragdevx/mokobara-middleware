.PHONY: build_all clean_all

LAMBDA_FUNCTIONS = productHandler

build_all:
	@for function in $(LAMBDA_FUNCTIONS); do \
		echo "Building $$function..."; \
		cd functions/$$function && GOOS=linux GOARCH=amd64 go build -o bootstrap main.go && cd - > /dev/null; \
	done

clean_all:
	@for function in $(LAMBDA_FUNCTIONS); do \
		echo "Cleaning $$function..."; \
		cd functions/$$function && rm -f main main.zip && cd - > /dev/null; \
	done
