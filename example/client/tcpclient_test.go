package main

import (
	"testing"
)

func TestTcpClient(m *testing.T) {

	// wsc, err := net.Dial("tcp", ":1234")
	// if err != nil {
	// 	panic(err)
	// }
	// ws := codec.NewPacketReadWriter(wsc)

	// pkt, release := codec.NewPacket()
	// defer release()

	// pkt.WriteBody([]byte("dsfsdfs"))
	// ws.WritePacket(pkt)

	// if pkt, release, err = ws.ReadPacket(); err != nil {
	// 	log.Fatal(err)
	// }
	// defer release()

	// fmt.Printf("Received: %s.\n", string(pkt.PacketBody()))
}
