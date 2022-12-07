#/bin/bash

go build;mv protoc-gen-gobbq ~/go/bin/
cp -r ./gogen/grpc-go-tpl /usr/local/.butchery/