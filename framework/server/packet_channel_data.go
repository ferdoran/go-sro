package server

import "github.com/ferdoran/go-sro/framework/network"

type PacketChannelData struct {
	*Session
	network.Packet
}
