# Build and publish Docker image for Playlister API
docker-publish:
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-t ianhowell.azurecr.io/playlister-downloader:latest \
		--push \
		.


# Sign in to Azure and log in to Azure Container Registry
acr-login:
	az login
	az acr login --name ianhowell -g default-rg
