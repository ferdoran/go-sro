package handler

import (
	"github.com/ferdoran/go-sro/framework/network"
	"github.com/ferdoran/go-sro/framework/network/opcode"
	"github.com/ferdoran/go-sro/framework/server"
	"github.com/ferdoran/go-sro/gateway-server/db"
	log "github.com/sirupsen/logrus"
)

type ShardlistRequestHandler struct {
	channel chan server.PacketChannelData
}

func InitShardlistRequestHandler() {
	handler := ShardlistRequestHandler{channel: server.PacketManagerInstance.GetQueue(opcode.ShardlistRequest)}
	go handler.Handle()
}

func (h *ShardlistRequestHandler) Handle() {
	for {
		packet := <-h.channel
		log.Println("CLIENT_GATEWAY_SHARD_LIST_REQUEST")
		shards := db.GetShards()
		// TODO make this configurable
		// Farm
		p := network.EmptyPacket()
		p.MessageID = opcode.ShardlistResponse

		p.WriteByte(1)                   // hasEntries
		p.WriteByte(1)                   // Farm.ID
		p.WriteString("SRO_DEUTSCHLAND") //Farm.Name

		p.WriteByte(0) // Divider / End of Farms

		// Shards
		for _, shard := range shards {
			p.WriteByte(1)                             // hasEntries
			p.WriteUInt16(uint16(shard.Id))            // Shard.ID
			p.WriteString(shard.Name)                  // Shard.Name
			p.WriteUInt16(uint16(shard.OnlinePlayers)) // Shard.OnlinePlayers
			p.WriteUInt16(uint16(shard.Capacity))      // Shard.PlayerCapacity
			p.WriteByte(byte(shard.Status))            // Shard.IsOperating
			p.WriteByte(1)                             // Shard.FarmID
		}
		p.WriteByte(0) // Divider / End of Shards

		_, err := packet.Session.Conn.Write(p.ToBytes())
		if err != nil {
			log.Panic(err)
		}
	}
}
