# 使用官方的 Golang 镜像作为基础镜像
FROM golang:1.22-alpine as release

# 在容器内部设置工作目录
WORKDIR /app

# 安装 tzdata 包
RUN apk add --no-cache tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 将服务器的 Go 代码复制到容器中
COPY . .

# 安装 Gin
RUN go mod tidy

# 编译 Go 代码
RUN go build -o /server main.go

# 设置容器启动时运行的命令
CMD ["/server"]

# 使用官方的 Golang 镜像作为基础镜像
FROM golang:1.22-alpine as debug

# 在容器内部设置工作目录
WORKDIR /app

# 安装 tzdata 包
RUN apk add --no-cache tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 将服务器的 Go 代码复制到容器中
COPY . .

# 安装 Gin
RUN go mod tidy

# 设置容器启动时运行的命令
CMD ["go", "run", "main.go"]
