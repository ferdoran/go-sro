package client

import (
	network2 "github.com/ferdoran/go-sro/framework/network"
	log "github.com/sirupsen/logrus"
)

type ModuleIdentificationHandler struct {
	client *Client
}

func (h *ModuleIdentificationHandler) Handle(packet network2.Packet) {
	moduleName, err := packet.ReadString()
	if err != nil {
		log.Error("Could not read module name")
	}
	isLocalModule, err := packet.ReadByte()
	if err != nil {
		log.Error("Could not read if module is local or not")
	}

	log.Debugf("Module identified %v. Local: %v\n", moduleName, isLocalModule)
	h.client.ConnectionEstablished <- true
}
