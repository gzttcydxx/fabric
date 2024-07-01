# 使用官方的 Golang 镜像作为基础镜像
FROM golang:1.22-alpine as base

# 在容器内部设置工作目录
WORKDIR /app

# 安装 tzdata 包
RUN apk add --no-cache tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 首先只复制go.mod和go.sum文件
COPY go.mod .
COPY go.sum .

# 下载依赖项，利用Docker缓存层
RUN go mod download

# 使用官方的 Golang 镜像作为基础镜像
FROM api-base as release

# 现在复制整个项目
COPY . .

# 编译 Go 代码
RUN go build -o /server main.go

# 设置容器启动时运行的命令
CMD ["/server"]

# 使用官方的 Golang 镜像作为基础镜像
FROM api-base as debug

# 现在复制整个项目
COPY . .

# 设置容器启动时运行的命令
CMD ["go", "run", "main.go"]
