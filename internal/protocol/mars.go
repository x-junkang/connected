package protocol

import (
	"bytes"
	"encoding/binary"
)

type MarsHeader struct {
	HeaderLength  uint32
	ClientVersion uint32
	Cmd           uint32
	Sequence      uint32
	BodyLength    uint32
}

type MarsMsg struct {
	*MarsHeader
	Opt  []byte // 可选头部字段，可能不为0
	Data []byte
}

const MarsHeaderLength = 4 * 5

func NewMarsMsg() *MarsMsg {
	return &MarsMsg{
		MarsHeader: &MarsHeader{
			HeaderLength:  MarsHeaderLength,
			ClientVersion: 1,
			Cmd:           1,
			Sequence:      0,
			BodyLength:    0,
		},
	}
}

func (mh *MarsMsg) GetDataLen() uint32 {
	return mh.BodyLength
}
func (mh *MarsMsg) GetMsgID() uint32 {
	return mh.Sequence
}
func (mh *MarsMsg) GetData() []byte {
	return mh.Data
}

func (mh *MarsMsg) GetHeaderLen() uint32 {
	return mh.HeaderLength
}

func (mh *MarsMsg) GetHeader() []byte {
	dataBuff := bytes.NewBuffer([]byte{})
	//写header
	if err := binary.Write(dataBuff, binary.LittleEndian, mh.MarsHeader); err != nil {
		return []byte{}
	}
	return dataBuff.Bytes()
}

func (mh *MarsMsg) SetMsgID(msgID uint32) {
	mh.Sequence = msgID
}
func (mh *MarsMsg) SetData(data []byte) {
	mh.Data = data
	mh.BodyLength = uint32(len(data))
}
func (mh *MarsMsg) SetDataLen(l uint32) {
	mh.BodyLength = l
}
