package server

import (
	log "github.com/sirupsen/logrus"
)

type KeyExchangeCompletedHandler struct {
	channel chan PacketChannelData
}

func InitKeyExchangeCompletedHandler() {
	queue := PacketManagerInstance.GetQueue(0x9000)
	handler := KeyExchangeCompletedHandler{channel: queue}
	go handler.Handle()
}

func (h *KeyExchangeCompletedHandler) Handle() {
	for {
		packet := <-h.channel
		if !packet.Session.Context.StartedHandshake && !packet.Session.Context.CompletedLocalSetup && !packet.Session.Context.CompletedRemoteSetup {
			log.Error("Handshake still not setup completely")
		}
		log.Debugln("Finished Handshake")
		packet.Session.Context.FinishedHandshake = true
	}
}
