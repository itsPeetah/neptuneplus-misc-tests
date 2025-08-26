IMAGE_BASE_NAME := itspeetah/np-prime-numbers-caller
IMAGE_TAG := latest

publish: build push
	

build:
	@echo "Building docker image..."
	docker buildx build --platform linux/amd64,linux/arm64 --target function -t $(IMAGE_BASE_NAME):$(IMAGE_TAG) .
	@echo "Done."

clean:
	@echo "Removing Docker images..."
	if docker images -q $(IMAGE_BASE_NAME):$(IMAGE_TAG) &> /dev/null; then \
		docker rmi $(IMAGE_BASE_NAME):$(IMAGE_TAG); \
		echo "Removed $(IMAGE_BASE_NAME):$(IMAGE_TAG)"; \
	else
		echo "Image $(IMAGE_BASE_NAME):$(IMAGE_TAG) does not exist."; \
	fi
	@echo "Clean up complete"

push:
	@echo "Pushing Docker imgage for function prime-numbers-caller..."
	docker image push $(IMAGE_BASE_NAME):$(IMAGE_TAG)
	@echo "Done."

deploy:
	@echo "Deploying functions..."
	kubectl apply -f ./config
	@echo "Done"

undeploy:
	@echo "Undeploying functions..."
	kubectl delete -f ./config
	@echo "Done"

deploy-%:
	@echo "Deploying functions..."
	kubectl apply -f ./config/$*
	@echo "Done"

undeploy-%:
	@echo "Deploying functions..."
	kubectl delete -f ./config/$*
	@echo "Done"