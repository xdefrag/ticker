# Check if we need to prepend docker commands with sudo
SUDO := $(shell docker version >/dev/null 2>&1 || echo "sudo")

# If TAG is not provided set default value
TAG ?= stellar/ticker:$(shell git rev-parse --short HEAD)$(and $(shell git status -s),-dirty-$(shell id -u -n))
# https://github.com/opencontainers/image-spec/blob/master/annotations.md
BUILD_DATE := $(shell date --rfc-3339=seconds)

docker-build:
	cd ../../ && \
	$(SUDO) docker build --label org.opencontainers.image.created="$(BUILD_DATE)" \
	-f services/ticker/docker/Dockerfile -t $(TAG) .

docker-push:
	cd ../../ && \
	$(SUDO) docker push $(TAG)

docker-build-dev:
	$(SUDO) docker build --label org.opencontainers.image.created="$(BUILD_DATE)" \
	-f docker/Dockerfile-dev -t $(TAG) .
