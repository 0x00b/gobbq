#/bin/bash

go build;mv protoc-gen-gobbq ~/go/bin/
cp -r ./gogen/bbs-go-tpl ~/.gobbq/

cp -r ./gogen/bbq-ts-client-tpl ~/.gobbq/