package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobCodec struct {
	conn io.ReadWriteCloser //由构建函数传入，通常是通过 TCP 或者 Unix 建立 socket 时得到的链接实例
	buf  *bufio.Writer      //buf 是为了防止阻塞而创建的带缓冲的 Writer
	dec  *gob.Decoder       //dec 和 enc 对应 gob 的 Decoder 和 Encoder
	enc  *gob.Encoder
}

func (c GobCodec) ReadHeader(header *Header) error {
	return c.dec.Decode(header)
}

func (c GobCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c GobCodec) Write(header *Header, body interface{}) (err error) {
	defer func() {
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()

	if err := c.enc.Encode(header); err != nil {
		log.Println("rpc codec: gob error encoding header:", err)
		return err
	}

	if err := c.enc.Encode(body); err != nil {
		log.Println("rpc codec: gob error encoding body:", err)
		return err
	}
	return nil
}

func (c *GobCodec) Close() error {
	return c.conn.Close()
}

var _ Codec = (*GobCodec)(nil)

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}
