#/bin/bash

# cd ~/gooooo/gobbq/cmd/bbq/proto/protoc-gen-gobbq; ./build.sh 
cd ~/gooooo/gobbq/components/proxy/proxypb;./gen.sh 
cd ~/gooooo/gobbq/components/gate/gatepb;./gen.sh 
cd ~/gooooo/gobbq/engine/entity;./bbqsysgen.sh
cd ~/gooooo/gobbq/frame/frameproto;./gen.sh 
cd ~/gooooo/gobbq/frame/test/testpb;./gen.sh 
cd ~/gooooo/gobbq/example/exampb;./gen.sh 