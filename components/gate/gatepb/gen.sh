#/bin/bash
 
protoc -I. -I../../../proto/bbq --go_out=paths=source_relative:. gate.proto

# tpl_dir绝对路径
protoc -I. -I../../../proto/bbq --gobbq_out=plugins=grpc,tpl_dir=/data/home/user00/.gobbq/bbs-go-tpl:. gate.proto

protoc -I. -I../../../proto/bbq --gobbq_out=plugins=grpc,tpl_dir=/data/home/user00/.gobbq/bbq-ts-client-tpl:. gate.proto
protoc -I. -I../../../proto/bbq  --ts_proto_out=. --ts_proto_opt=outputClientImpl=false,oneof=unions,snakeToCamel=false,useExactTypes=false ./gate.proto
