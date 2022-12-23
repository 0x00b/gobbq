FROM ccr.ccs.tencentyun.com/aeras/aeras-ubuntu:v3

# docker build . -t csighub.tencentyun.com/genejli/butchery:v0.0.1
# docker run -it  -v /your/app/path:/app csighub.tencentyun.com/genejli/butchery:v0.0.1  bbq new your/app/name

# USER root
WORKDIR /tmp
# install protobuf from source
RUN apt-get update && \
    apt-get -y install git unzip build-essential autoconf libtool wget curl
RUN git clone https://github.com/google/protobuf.git && \
    cd protobuf && \
    ./autogen.sh && \
    ./configure && \
    make && \
    make install && \
    ldconfig && \
    make clean
    # cd .. && 
    # rm -r protobuf

# go
RUN wget https://dl.google.com/go/go1.17.1.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.17.1.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin:/root/go/bin

# 准备模版 
# bbq install 
# cp /usr/local/.butchery .butchery -r
COPY .butchery /usr/local/.butchery

# 生成代码工具
COPY bbq/bbq /usr/local/go/bin/
COPY protoc-gen-gobbq/protoc-gen-gobbq /usr/local/go/bin/

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install github.com/envoyproxy/protoc-gen-validate@latest
# RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# 安装代码检查工具 
RUN wget https://mirrors.tencent.com/repository/generic/vlib/tools/gonote/release/linux/latest/gonote && cp gonote /root/go/bin
RUN curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin 
RUN	go install golang.org/x/tools/cmd/goimports@latest

# python3
RUN apt-get -y install python3 python3-pip
# 免ssh交互
RUN apt-get -y install expect

WORKDIR /app

CMD [ "bbq" ]

