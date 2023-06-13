#/bin/bash

protoc --go_out=paths=source_relative:. bbq.proto

protoc -I.  --ts_proto_out=. --ts_proto_opt=outputServices=generic-definitions,outputClientImpl=false,oneof=unions,snakeToCamel=false,esModuleInterop=true,useExactTypes=false,forceLong=long ./bbq.proto
