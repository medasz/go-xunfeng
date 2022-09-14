# Go语言环境基础镜像
FROM golang:1.20.7 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

# 指定工作目录
WORKDIR /app

# 将源码拷贝到镜像中
COPY . .

# 编译镜像时，运行 go build 编译生成 server 程序
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o xunfeng

# 运行阶段指定scratch作为基础镜像
FROM scratch

WORKDIR /app

# 将上一个阶段publish文件夹下的所有文件复制进来
COPY --from=builder /app/publish .

# 为了防止代码中请求https链接报错，我们需要将证书纳入到scratch中
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/cert

# 指定运行时环境变量
ENV GIN_MODE=release \
    PORT=80

EXPOSE 80

ENTRYPOINT ["./toc-generator"]