FROM golang:1.13 AS build-dist
ENV GOPROXY='https://mirrors.aliyun.com/goproxy'
WORKDIR /data/release
COPY . .
RUN go build

FROM centos:latest as prod
WORKDIR /data/go-websocket
COPY --from=build-dist /data/release/go-websocket ./
COPY --from=build-dist /data/release/conf /data/go-websocket/conf

EXPOSE 6000

CMD ["/data/go-websocket/go-websocket","-c","./conf/app.ini"]