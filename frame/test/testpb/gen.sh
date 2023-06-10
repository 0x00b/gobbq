#/bin/bash
 
protoc -I. -I../../../proto/bbq --go_out=paths=source_relative:. testpb.proto

# tpl_dir绝对路径
protoc -I. -I../../../proto/bbq --gobbq_out=plugins=grpc,tpl_dir=/data/home/user00/.gobbq/bbs-go-tpl:. testpb.proto

protoc -I. -I../../../proto/bbq --gobbq_out=plugins=grpc,tpl_dir=/data/home/user00/.gobbq/bbq-ts-client-tpl:. testpb.proto

# npm install ts-proto
# export  PATH=$PATH:/root/node_modules/.bin
protoc -I. -I../../../proto/bbq  --ts_proto_out=. --ts_proto_opt=outputClientImpl=false,oneof=unions,snakeToCamel=false,useExactTypes=false ./testpb.proto

# protoc \
#     -I. -I../../proto/bbq\
#     --plugin="protoc-gen-ts=../../node_modules/.bin/protoc-gen-ts" \
#     --js_out="import_style=commonjs,binary:." \
#     --ts_out="." \
#     exam.proto

 

# protoc \
#     -I. -I../../proto/bbq \
#     --plugin="protoc-gen-ts=../../node_modules/.bin/protoc-gen-ts" \
#     --plugin=protoc-gen-grpc=../../node_modules/.bin/grpc_tools_node_protoc_plugin \
#     --js_out="import_style=commonjs,binary:." \
#     --ts_out="service=grpc-node:." \
#     --grpc_out="." \
#     exam.proto
