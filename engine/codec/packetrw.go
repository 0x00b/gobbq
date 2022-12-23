package codec

import (
	"context"
	"fmt"
	"hash/crc32"
	"io"
	"runtime"

	"github.com/0x00b/gobbq/erro"
	"github.com/pkg/errors"
)

var (
	ErrPacketBodyTooLarge = io.ErrShortWrite
	ErrPacketBodyTooSmall = io.ErrUnexpectedEOF

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

	rw io.ReadWriter

	headerBuff [packetHeaderSize]byte

	// write will inc 1
	writeMsgCnt uint32
	// read will inc 1
	readMsgCnt uint32
}

// NewPacketReadWriter creates a packet connection based on network connection
func NewPacketReadWriter(ctx context.Context, rw io.ReadWriter) *PacketReadWriter {
	return NewPacketReadWriterWithConfig(ctx, rw, DefaultConfig())
}

func NewPacketReadWriterWithConfig(ctx context.Context, rw io.ReadWriter, cfg *Config) *PacketReadWriter {
	if rw == nil {
		panic(fmt.Errorf("conn is nil"))
	}

	if cfg == nil {
		cfg = DefaultConfig()
	}

	pc := &PacketReadWriter{
		rw:     rw,
		Config: *cfg,
	}

	return pc
}

// WritePacket write packet data to pc.rw, need to initialize the packet by yourself
func (pc *PacketReadWriter) WritePacket(packet *Packet) error {
	pdata := packet.Data()
	err := writeFull(pc.rw, pdata)
	if err != nil {
		return err
	}

	if packet.GetPacketFlags().Has(FlagDataChecksumIEEE) {
		var crc32Buffer [4]byte
		packetBodyCrc := crc32.ChecksumIEEE(pdata)
		packetEndian.PutUint32(crc32Buffer[:], packetBodyCrc)
		return writeFull(pc.rw, crc32Buffer[:])
	} else {
		return nil
	}
}

// recv receives the next packet
func (pc *PacketReadWriter) ReadPacket() (*Packet, error) {
	var err error

	_, err = io.ReadFull(pc.rw, pc.headerBuff[:])
	if err != nil {
		return nil, err
	}

	packetBodySize := packetEndian.Uint32(pc.headerBuff[:4])
	if packetBodySize > MaxPacketBodyLength {
		return nil, errPacketBodyTooLarge
	}

	// allocate a packet to receive packetBody
	packet := NewPacket()
	packet.SetPacketType(PacketType(pc.headerBuff[4]))
	packet.SetPacketFlags(Flags(pc.headerBuff[5]))
	packet.setMsgHeaderLen(packetEndian.Uint32(pc.headerBuff[6:]))
	packet.Src = pc

	//extendPacketBody 返回的时候已经把header buff排除了
	packetBody := packet.extendPacketBody(int(packetBodySize))
	_, err = io.ReadFull(pc.rw, packetBody)
	if err != nil {
		return nil, err
	}

	packet.setPacketBodyLen(packetBodySize)
	pc.readMsgCnt++

	// receive checksum (uint32)
	if packet.GetPacketFlags().Has(FlagDataChecksumIEEE) {
		_, err = io.ReadFull(pc.rw, pc.headerBuff[:4])
		if err != nil {
			return nil, err
		}

		packetBodyCrc := crc32.ChecksumIEEE(packet.Data())
		if packetBodyCrc != packetEndian.Uint32(pc.headerBuff[:4]) {
			return nil, errChecksumError
		}
	}

	return packet, nil
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
