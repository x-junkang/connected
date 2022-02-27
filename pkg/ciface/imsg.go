package ciface

type IMessage interface {
	GetDataLen() uint32 //获取消息数据段长度
	GetMsgID() uint64   //获取消息ID
	GetData() []byte    //获取消息内容

	SetMsgID(uint64)   //设计消息ID
	SetData([]byte)    //设计消息内容
	SetDataLen(uint32) //设置消息数据段长度
}
