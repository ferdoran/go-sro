package server

import (
	"github.com/ferdoran/go-sro/framework/network/opcode"
	log "github.com/sirupsen/logrus"
)

type BackendConnectionData struct {
	*Session
	ModuleID string
}

type BackendConnectionHandler struct {
	BackendConnected chan BackendConnectionData
	backendModules   map[string]string
	channel          chan PacketChannelData
}

func InitBackendConnectionHandler(backendConnectedChannel chan BackendConnectionData, backedModules map[string]string) {
	handler := &BackendConnectionHandler{
		BackendConnected: backendConnectedChannel,
		backendModules:   backedModules,
		channel:          PacketManagerInstance.GetQueue(opcode.BackendAuthentication),
	}

	go handler.Handle()
}

func (h *BackendConnectionHandler) Handle() {
	for {
		packet := <-h.channel
		serverModuleId, err := packet.ReadString()
		if err != nil {
			log.Error("Could not read server name")
		}

		secret, err := packet.ReadString()
		if err != nil {
			log.Error("Could not read secret")
		}

		if moduleSecret, exists := h.backendModules[serverModuleId]; exists {
			if moduleSecret == secret {
				log.Infof("%s connected", serverModuleId)
				h.BackendConnected <- BackendConnectionData{
					Session:  packet.Session,
					ModuleID: serverModuleId,
				}
			} else {
				log.Warnf("wrong secret for %s: %s. closing connection", serverModuleId, secret)
				packet.Session.Conn.Close()
			}
		} else {
			log.Warnf("unknown backend module tried to connect: %s. closing connection", serverModuleId)
			packet.Session.Conn.Close()
		}
	}
}
