package main

import (
	"bufio"
	"github.com/ferdoran/go-sro/framework/boot"
	"github.com/ferdoran/go-sro/framework/logging"
	"github.com/ferdoran/go-sro/framework/network"
	"github.com/ferdoran/go-sro/framework/server"
	"github.com/ferdoran/go-sro/gateway-server/clients"
	"github.com/ferdoran/go-sro/gateway-server/config"
	"github.com/ferdoran/go-sro/gateway-server/handler"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math/rand"
	"os"
	"time"
)

type GatewayServer struct {
	Server            server.Server
	failedLogins      map[string]int
	AgentServerClient *clients.AgentServerClient
}

func NewGatewayServer() GatewayServer {
	s := server.NewEngine(
		viper.GetString(config.GatewayHost),
		viper.GetInt(config.GatewayPort),
		network.EncodingOptions{
			None:         false,
			Disabled:     false,
			Encryption:   true,
			EDC:          true,
			KeyExchange:  true,
			KeyChallenge: false,
		},
	)
	s.ModuleID = viper.GetString(config.GatewayModuleId)
	boot.RegisterComponent("packethandler", handler.InitPatchRequestHandler, 2)
	boot.RegisterComponent("packethandler", handler.InitShardlistRequestHandler, 2)
	boot.RegisterComponent("packethandler", handler.InitShardlistPingHandler, 2)
	boot.RegisterComponent("packethandler", handler.InitNoticeRequestHandler, 2)
	agentClient := clients.NewAgentServerClient()
	return GatewayServer{s, make(map[string]int), agentClient}
}

func (g *GatewayServer) Start() {
	go g.Server.Start()
	g.handlePackets()
}

func (g *GatewayServer) handlePackets() {
	handler.InitLoginRequestHandler(g.failedLogins, g.AgentServerClient)

	for {
		data := <-g.Server.PacketChannel
		switch data.MessageID {
		default:
			log.Debugf("Unhandled packet %v\n", data.Packet)
		}
	}
}

func main() {
	config.Initialize()
	logging.Init()
	boot.SetPhases("packethandler", "server")
	reader := bufio.NewReader(os.Stdin)

	rand.Seed(time.Now().UnixNano())

	gw := NewGatewayServer()
	startGameServer := func() {
		gw.Start()
	}

	boot.RegisterComponent("server", startGameServer, 1)
	boot.Boot()
	log.Println("Press Enter to exit...")
	reader.ReadString('\n')
}
