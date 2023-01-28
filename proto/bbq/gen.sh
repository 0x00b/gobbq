#/bin/bash
 
protoc -I. -I../ --go_out=paths=source_relative:. bbq.proto

# tpl_dir绝对路径
protoc -I. -I../ --gobbq_out=plugins=grpc,tpl_dir=/data/home/user00/.gobbq/grpc-go-tpl:. bbq.proto
