FROM golang:1.19.4

WORKDIR /app

COPY gate /app/gate
COPY gate.yaml /app/gate.yaml
ENV TZ=Asia/Shanghai
RUN chmod +x /app/gate


EXPOSE 12347/tcp

ENTRYPOINT ["/app/gate"]

# docker run -v /root/code/xx/conf:/app/conf ccr.ccs.tencentyun.com/xx/gate:0.0.1


