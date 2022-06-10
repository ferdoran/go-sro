package handler

import (
	"github.com/ferdoran/go-sro/framework/network"
	"github.com/ferdoran/go-sro/framework/network/opcode"
	"github.com/ferdoran/go-sro/framework/server"
	"github.com/ferdoran/go-sro/gateway-server/clients"
	"github.com/ferdoran/go-sro/gateway-server/config"
	"github.com/ferdoran/go-sro/gateway-server/db"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math/rand"
)

type LoginRequestHandler struct {
	FailedLogins      map[string]int
	AgentServerClient *clients.AgentServerClient
	channel           chan server.PacketChannelData
}

func InitLoginRequestHandler(failedLogins map[string]int, agentClient *clients.AgentServerClient) {
	queue := server.PacketManagerInstance.GetQueue(opcode.LoginRequest)
	handler := &LoginRequestHandler{FailedLogins: failedLogins, AgentServerClient: agentClient, channel: queue}
	go handler.Handle()
}

func (h *LoginRequestHandler) Handle() {
	for {
		packet := <-h.channel
		contentId, err := packet.ReadByte()
		if err != nil {
			log.Panicln("Failed to read contentId")
		}

		username, err := packet.ReadString()
		if err != nil {
			log.Panicln("Failed to read username")
		}

		password, err := packet.ReadString()
		if err != nil {
			log.Panicln("Failed to read password")
		}

		shardId, err := packet.ReadUInt16()
		if err != nil {
			log.Panicln("Failed to read shardId")
		}

		// TODO implement login
		// TODO add captcha anywhere in the future
		log.Tracef("Login Request: %d %s %d\n", contentId, username, shardId)
		p := network.EmptyPacket()
		p.MessageID = 0xA102
		if isValidLogin, accountId := db.DoLogin(username, password); isValidLogin {
			// generate agent server token
			token := h.generateAgentToken(username, password, accountId, shardId)
			// Valid login
			p.WriteByte(1)
			p.WriteUInt32(token) // AgentServer.Token
			p.WriteString(viper.GetString(config.AgentPublicIp))
			p.WriteUInt16(uint16(viper.GetInt(config.AgentPort)))
			h.FailedLogins[username] = 0

		} else {
			// TODO we need a reason why the login failed. It might be a server issue as well, at least according to doc
			// TODO Add error code 2 - Account banned
			// TODO Add error code 3 - Custom Error Message
			// Invalid login
			h.FailedLogins[username]++

			if h.FailedLogins[username] < 4 {
				p.WriteByte(2)
				p.WriteByte(1)                                           // Error Code = 1 = Counter
				p.WriteUInt32(4)                                         // Max attempts
				p.WriteUInt32(uint32(h.FailedLogins[packet.Session.ID])) // Current attempts
			} else {
				// TODO lock account for x minutes
			}
		}

		packet.Session.Conn.Write(p.ToBytes())
	}
}

func (h *LoginRequestHandler) generateAgentToken(username, password string, accountId int, shardId uint16) uint32 {
	token := rand.Uint32()
	request := clients.LoginTokenRequest{
		Username:  username,
		Password:  password,
		Token:     token,
		AccountID: uint32(accountId),
		ShardID:   shardId,
	}
	h.AgentServerClient.SendLoginToken(request)
	return token
}
