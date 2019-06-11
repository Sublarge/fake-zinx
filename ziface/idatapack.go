package ziface

type IDataPack interface {
	GetHeaderLen() uint32
	Pack(msg IMessage) ([]byte, error)
	Unpack(data []byte) (IMessage, error)
}
