.PHONY: all clean

# Define the base image name
IMAGE_BASE_NAME := itspeetah/npta-simple-calls

# List of functions
FUNCTIONS := a b c w2 rw d

# Default target: build all images
buildall: $(addprefix build-,$(FUNCTIONS))

# Target to build a specific function's Docker image
build-%:
	@echo "Building Docker image for function $*..."
	docker build --target function-$* -t $(IMAGE_BASE_NAME)-$*:latest .
	@echo "Successfully built $(IMAGE_BASE_NAME)-$*:latest"

# Clean up: remove all built Docker images
clean:
	@echo "Removing Docker images..."
	@for func in $(FUNCTIONS); do \
		if docker images -q $(IMAGE_BASE_NAME)-$$func:latest &> /dev/null; then \
			docker rmi $(IMAGE_BASE_NAME)-$$func:latest; \
			echo "Removed $(IMAGE_BASE_NAME)-$$func:latest"; \
		else \
			echo "Image $(IMAGE_BASE_NAME)-$$func:latest does not exist."; \
		fi \
	done
	@echo "Clean up complete."

deployall : $(addprefix deploy-,$(FUNCTIONS))

deploy-%:
	@echo "Deploying Docker imgage for function $*..."
	docker image push $(IMAGE_BASE_NAME)-$*:latest
	@echo "Successfully pushed $(IMAGE_BASE_NAME)-$*:latest"