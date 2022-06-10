package client

import (
	"fmt"
	"github.com/ferdoran/go-sro/framework/network"
	"github.com/ferdoran/go-sro/framework/security/blowfish"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net"
	"sync"
	"time"
)

type Client struct {
	host                      string
	port                      int
	Conn                      net.Conn
	Ctx                       SecurityContext
	EncodingOptions           network.EncodingOptions
	Cipher                    *blowfish.Cipher
	Sequence                  network.MessageSequence
	Crc                       network.MessageCRC
	ModuleID                  string
	PacketHandlers            map[uint16]PacketHandler
	IncomingPacketChannel     chan network.Packet
	OutgoingPacketChannel     chan network.Packet
	ConnectionClosed          chan bool
	ConnectionEstablished     chan bool
	Reconnected               chan bool
	mutex                     *sync.Mutex
	packetMutex               *sync.Mutex
	AutoReconnect             bool
	reconnectCount            int
	MaxReconnectAttempts      int
	HandlingOutgoingPackets   bool
	ReconnectTimeoutInSeconds time.Duration
}

func NewClient(host string, port int, moduleId string) *Client {
	c := Client{
		host:                      host,
		port:                      port,
		Conn:                      nil,
		Ctx:                       SecurityContext{},
		EncodingOptions:           network.EncodingOptions{},
		Cipher:                    nil,
		Sequence:                  network.MessageSequence{},
		Crc:                       network.MessageCRC{},
		ModuleID:                  moduleId,
		PacketHandlers:            make(map[uint16]PacketHandler),
		IncomingPacketChannel:     make(chan network.Packet, 8),
		OutgoingPacketChannel:     make(chan network.Packet, 8),
		mutex:                     &sync.Mutex{},
		packetMutex:               &sync.Mutex{},
		MaxReconnectAttempts:      100,
		ConnectionEstablished:     make(chan bool, 1),
		Reconnected:               make(chan bool, 1),
		ConnectionClosed:          make(chan bool, 1),
		ReconnectTimeoutInSeconds: 10,
	}
	rand.Seed(time.Now().UnixNano())
	c.Ctx.Private = rand.Uint32() & 0xFFFFFFFF
	c.PacketHandlers[0x5000] = &KeyExchangeHandler{&c}
	c.PacketHandlers[0x2001] = &ModuleIdentificationHandler{&c}
	c.PacketHandlers[0x600D] = NewMassivePacketHandler(&c)
	return &c
}

func (c *Client) Connect() {
	conn, err := net.Dial("tcp4", fmt.Sprintf("%s:%d", c.host, c.port))

	if err != nil {
		log.Debugf("failed to connect, reason: %s", err.Error())
		c.Reconnect()
	} else {
		log.Debugf("Connected to %v\n", conn.RemoteAddr().String())
		if c.reconnectCount > 0 {
			c.Reconnected <- true
		}
		c.reconnectCount = 0

		c.Conn = conn
		go c.startHandling()
		select {
		case successful := <-c.ConnectionEstablished:
			if successful {
				go c.handleOutgoingPackets()
				log.Debugf("Connection established")
				return
			}
		}
	}
}

func (c *Client) Reconnect() {
	if c.AutoReconnect && c.reconnectCount < c.MaxReconnectAttempts {
		c.mutex.Lock()
		c.reconnectCount++
		c.Ctx = SecurityContext{}
		c.Cipher = nil
		log.Infof("Reconnecting... Attempt [%v/%v]\n", c.reconnectCount, c.MaxReconnectAttempts)
		c.mutex.Unlock()
		time.Sleep(time.Second * c.ReconnectTimeoutInSeconds)
		c.Connect()
	} else {
		log.Infof("Reconnect disabled [%v] or maximum attempts [%v] reached. Closing connection\n", !c.AutoReconnect, c.MaxReconnectAttempts)
		c.Conn.Close()
		c.ConnectionClosed <- true
	}
}

func (c *Client) handleOutgoingPackets() {
	if c.HandlingOutgoingPackets {
		log.Debugf("Already handling outgoing packets...")
		return
	}
	c.HandlingOutgoingPackets = true
	mutex := &sync.Mutex{}
	log.Debugf("Handling outgoing packets...")
	for {
		select {
		case packet := <-c.OutgoingPacketChannel:
			mutex.Lock()
			log.Debugf("Sending message: %02X\n", packet.MessageID)
			c.sendPacket(packet)
			mutex.Unlock()
		case <-c.ConnectionClosed:
			c.HandlingOutgoingPackets = false
			log.Debugf("Stop handling outgoing packets...")
			return
		}
	}

}

func (c *Client) startHandling() {
	reader := network.NewPacketReader(c.Conn)
	for {
		packets, err := reader.ReadPackets()
		if err != nil {
			log.Errorf("Error: %v\n", err)
			c.Conn.Close()
			c.ConnectionClosed <- true
			c.Reconnect()
			break
		}

		for _, p := range packets {
			c.handleIncomingPacket(p)
		}
	}
}

func (c *Client) handleIncomingPacket(packet network.Packet) {
	packetDec := packet.Decrypt(c.Cipher)
	switch handler := c.PacketHandlers[packetDec.MessageID]; handler {
	case nil:
		c.IncomingPacketChannel <- packetDec
	default:
		handler.Handle(packetDec)
	}
}

func (c *Client) sendPacket(p network.Packet) {
	c.packetMutex.Lock()
	defer c.packetMutex.Unlock()
	p.ServerPacket = false
	p.WriteMessageID()
	p.WriteSize()

	p.Sequence = c.Sequence.Next()
	p.IsSequenceInitialized = true
	p.WriteSequence()
	p.CRC = c.Crc.Compute(p.Buffer[:p.DataSize+int(network.HeaderSize)])
	p.WriteCRC()
	p.Encrypt(c.Cipher)
	numBytes, err := c.Conn.Write(p.Buffer[:p.DataSize+int(network.HeaderSize)])
	if err != nil {
		log.Panicf("err sending msg: %v", err)
	}

	if numBytes != p.DataSize+int(network.HeaderSize) {
		log.Panicf("sent %d bytes, want = %d", numBytes, p.DataSize+int(network.HeaderSize))
	}
}
