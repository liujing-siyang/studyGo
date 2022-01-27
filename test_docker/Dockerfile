#本地路径就是docker build 所在文件夹的路径

FROM golang:alpine

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 设置RUN CMD ENTRYPOINT ADD COPY指令的工作目录
# 移动到工作目录（容器目录）：绝对路径
WORKDIR /root/dockerTest/test_docker_build

# 将从本地源复制文件。表示当前目录中的文件，到 WORKDIR 定义的位置。也就是将代码复制到容器中
COPY . .

# 在容器中将复制过来的代码编译成二进制可执行文件test_docker_app
RUN go build -o test_docker_app .


# 声明服务端口
EXPOSE 8000

# 启动容器时运行的命令
CMD ["./test_docker_app"]