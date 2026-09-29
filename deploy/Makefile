ROOT_DIR := $(abspath ..)
CLOUD_DIR := $(ROOT_DIR)/cloud

IMAGE_NAME := openindustrial-cloud
CONTAINER_NAME := openindustrial-cloud
IMAGE_TAG := latest

.PHONY: build run stop restart logs clean

# ==============================
# Build
# ==============================

# 编译后端 Docker 镜像
# 前端 console 不参与构建
build:
	@echo "==> Building cloud backend image..."
	docker build \
		-t $(IMAGE_NAME):$(IMAGE_TAG) \
		$(CLOUD_DIR)
	@echo "==> Cloud backend image built: $(IMAGE_NAME):$(IMAGE_TAG)"

# ==============================
# Run
# ==============================

# 更新后端镜像并重新运行容器
run: build
	@echo "==> Removing old cloud container..."
	-docker rm -f $(CONTAINER_NAME)

	@echo "==> Starting cloud container..."
	docker run -d \
		--name $(CONTAINER_NAME) \
		-p 8080:8080 \
		$(IMAGE_NAME):$(IMAGE_TAG)

	@echo "==> Cloud container started."

# ==============================
# Container
# ==============================

stop:
	@echo "==> Stopping cloud container..."
	-docker stop $(CONTAINER_NAME)

restart:
	@$(MAKE) run

logs:
	docker logs -f $(CONTAINER_NAME)

# ==============================
# Clean
# ==============================

clean:
	@echo "==> Removing cloud container..."
	-docker rm -f $(CONTAINER_NAME)

	@echo "==> Removing cloud image..."
	-docker rmi $(IMAGE_NAME):$(IMAGE_TAG)