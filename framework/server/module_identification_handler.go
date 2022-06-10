package server

import (
	"github.com/ferdoran/go-sro/framework/network"
	log "github.com/sirupsen/logrus"
)

type ModuleIdentificationHandler struct {
	channel chan PacketChannelData
}

func InitModuleIdentificationHandler() {
	queue := PacketManagerInstance.GetQueue(0x2001)
	handler := ModuleIdentificationHandler{channel: queue}
	go handler.Handle()
}

func (mih *ModuleIdentificationHandler) Handle() {
	for {
		packet := <-mih.channel
		moduleName, err := packet.ReadString()
		if err != nil {
			log.Error("Could not read module name")
		}
		isLocalModule, err := packet.ReadByte()
		if err != nil {
			log.Error("Could not read if module is local or not")
		}

		log.Debugf("Module identified %v. Local: %v", moduleName, isLocalModule)

		p := network.EmptyPacket()
		p.MessageID = 0x2001
		p.WriteString(packet.Session.ServerModuleID)
		p.WriteByte(0)
		packet.Session.Conn.Write(p.ToBytes())
	}
}
