package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/x-junkang/connected/internal/config"
	"github.com/x-junkang/connected/pkg/ciface"
)

//DataPack 封包拆包类实例，暂时不需要成员
type DataPack struct{}

//NewDataPack 封包拆包实例初始化方法
func NewDataPack() ciface.Packet {
	return &DataPack{}
}

//GetHeadLen 获取包头长度方法
func (dp *DataPack) GetHeadLen() uint32 {
	return MarsHeaderLength
}

//Pack 封包方法(压缩数据)
func (dp *DataPack) Pack(msg ciface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//写header
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetHeader()); err != nil {
		return nil, err
	}

	//写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

//Unpack 拆包方法(解压数据)
func (dp *DataPack) Unpack(binaryData []byte) (ciface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head的信息，得到dataLen和msgID
	msg := &MarsMsg{}

	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.MarsHeader); err != nil {
		return nil, err
	}

	//判断dataLen的长度是否超出我们允许的最大包长度
	if config.GlobalObject.MaxPacketSize > 0 && msg.BodyLength > config.GlobalObject.MaxPacketSize {
		return nil, errors.New("too large msg data received")
	}

	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil
}
