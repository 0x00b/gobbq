package snowflake

import (
	"encoding/base64"
	"encoding/binary"
	"os"
	"sync/atomic"
	"time"

	"github.com/0x00b/gobbq/tool/host"
)

const (
	// UUID_LENGTH is length of a UUID
	UUID_LENGTH = 16
	encodeUUID  = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_."
)

var (
	// _UUIDEncoding is encoding for UUID
	_UUIDEncoding = base64.NewEncoding(encodeUUID).WithPadding(base64.NoPadding)
)

// GenUUID generates a new unique ObjectId.
func GenUUID() string {
	var b = make([]byte, 12)
	// Timestamp, 4 bytes, big endian
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	// Machine, first 3 bytes of md5(hostname)
	b[4] = machineId[0]
	b[5] = machineId[1]
	b[6] = machineId[2]
	// Pid, 2 bytes, specs don't specify endianness, but we use big endian.
	pid := os.Getpid()
	b[7] = byte(pid >> 8)
	b[8] = byte(pid)
	// Increment, 3 bytes, big endian
	i := atomic.AddUint32(&objectIdCounter, 1)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)

	return _UUIDEncoding.EncodeToString(b)
}

func GenFixedUUID(b []byte) string {
	bl := len(b)
	if bl > 12 {
		b = b[:12]
	} else if bl < 12 {
		nb := make([]byte, 12)
		copy(nb[12-bl:], b)
		b = nb
	}

	return _UUIDEncoding.EncodeToString(b)
}

// objectIdCounter is atomically incremented when generating a new ObjectId
// using NewObjectId() function. It's used as a counter part of an id.
var objectIdCounter uint32

// machineId stores machine id generated once and used in subsequent calls
// to NewObjectId function.
var machineId = host.GetHostName()
