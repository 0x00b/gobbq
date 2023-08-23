FROM golang:1.19.4

WORKDIR /app

COPY proxy /app/proxy
RUN chmod +x /app/proxy

COPY proxy.yaml /app/proxy.yaml

ENV TZ=Asia/Shanghai

EXPOSE 12347/tcp

ENTRYPOINT ["/app/proxy"]

# docker run -v /root/code/xx/conf:/app/conf ccr.ccs.tencentyun.com/xx/proxy:0.0.1

# docker build -t --network=host csighub.tencentyun.com/xxx/xxx:v1 .