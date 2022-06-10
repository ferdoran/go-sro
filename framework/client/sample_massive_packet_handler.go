package client

import (
	network2 "github.com/ferdoran/go-sro/framework/network"
	log "github.com/sirupsen/logrus"
)

type SampleMassivePacketHandler struct {
}

func (h *SampleMassivePacketHandler) Handle(packet network2.Packet) {
	str1, _ := packet.ReadString()
	str2, _ := packet.ReadString()

	log.Infof("Received 0x1234 message [%s] [%s]\n", str1, str2)
}
