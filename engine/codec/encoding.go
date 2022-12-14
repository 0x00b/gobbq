// Package encoding defines the interface for the compressor and codec, and
// functions to register and retrieve compressors and codecs.
//
// # Experimental
//
// Notice: This package is EXPERIMENTAL and may be changed or removed in a
// later release.
package codec

import (
	"io"

	"github.com/0x00b/gobbq/bbqpb"
)

// Compressor is used for compressing and decompressing when sending or
// receiving packets.
type Compressor interface {
	// Compress writes the data written to wc to w after compressing it.  If an
	// error occurs while initializing the compressor, that error is returned
	// instead.
	Compress(w io.Writer) (io.WriteCloser, error)
	// Decompress reads data from r, decompresses it, and provides the
	// uncompressed data via the returned io.Reader.  If an error occurs while
	// initializing the decompressor, that error is returned instead.
	Decompress(r io.Reader) (io.Reader, error)
	// Name is the name of the compression codec and is used to set the content
	// coding header.  The result must be static; the result cannot change
	// between calls.
	Type() bbqpb.CompressType
	// If a Compressor implements
	// DecompressedSize(compressedBytes []byte) int, gRPC will call it
	// to determine the size of the buffer allocated for the result of decompression.
	// Return -1 to indicate unknown size.
	//
	// Experimental
	//
	// Notice: This API is EXPERIMENTAL and may be changed or removed in a
	// later release.
}

var registeredCompressor = make(map[bbqpb.CompressType]Compressor)

// RegisterCompressor registers the compressor with gRPC by its name.  It can
// be activated when sending an RPC via grpc.UseCompressor().  It will be
// automatically accessed when receiving a packet based on the content coding
// header.  Servers also use it to send a response with the same encoding as
// the request.
//
// NOTE: this function must only be called during initialization time (i.e. in
// an init() function), and is not thread-safe.  If multiple Compressors are
// registered with the same name, the one registered last will take effect.
func RegisterCompressor(c Compressor) {
	registeredCompressor[c.Type()] = c
	// grpcutil.RegisteredCompressorNames = append(grpcutil.RegisteredCompressorNames, c.Name())
}

// GetCompressor returns Compressor for the given compressor name.
func GetCompressor(name bbqpb.CompressType) Compressor {
	return registeredCompressor[name]
}

// Codec defines the interface gRPC uses to encode and decode packets.  Note
// that implementations of this interface must be thread safe; a Codec's
// methods can be called from concurrent goroutines.
type Codec interface {
	// Marshal returns the wire format of v.
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal parses the wire format into v.
	Unmarshal(data []byte, v interface{}) error
	// Name returns the name of the Codec implementation. The returned string
	// will be used as part of content type in transmission.  The result must be
	// static; the result cannot change between calls.
	Type() bbqpb.ContentType
}

var registeredCodecs = make(map[bbqpb.ContentType]Codec)

// RegisterCodec registers the provided Codec for use with all gRPC clients and
// servers.
//
// The Codec will be stored and looked up by result of its Name() method, which
// should match the content-subtype of the encoding handled by the Codec.  This
// is case-insensitive, and is stored and looked up as lowercase.  If the
// result of calling Name() is an empty string, RegisterCodec will panic. See
// Content-Type on
// https://github.com/grpc/grpc/blob/master/doc/bbqpbCOL-HTTP2.md#requests for
// more details.
//
// NOTE: this function must only be called during initialization time (i.e. in
// an init() function), and is not thread-safe.  If multiple Codecs are
// registered with the same name, the one registered last will take effect.
func RegisterCodec(codec Codec) {
	if codec == nil {
		panic("cannot register a nil Codec")
	}
	registeredCodecs[codec.Type()] = codec
}

// GetCodec gets a registered Codec by content-subtype, or nil if no Codec is
// registered for the content-subtype.
//
// The content-subtype is expected to be lowercase.
func GetCodec(contentType bbqpb.ContentType) Codec {
	return registeredCodecs[contentType]
}
