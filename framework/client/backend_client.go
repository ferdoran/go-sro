package client

import (
	"github.com/ferdoran/go-sro/framework/network"
	"github.com/ferdoran/go-sro/framework/network/opcode"
)

type BackendClient struct {
	*Client
	BackendConnection chan string
	Connected         bool
	secret            string
}

func NewBackendClient(host string, port int, moduleId, secret string) *BackendClient {
	client := NewClient(host, port, moduleId)
	backendClient := &BackendClient{
		Client:    client,
		Connected: false,
		secret:    secret,
	}
	backendClient.AutoReconnect = true
	return backendClient
}

func (bc *BackendClient) Connect() {
	bc.Client.Connect()
	bc.BackendAuthentication()
	go func() {
		for {
			select {
			case <-bc.Reconnected:
				bc.BackendAuthentication()
			}
		}
	}()
}

func (bc *BackendClient) BackendAuthentication() {
	p := network.EmptyPacket()
	p.MessageID = opcode.BackendAuthentication
	p.WriteString(bc.ModuleID)
	p.WriteString(bc.secret)
	bc.Client.OutgoingPacketChannel <- p
	bc.Connected = true
}
