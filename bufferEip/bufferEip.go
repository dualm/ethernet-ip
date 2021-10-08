package bufferEip

import (
	"bytes"
	"encoding/binary"
	"io"

	"sync"
)

type BufferEip struct {
	buffer *bytes.Buffer
	err    error
}

func (b *BufferEip) WriteLittle(target interface{}) {
	b.err = binary.Write(b.buffer, binary.LittleEndian, target)
}

func (b *BufferEip) WriteBig(target interface{}) {
	b.err = binary.Write(b.buffer, binary.LittleEndian, target)
}

func (b *BufferEip) ReadBig(target interface{}) {
	b.err = binary.Read(b.buffer, binary.BigEndian, target)
}

func (b *BufferEip) ReadLittle(target interface{}) {
	b.err = binary.Read(b.buffer, binary.LittleEndian, target)
}

func (b *BufferEip) Reset() {
	b.buffer.Reset()
}

func (b *BufferEip) Error() error {
	if b.err == io.EOF{
		return nil
	}

	return b.err
}

func (b BufferEip) Bytes() []byte {
	return b.buffer.Bytes()
}

func (b BufferEip) Len() int {
	return b.buffer.Len()
}

func (b *BufferEip)Put() {
	eipPool.Put(b)
}


func New(data []byte) *BufferEip{
	buf := bytes.NewBuffer(data)

	return &BufferEip{
		buffer: buf,
		err:    nil,
	}
}

var eipPool = sync.Pool{
	New: func() interface{} {
		return &BufferEip{
			buffer: bytes.NewBuffer(nil),
			err:    nil,
		}
	},
}
