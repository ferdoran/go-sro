package clients

import (
	"github.com/ferdoran/go-sro/framework/client"
	"github.com/ferdoran/go-sro/framework/network"
	"github.com/ferdoran/go-sro/framework/network/opcode"
	"github.com/ferdoran/go-sro/gateway-server/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type AgentServerClient struct {
	PacketChannel chan network.Packet
	Connected     bool
	*client.BackendClient
}

type LoginTokenRequest struct {
	Username  string
	Password  string
	Token     uint32
	AccountID uint32
	ShardID   uint16
}

func NewAgentServerClient() *AgentServerClient {
	c := client.NewBackendClient(
		viper.GetString(config.AgentHost),
		viper.GetInt(config.AgentPort),
		viper.GetString(config.GatewayModuleId),
		viper.GetString(config.GatewaySecret),
	)
	logrus.Infoln("connecting to agent server")
	c.Connect()
	ac := AgentServerClient{PacketChannel: c.IncomingPacketChannel, BackendClient: c}
	return &ac
}

func (ac *AgentServerClient) SendLoginToken(request LoginTokenRequest) {
	p := network.EmptyPacket()
	p.MessageID = opcode.GatewayLoginTokenRequest
	p.WriteUInt32(request.AccountID)
	p.WriteString(request.Username)
	p.WriteString(request.Password)
	p.WriteUInt32(request.Token)
	p.WriteUInt16(request.ShardID)
	ac.OutgoingPacketChannel <- p
}
