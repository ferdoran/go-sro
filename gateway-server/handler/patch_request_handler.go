package handler

import (
	"github.com/ferdoran/go-sro/framework/network"
	"github.com/ferdoran/go-sro/framework/network/opcode"
	"github.com/ferdoran/go-sro/framework/server"
)

type PatchRequestHandler struct {
	channel chan server.PacketChannelData
}

func InitPatchRequestHandler() {
	handler := PatchRequestHandler{channel: server.PacketManagerInstance.GetQueue(opcode.PatchRequest)}
	go handler.Handle()
}

func (h *PatchRequestHandler) Handle() {
	for {
		packet := <-h.channel
		// TODO an actual check for patches
		// TODO add patch error codes
		//  InvalidVersion = 1,
		//  UPDATE = 2,
		//  NotInService = 3,
		//  AbnormalModule = 4,
		//  PatchDisabled = 5,
		p := network.EmptyPacket()
		p.MessageID = 0x600D
		p.WriteByte(1)
		p.WriteUInt16(1)
		p.WriteUInt16(opcode.PatchResponse)

		packet.Session.Conn.Write(p.ToBytes())
		p = network.EmptyPacket()
		p.MessageID = 0x600D
		p.WriteByte(0)
		p.WriteUInt16(1)

		packet.Session.Conn.Write(p.ToBytes())
	}
}
