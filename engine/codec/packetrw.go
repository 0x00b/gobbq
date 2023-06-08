package codec

import (
	"fmt"
	"io"
	"runtime"
	"sync"

	"github.com/0x00b/gobbq/erro"
	"github.com/0x00b/gobbq/xlog"
	"github.com/pkg/errors"
)

var (
	errPacketBodyTooLarge = errors.Errorf("packetBody too large")
	errChecksumError      = errors.Errorf("checksum error")
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
	rw    io.ReadWriter

	// write will inc 1
	writeMsgCnt uint32
	// read will inc 1
	readMsgCnt uint32
}

// NewPacketReadWriter creates a packet connection based on network connection
func NewPacketReadWriter(rw io.ReadWriter) *PacketReadWriter {
	return NewPacketReadWriterWithConfig(rw, DefaultConfig())
}

func NewPacketReadWriterWithConfig(rw io.ReadWriter, cfg *Config) *PacketReadWriter {
	if rw == nil {
		panic(fmt.Errorf("conn is nil"))
	}

	if cfg == nil {
		cfg = DefaultConfig()
	}

	pc := &PacketReadWriter{
		rw:     rw,
		Config: *cfg,
		rwMtx:  &sync.Mutex{},
	}

	return pc
}

// SendPackt write packet data to pc.rw, need to initialize the packet by yourself
func (pc *PacketReadWriter) SendPackt(packet *Packet) error {
	pdata := packet.Data()

	xlog.Traceln("send raw:", packet.String())

	pc.rwMtx.Lock()
	defer pc.rwMtx.Unlock()

	// todo 合并 ，不要分两次 writeFull, 优化不用append
	err := writeFull(pc.rw, append(packetEndian.AppendUint32(nil, uint32(len(pdata))), pdata...))
	if err != nil {
		return err
	}
	xlog.Traceln("send raw done")

	pc.writeMsgCnt++

	// if HasFlags(packet.Header.CheckFlags, FlagDataChecksumIEEE) {
	// 	var crc32Buffer [4]byte
	// 	packetBodyCrc := crc32.ChecksumIEEE(pdata)
	// 	packetEndian.PutUint32(crc32Buffer[:], packetBodyCrc)
	// 	return writeFull(pc.rw, crc32Buffer[:])
	// }
	return nil
}

// recv receives the next packet
func (pc *PacketReadWriter) ReadPacket() (*Packet, ReleasePkt, error) {
	var err error

	var tempBuff [4]byte

	xlog.Traceln("recv raw 1 ")
	_, err = io.ReadFull(pc.rw, tempBuff[:])
	xlog.Traceln("recv raw 2 ")
	if err != nil {
		return nil, nil, err
	}

	packetDataSize := packetEndian.Uint32(tempBuff[:])
	if packetDataSize > MaxPacketBodyLength {
		return nil, nil, errPacketBodyTooLarge
	}

	// allocate a packet to receive packetBody
	packet, release := NewPacket()
	packet.Src = pc
	packet.totalLen = packetDataSize

	// xlog.Traceln("recv raw 3 ")
	//extendPacketBody 返回的时候已经把header buff排除了
	packetData := packet.extendPacketData(packetDataSize)
	xlog.Traceln("recv raw 4 ")
	_, err = io.ReadFull(pc.rw, packetData)
	if err != nil {
		release()
		return nil, nil, err
	}

	xlog.Traceln("recv raw 5 ")

	packet.headerLen = packetEndian.Uint32(packetData[:4])

	// xlog.Traceln("recv raw:", packet.headerLen, string(packetData))

	// header, headerlen包含自己本身的长度（4个字节），后面才是真正的header内容
	err = DefaultCodec.Unmarshal(packetData[4:packet.headerLen], packet.Header)
	if err != nil {
		release()
		return nil, nil, err
	}

	xlog.Traceln("recv raw:", packet.String())

	pc.readMsgCnt++

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

	return packet, release, nil
}

func writeFull(conn io.Writer, data []byte) error {
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
