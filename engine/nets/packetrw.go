package nets

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"sync"

	"github.com/0x00b/gobbq/engine/codec"
	"github.com/0x00b/gobbq/erro"
	"github.com/0x00b/gobbq/xlog"
	"github.com/pkg/errors"
)

var (
	errPacketBodyTooLarge = errors.Errorf("packetBody too large")
	// errChecksumError      = errors.Errorf("checksum error")
)

type Config struct {
	// flags Flags
}

func DefaultConfig() *Config {
	return &Config{
		// flags: 0,
	}
}

// PacketReadWriter is a connection that send and receive data packets upon a network stream connection
type PacketReadWriter struct {
	Config Config

	rwMtx *sync.Mutex
	conn  *Conn
}

// NewPacketReadWriter creates a packet connection based on network connection
func NewPacketReadWriter(conn *Conn) *PacketReadWriter {
	return NewPacketReadWriterWithConfig(conn, DefaultConfig())
}

func NewPacketReadWriterWithConfig(conn *Conn, cfg *Config) *PacketReadWriter {
	if conn == nil {
		panic(fmt.Errorf("conn is nil"))
	}

	if cfg == nil {
		cfg = DefaultConfig()
	}

	pc := &PacketReadWriter{
		conn:   conn,
		Config: *cfg,
		rwMtx:  &sync.Mutex{},
	}

	return pc
}

// recv receives the next packet
func (pc *PacketReadWriter) ReadPacket() (*Packet, error) {
	var err error

	var tempBuff [4]byte

	// xlog.Traceln("recv raw 1 ")
	_, err = io.ReadFull(pc.conn.rwc, tempBuff[:])
	// xlog.Traceln("recv raw 2 ")
	if err != nil {
		return nil, err
	}

	packetDataSize := packetEndian.Uint32(tempBuff[:])
	if packetDataSize > MaxPacketBodyLength {
		return nil, errPacketBodyTooLarge
	}
	if packetDataSize < 4 {
		return nil, errors.New("bad packet")
	}

	// allocate a packet to receive packetBody
	packet := NewPacket()
	packet.Src = pc.conn
	packet.totalLen = packetDataSize

	// xlog.Traceln("recv raw 3 ")
	//extendPacketBody 返回的时候已经把header buff排除了
	packetData := packet.extendPacketData(packetDataSize)
	// xlog.Traceln("recv raw 4 ")
	_, err = io.ReadFull(pc.conn.rwc, packetData)
	if err != nil {
		packet.Release()
		return nil, err
	}

	// xlog.Traceln("recv raw 5 ", packetDataSize, cap(packetData))

	packet.headerLen = packetEndian.Uint32(packetData[:4])

	// xlog.Traceln("recv raw:", packet.headerLen, string(packetData))

	// header, headerlen包含自己本身的长度（4个字节），后面才是真正的header内容
	err = codec.DefaultCodec.Unmarshal(packetData[4:packet.headerLen], packet.Header)
	if err != nil {
		packet.Release()
		return nil, err
	}

	xlog.Traceln("recv raw:", packet.String())

	// receive checksum (uint32)
	// if HasFlags(packet.Header.CheckFlags, FlagDataChecksumIEEE) {
	// 	_, err = io.ReadFull(pc.rw, tempBuff[:4])
	// 	if err != nil {
	// 		release()
	// 		return nil, nil, err
	// 	}

	// 	packetBodyCrc := crc32.ChecksumIEEE(packet.Data())
	// 	if packetBodyCrc != packetEndian.Uint32(tempBuff[:4]) {
	// 		release()
	// 		return nil, nil, errChecksumError
	// 	}
	// }

	return packet, nil
}

func writeFull(conn net.Conn, data []byte) error {
	left := len(data)
	for left > 0 {
		n, err := conn.Write(data)
		if n == left && err == nil { // handle most common case first
			return nil
		}

		if n > 0 {
			data = data[n:]
			left -= n
		}

		if err != nil {
			if !erro.IsTemporary(err) {
				return err
			} else {
				runtime.Gosched()
			}
		}
	}
	return nil
}
