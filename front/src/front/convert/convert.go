package convert

import (
	"front/errstack"
	"reflect"

	"encoding/binary"

	"bytes"
	"fmt"

	"github.com/ugorji/go/codec"
)

type Converter struct {
	cmdId      uint32
	unpackData []byte
	packedData []byte
	isPacked   bool
	Message    string
}

func (c *Converter) PackedData() []byte {
	return c.packedData
}

func (c *Converter) IsPacked() bool {
	return c.isPacked
}

func (c *Converter) CommandId() uint32 {
	return c.cmdId
}

func Create(recieve []byte) *Converter {
	cmdId, unpackData := filter(recieve)
	return &Converter{
		cmdId:      cmdId,
		unpackData: unpackData,
	}
}

/*
  Message UnPacking

  argument
    data : UnPackしたいデータ
    out  : データ出力先

  return
    失敗時にエラー
*/
func (c *Converter) Unpack(out interface{}) errstack.Stacker {
	ew := errstack.NewErrWriter()

	// decode(codec)
	mh := &codec.MsgpackHandle{RawToString: true}
	dec := codec.NewDecoderBytes(c.unpackData, mh)
	e := dec.Decode(out)

	if e != nil {
		return ew.Write(e)
	}

	return ew
}

/*
  Message Packing

  argument
    data : Packingしたいデータ

  return
    Packed Data、エラー
*/
func (c *Converter) Pack(responseCmdId int, data interface{}) errstack.Stacker {
	ew := errstack.NewErrWriter()
	var packedData []byte

	// encode(codec)
	mh := &codec.MsgpackHandle{}
	mh.MapType = reflect.TypeOf(data)
	encoder := codec.NewEncoderBytes(&packedData, mh)
	e := encoder.Encode(data)
	if e != nil {
		return ew.Write(e)
	}

	// command_id covert int to byte
	buf := new(bytes.Buffer)
	e = binary.Write(buf, binary.LittleEndian, uint32(responseCmdId))
	if e != nil {
		fmt.Println("binary.Write failed:", e)
		return ew.Write(e)
	}
	result := append(buf.Bytes(), packedData...)

	fmt.Println("lenth : ", len(buf.Bytes()))
	fmt.Println("res lenth : ", len(result))

	c.packedData = result
	c.isPacked = true
	return ew
}

func filter(raw []byte) (uint32, []byte) {
	// todo : datasize error
	// length
	cmdByteLen := 4
	// dataLen := len(raw)

	//
	cmd := raw[:cmdByteLen]
	data := raw[cmdByteLen:]

	// convert to command id
	cmdId := binary.LittleEndian.Uint32(cmd)
	return cmdId, data
}
