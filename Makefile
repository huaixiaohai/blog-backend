GO = go
# >>>>>>>>>>> 自定义常量 >>>>>>>>>>>>>
# 定义项目基本信息
CURRENT_GIT_REPO  := blog
COMMONENVVAR      ?= GOOS=linux GOARCH=amd64
BUILDENVVAR       ?= CGO_ENABLED=0
TARGET_DIR        ?= $(CURDIR)/target

# >>>>>>>>>>> 必须包含的命令 >>>>>>>>>
# 定义环境变量
export GOBIN  := $(CURDIR)/bin

# 构建并编译出静态可执行文件
all: linux_build

# 单元测试
test:
	@echo ">> Go test module for server"
	go test -v $(CURDIR)/handler

# 生成可执行文件
build:
	# go build -o $(TARGET_DIR)/grpc-visit
	go build -o $(GOBIN)/$(CURRENT_GIT_REPO)

# 交叉编译出linux下的静态可执行文件
linux_build:
	$(COMMONENVVAR) $(BUILDENVVAR) make build

lint:
	golangci-lint run

# 编译Docker
docker: linux_build
	@echo ">> build docker artifact"
	docker build --rm --no-cache -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) $(CURDIR)

# push编译成功的docker到仓库中
docker_push: docker
	@echo ">> docker push $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) $(CURDIR)"
	docker push $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

# 清除所有编译生成的文件
clean:
	@rm -rf bin target

.PHONY: build linux_build lint all test clean docker docker_push
