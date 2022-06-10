package server

type KeepAliveHandler struct {
	channel chan PacketChannelData
}

func InitKeepAliveHandler() {
	queue := PacketManagerInstance.GetQueue(0x2002)
	handler := KeepAliveHandler{channel: queue}
	go handler.Handle()
}

func (h KeepAliveHandler) Handle() {
	// Do Nothing
}
