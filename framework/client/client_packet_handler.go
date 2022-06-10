package client

import "github.com/ferdoran/go-sro/framework/network"

type PacketHandler interface {
	Handle(packet network.Packet)
}
