package client

import (
	network2 "github.com/ferdoran/go-sro/framework/network"
	log "github.com/sirupsen/logrus"
)

type MassivePacketHandler struct {
	client              *Client
	currentMessageID    uint16
	currentMessageCount uint16
	currentMessageBody  [][]byte
}

func NewMassivePacketHandler(session *Client) *MassivePacketHandler {
	handler := &MassivePacketHandler{}
	handler.client = session
	return handler
}

func (h *MassivePacketHandler) Handle(packet network2.Packet) {
	isHeader, err := packet.ReadBool()
	if err != nil {
		log.Error("Failed to read massive packet isHeader")
	}

	if isHeader {
		// HEAD
		h.currentMessageCount, err = packet.ReadUInt16()
		if err != nil {
			log.Error("Failed to read massive packet count")
		}
		h.currentMessageID, err = packet.ReadUInt16()
		if err != nil {
			log.Error("Failed to read massive packet message id")
		}
		h.currentMessageBody = make([][]byte, 0)
		log.Infof("Received massive packet isHeader with message ID %v\n", h.currentMessageID)
		return
	}
	// DATA / BODY
	if h.currentMessageID != 0 && h.currentMessageCount > 0 {
		// parse as much as you can
		h.currentMessageBody = append(h.currentMessageBody, packet.Data[1:])
		h.currentMessageCount--
		if h.currentMessageCount == 0 {
			// construct final packet
			finalPacket := h.constructFinalPacket()
			log.Tracef("Finished massive packet with message ID %v\n", h.currentMessageID)
			log.Tracef(finalPacket.String())
			h.currentMessageID = 0
			h.currentMessageCount = 0
			h.currentMessageBody = nil
			h.client.IncomingPacketChannel <- finalPacket
			//h.client.PacketHandlers[finalPacket.MessageID].Handle(finalPacket)
		}
	}
}

func (h *MassivePacketHandler) constructFinalPacket() network2.Packet {
	finalBuffer := make([]byte, 0)
	for i := 0; i < len(h.currentMessageBody); i++ {
		finalBuffer = append(finalBuffer, h.currentMessageBody[i]...)
	}
	packet := network2.EmptyPacket()
	packet.MessageID = h.currentMessageID
	packet.Buffer = append(packet.Buffer[:network2.DataOffset], finalBuffer...)
	packet.WriteMessageID()
	return network2.NewPacket(packet.Buffer)
}
