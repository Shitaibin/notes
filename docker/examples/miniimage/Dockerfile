# 构建阶段
FROM golang:alpine
COPY hello.go .
RUN go build hello.go

# 运行阶段镜像
FROM scratch
COPY --from=0 /go/hello .
CMD ["./hello"]