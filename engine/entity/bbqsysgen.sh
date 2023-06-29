#/bin/bash
 
protoc -I. -I../../proto/bbq --go_out=paths=source_relative:. bbqsys.proto

# tpl_dir绝对路径
protoc -I. -I../../proto/bbq --gobbq_out=plugins=grpc,tpl_dir=/data/home/user00/.gobbq/bbs-go-tpl:. bbqsys.proto

# protoc -I. -I../bbq --gobbq_out=plugins=grpc,tpl_dir=/data/home/user00/.gobbq/bbq-ts-client-tpl:. sys.proto
# protoc -I. -I../bbq  --ts_proto_out=. --ts_proto_opt=outputServices=generic-definitions,outputClientImpl=false,oneof=unions,snakeToCamel=false,esModuleInterop=true,useExactTypes=false,forceLong=long ./sys.proto

sed -i "s/entity\.//g" bbqsys.bbq.go
sed -i "s/\"github.com\/0x00b\/gobbq\/engine\/entity\"//g" bbqsys.bbq.go
