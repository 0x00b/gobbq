#/bin/bash
 
protoc -I. -I../bbq --go_out=paths=source_relative:. sys.proto

# tpl_dir绝对路径
protoc -I. -I../bbq --gobbq_out=plugins=grpc,tpl_dir=/data/home/user00/.gobbq/bbs-go-tpl:. sys.proto

# protoc -I. -I../bbq --gobbq_out=plugins=grpc,tpl_dir=/data/home/user00/.gobbq/bbq-ts-client-tpl:. sys.proto
# protoc -I. -I../bbq  --ts_proto_out=. --ts_proto_opt=outputServices=generic-definitions,outputClientImpl=false,oneof=unions,snakeToCamel=false,esModuleInterop=true,useExactTypes=false,forceLong=long ./sys.proto
