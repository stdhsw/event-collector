IMAGE_NAME := event-collector
IMAGE_TAG  := latest

.PHONY: build

# Docker 이미지를 빌드한다
build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .
