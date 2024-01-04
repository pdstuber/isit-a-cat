.DEFAULT_GOAL := build

IMAGE ?= isit-a-cat:latest

.PHONY: build
build:
	@docker buildx create --use --name=crossplat --node=crossplat && \
	docker buildx build \
		--progress auto \
		--output "type=docker,push=false" \
		--tag ${IMAGE} \
		--file build/Dockerfile \
		.

.PHONY: run
run:
	docker run -e MODEL_PATH="/model" -e TELEGRAM_BOT_TOKEN ${IMAGE} run