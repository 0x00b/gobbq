#/bin/bash
 
protoc -I. -I../../proto/bbq --go_out=paths=source_relative:. exam.proto

# tpl_dir绝对路径
protoc -I. -I../../proto/bbq --gobbq_out=plugins=grpc,tpl_dir=/data/home/user00/.gobbq/grpc-go-tpl:. exam.proto

# npm install ts-proto
# export  PATH=$PATH:/root/node_modules/.bin
protoc -I. -I../../proto/bbq  --ts_proto_out=. --ts_proto_opt=outputServices=generic-definitions,outputClientImpl=false,oneof=unions,snakeToCamel=false,useExactTypes=false ./exam.proto
