cd ~/gooooo/gobbq/components/proxy ; rm proxy.log; go build -gcflags "-l -N";./proxy
cd ~/gooooo/gobbq/example/server2; rm proxy proxy.log;cp ../../components/proxy/proxy ./ ; ./proxy
cd ~/gooooo/gobbq/components/gate; rm gate.log; go build -gcflags "-l -N"; ./gate
cd ~/gooooo/gobbq/example/server; rm server.log;go build -gcflags "-l -N";./server
cd ~/gooooo/gobbq/example/server2; rm server2.log;go build -gcflags "-l -N";./server2 
cd ~/gooooo/gobbq/example/client; rm client.log p.log; go test wsclient_test.go -v 2>&1 >p.log