#/bin/bash
 
protoc -I. -I../../../proto/bbq --go_out=paths=source_relative:. gate.proto

# tpl_dir绝对路径
protoc -I. -I../../../proto/bbq --gobbq_out=plugins=grpc,tpl_dir=/data/home/user00/.gobbq/grpc-go-tpl:. gate.proto
