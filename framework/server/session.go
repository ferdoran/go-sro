package server

import (
	"github.com/ferdoran/go-sro/framework/network"
	"github.com/ferdoran/go-sro/framework/security/blowfish"
	utils2 "github.com/ferdoran/go-sro/framework/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"net"
	"sync"
)

type Session struct {
	ID          string
	Conn        net.Conn
	Context     SecurityContext
	UserContext UserContext
	network.EncodingOptions
	Sequence            network.MessageSequence
	Crc                 network.MessageCRC
	Cipher              *blowfish.Cipher
	ServerPacketChannel chan PacketChannelData
	ConnectionClosed    chan *Session
	mutex               *sync.Mutex
	ServerModuleID      string
	packetChannel       chan network.Packet
}

func NewSession(
	conn net.Conn,
	options network.EncodingOptions,
	packetChannel chan PacketChannelData,
	closedConnectionChannel chan *Session,
	backendConnectionChannel chan BackendConnectionData,
	serverModuleId string,
) *Session {
	sessionId, err := uuid.NewRandom()
	if err != nil {
		log.Error("Failed to create session id")
	}

	session := Session{
		ID:                  sessionId.String(),
		Conn:                conn,
		EncodingOptions:     options,
		Context:             NewContext(),
		ServerPacketChannel: packetChannel,
		mutex:               &sync.Mutex{},
		ConnectionClosed:    closedConnectionChannel,
		ServerModuleID:      serverModuleId,
		packetChannel:       make(chan network.Packet, 1024),
	}

	return &session
}

func (s *Session) StartHandling() {
	defer s.Conn.Close()
	if !s.Context.StartedHandshake {
		s.initialize()
	}
	go func() {
		for {
			select {
			case p := <-s.packetChannel:
				s.handleIncomingPacket(p)
			}
		}
	}()
	reader := network.NewPacketReader(s.Conn)
	for {
		packets, err := reader.ReadPackets()
		if err != nil {
			if err != io.EOF {
				log.Errorln(err)
			}

			log.Debugf("Connection closed: %s\n", s.Conn.RemoteAddr().String())
			s.ConnectionClosed <- s
			break
		}

		for _, packet := range packets {
			//log.Tracef(packet.String())
			s.packetChannel <- packet
		}
	}
}

func (s *Session) handleIncomingPacket(p network.Packet) {
	// 1. Decode if Encrypted
	//log.Debugf("Received Packet: %s\n", p.String())
	pDec := p.Decrypt(s.Cipher)
	// 2. Check Sequence & CRC
	err := s.validateEDC(pDec)
	if err != nil {
		log.Errorf("Error: %v\n", err)
		return
	}
	// 3. Parse/ Handle packet
	switch queue := PacketManagerInstance.queues[pDec.MessageID]; queue {
	case nil:
		log.Tracef("Unknown packet received: \n%v\n", pDec.String())
		s.mutex.Lock()
		s.ServerPacketChannel <- PacketChannelData{Session: s, Packet: pDec}
		s.mutex.Unlock()
	default:
		queue <- PacketChannelData{
			Session: s,
			Packet:  pDec,
		}
	}
}

func (s *Session) validateEDC(p network.Packet) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	nextSeqVal := s.Sequence.Next()
	if nextSeqVal != p.Sequence {
		return errors.Errorf("Invalid Sequence. Got = %v, want %v", p.Sequence, nextSeqVal)
	}
	oldCrc := p.Buffer[network.CRCOffset]
	p.Buffer[network.CRCOffset] = 0
	nextCrcVal := s.Crc.Compute(p.Buffer[:p.DataSize+int(network.HeaderSize)])
	if nextCrcVal != oldCrc {
		err := errors.Errorf("Invalid CRC. Got = %v, want %v", p.CRC, nextCrcVal)
		return err
		//return errors.Errorf("Invalid CRC. Got = %v, want %v", p.CRC, nextCrcVal)
	}
	return nil
}

func (s *Session) initialize() {
	log.Debugln("Writing challenge")
	packet := network.EmptyPacket()
	packet.MessageID = 0x5000
	packet.WriteByte(s.EncodingOptions.GetEncodingOptionsByte())

	s.writeEncryptionSetup(&packet)
	s.writeEDC(&packet)
	s.setupKeyExchange(&packet)

	bytesToSend := packet.ToBytes()
	s.Conn.Write(bytesToSend)
	log.Debugln("Finished init")
}

func (s *Session) writeEncryptionSetup(p *network.Packet) error {
	if s.EncodingOptions.Encryption {
		cipher, err := blowfish.NewCipher(utils2.Uint64ToByteArray(s.Context.InitialBlowfishKey))
		if err != nil {
			return err
		}
		s.Cipher = cipher
		p.WriteUInt64(s.Context.InitialBlowfishKey)
		s.Context.StartedHandshake = true
	}
	return nil
}

func (s *Session) writeEDC(p *network.Packet) error {
	if s.EncodingOptions.EDC {
		sequenceSeed := byte(rand.Int31() & 0xFF)
		s.Sequence = network.NewMessageSequence(uint32(sequenceSeed))
		p.WriteUInt32(uint32(sequenceSeed))
		log.Debugf("Setup Sequence with seed %v\n", sequenceSeed)
		crc := byte(rand.Int31() & 0xFF)
		s.Crc = network.NewMessageCRC(uint32(crc))
		p.WriteUInt32(uint32(crc))

		log.Debugf("Setup CRC with seed %v\n", crc)
	}
	return nil
}

func (s *Session) setupKeyExchange(p *network.Packet) error {
	if s.KeyExchange {
		p.WriteUInt64(s.Context.HandshakeBlowfishKey)
		p.WriteUInt32(s.Context.Generator)
		p.WriteUInt32(s.Context.Prime)
		p.WriteUInt32(s.Context.LocalPublic)
	}
	return nil
}
