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
	ErrMessageBodyTooLarge = io.ErrShortWrite
	ErrMessageBodyTooSmall = io.ErrUnexpectedEOF

	errMessageBodyTooLarge = errors.Errorf("messageBody too large")
	errChecksumError       = errors.Errorf("checksum error")
)

type Config struct {
	// flags Flags
}

func DefaultConfig() *Config {
	return &Config{
		// flags: 0,
	}
}

// MessageReadWriter is a connection that send and receive data messages upon a network stream connection
type MessageReadWriter struct {
	Config Config

	rw io.ReadWriter

	headerBuff [messageHeaderSize]byte

	// write will inc 1
	writeMsgCnt uint32
	// read will inc 1
	readMsgCnt uint32
}

// NewMessageReadWriter creates a message connection based on network connection
func NewMessageReadWriter(ctx context.Context, rw io.ReadWriter) *MessageReadWriter {
	return NewMessageReadWriterWithConfig(ctx, rw, DefaultConfig())
}

func NewMessageReadWriterWithConfig(ctx context.Context, rw io.ReadWriter, cfg *Config) *MessageReadWriter {
	if rw == nil {
		panic(fmt.Errorf("conn is nil"))
	}

	if cfg == nil {
		cfg = DefaultConfig()
	}

	pc := &MessageReadWriter{
		rw:     rw,
		Config: *cfg,
	}

	return pc
}

// WriteMessage write message data to pc.rw, need to initialize the message by yourself
func (pc *MessageReadWriter) WriteMessage(message *Message) error {
	pdata := message.Data()
	err := writeFull(pc.rw, pdata)
	if err != nil {
		return err
	}

	if message.GetMessageFlags().Has(FlagDataChecksumIEEE) {
		var crc32Buffer [4]byte
		messageBodyCrc := crc32.ChecksumIEEE(pdata)
		messageEndian.PutUint32(crc32Buffer[:], messageBodyCrc)
		return writeFull(pc.rw, crc32Buffer[:])
	} else {
		return nil
	}
}

// recv receives the next message
func (pc *MessageReadWriter) ReadMessage() (*Message, error) {
	var err error

	_, err = io.ReadFull(pc.rw, pc.headerBuff[:])
	if err != nil {
		return nil, err
	}

	messageBodySize := messageEndian.Uint32(pc.headerBuff[:4])
	if messageBodySize > MaxMessageBodyLength {
		return nil, errMessageBodyTooLarge
	}

	// allocate a message to receive messageBody
	message := NewMessage()
	message.SetMessageType(MessageType(pc.headerBuff[4]))
	message.SetMessageFlags(Flags(pc.headerBuff[5]))
	message.SetMessageID(messageEndian.Uint32(pc.headerBuff[6:]) & (1<<31 - 1))
	message.Src = pc

	messageParser, ok := messageParsers[message.GetMessageType()]
	if ok {
		n, err := messageParser(pc, message)
		if err != nil {
			return nil, err
		}
		if n > messageBodySize {
			return nil, errors.New("invalid message")
		}
		messageBodySize -= n
	}

	//extendMessageBody 返回的时候已经把header buff排除了
	messageBody := message.extendMessageBody(int(messageBodySize))
	_, err = io.ReadFull(pc.rw, messageBody)
	if err != nil {
		return nil, err
	}

	message.setMessageBodyLen(messageBodySize)
	pc.readMsgCnt++

	// receive checksum (uint32)
	if message.GetMessageFlags().Has(FlagDataChecksumIEEE) {
		_, err = io.ReadFull(pc.rw, pc.headerBuff[:4])
		if err != nil {
			return nil, err
		}

		messageBodyCrc := crc32.ChecksumIEEE(message.Data())
		if messageBodyCrc != messageEndian.Uint32(pc.headerBuff[:4]) {
			return nil, errChecksumError
		}
	}

	return message, nil
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
