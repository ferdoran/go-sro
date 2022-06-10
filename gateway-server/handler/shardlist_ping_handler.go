package handler

import (
	"github.com/ferdoran/go-sro/framework/network"
	"github.com/ferdoran/go-sro/framework/network/opcode"
	"github.com/ferdoran/go-sro/framework/server"
	"github.com/ferdoran/go-sro/framework/utils"
	"github.com/ferdoran/go-sro/gateway-server/config"
	"github.com/spf13/viper"
	"net"
)

type ShardlistPingHandler struct {
	channel chan server.PacketChannelData
}

func InitShardlistPingHandler() {
	handler := ShardlistPingHandler{channel: server.PacketManagerInstance.GetQueue(opcode.ShardlistPing)}
	go handler.Handle()
}

func (h *ShardlistPingHandler) Handle() {
	for {
		packet := <-h.channel
		p := network.EmptyPacket()
		p.MessageID = opcode.ShardlistPong
		// TODO can also return an error code
		p.WriteByte(1) // result = 1 = Successful, else error
		p.WriteByte(1) // result Farm.ID

		p.WriteUInt32(utils.ByteArrayToUint32(net.ParseIP(viper.GetString(config.AgentPublicIp))[12:16])) // Farm.IP

		packet.Session.Conn.Write(p.ToBytes())
	}
}
