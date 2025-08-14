# Load configuration from parent directory
include ../.env
export

# Build Docker image locally
docker-build:
	docker build -t $(REGISTRY_URL)/$(REGISTRY_PROJECT)-downloader:$(IMAGE_TAG) .

# Build and publish Docker image for Playlister Downloader
docker-publish:
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-t $(REGISTRY_URL)/$(REGISTRY_PROJECT)-downloader:$(IMAGE_TAG) \
		--push \
		.

# Azure Container Registry login helper
acr-login:
	az login
	az acr login --name $(ACR_NAME) -g $(ACR_RESOURCE_GROUP)

.PHONY: docker-build docker-publish acr-login
