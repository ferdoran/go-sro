package server

import (
	"encoding/binary"
	"github.com/ferdoran/go-sro/framework/network"
	"github.com/ferdoran/go-sro/framework/security"
	"github.com/ferdoran/go-sro/framework/security/blowfish"
	utils2 "github.com/ferdoran/go-sro/framework/utils"
	log "github.com/sirupsen/logrus"
)

type KeyExchangeHandler struct {
	channel chan PacketChannelData
}

func InitKeyExchangeHandler() {
	queue := PacketManagerInstance.GetQueue(0x5000)
	handler := KeyExchangeHandler{channel: queue}
	go handler.Handle()
}

func (keh *KeyExchangeHandler) Handle() {
	for {
		var err error
		packet := <-keh.channel
		if !packet.Session.Context.StartedHandshake {
			log.Error("Received handshake packet before handshake started")
		}

		// 1. Parse Handshake Response
		packet.Context.RemotePublic, err = packet.ReadUInt32()
		if err != nil {
			log.Error("Could not read remote public")
		}

		packet.Session.Context.RemoteSignature, err = packet.ReadUInt64()
		if err != nil {
			log.Error("Could not read remote challenge")
		}

		// 2. Calculate Common Secret
		packet.Session.Context.CommonSecret = security.G_pow_X_mod_P(packet.Session.Context.RemotePublic, packet.Session.Context.Private, packet.Session.Context.Prime)

		// 3. derive new blowfish key
		newBlowfishKey := security.CalculateKey(packet.Session.Context.CommonSecret, packet.Session.Context.LocalPublic, packet.Session.Context.RemotePublic)
		packet.Session.Cipher, err = blowfish.NewCipher(newBlowfishKey)
		if err != nil {
			log.Error(err)
		}

		// 4. Decrypt remote challenge
		keh.decryptRemoteChallenge(packet)

		// 5. Compute signature
		expectedSignature := utils2.ByteArrayToUint64(
			security.CalculateChallenge(packet.Session.Context.CommonSecret, packet.Session.Context.RemotePublic, packet.Session.Context.LocalPublic))

		// 6. Validate Remote signature
		if expectedSignature != packet.Session.Context.RemoteSignature {
			log.Errorf("Invalid client signature. Got = %v, want %v", packet.Session.Context.LocalSignature, packet.Session.Context.RemoteSignature)
		}

		// 7. Calculate Local Challenge
		keh.calculateNewLocalChallenge(packet)

		// 8. Generate final key
		keh.deriveFinalKey(packet)

		// 9. Send challenge
		packet2 := network.EmptyPacket()
		packet2.MessageID = 0x5000
		packet2.WriteByte(network.EncodingKeyChallenge)
		packet2.WriteUInt64(packet.Session.Context.LocalSignature)

		packet.Session.Conn.Write(packet2.ToBytes())

	}
}

func (keh KeyExchangeHandler) decryptRemoteChallenge(packet PacketChannelData) {
	decryptedRemoteChallenge := make([]byte, 8)
	binary.LittleEndian.PutUint64(decryptedRemoteChallenge, packet.Session.Context.RemoteSignature)
	packet.Session.Cipher.DecryptRev(decryptedRemoteChallenge, decryptedRemoteChallenge)
	packet.Session.Context.RemoteSignature = utils2.ByteArrayToUint64(decryptedRemoteChallenge)
}

func (keh KeyExchangeHandler) calculateNewLocalChallenge(packet PacketChannelData) {
	localChallenge := security.CalculateChallenge(packet.Session.Context.CommonSecret, packet.Session.Context.LocalPublic, packet.Session.Context.RemotePublic)
	packet.Session.Cipher.EncryptRev(localChallenge, localChallenge)
	packet.Session.Context.LocalSignature = utils2.ByteArrayToUint64(localChallenge)
}

func (keh KeyExchangeHandler) deriveFinalKey(packet PacketChannelData) {
	finalKey := security.KeyTransformValue(utils2.Uint64ToByteArray(packet.Session.Context.HandshakeBlowfishKey), packet.Session.Context.CommonSecret, 3)
	packet.Session.Cipher, _ = blowfish.NewCipher(finalKey)
	packet.Session.Context.HandshakeBlowfishKey = utils2.ByteArrayToUint64(finalKey)
	packet.Session.Context.CompletedLocalSetup = true
	packet.Session.Context.CompletedRemoteSetup = true
}
