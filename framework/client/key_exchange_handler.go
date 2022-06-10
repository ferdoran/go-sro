package client

import (
	network2 "github.com/ferdoran/go-sro/framework/network"
	security2 "github.com/ferdoran/go-sro/framework/security"
	blowfish2 "github.com/ferdoran/go-sro/framework/security/blowfish"
	utils2 "github.com/ferdoran/go-sro/framework/utils"
	log "github.com/sirupsen/logrus"
)

type KeyExchangeHandler struct {
	client *Client
}

func (h *KeyExchangeHandler) Handle(packet network2.Packet) {
	if !h.client.Ctx.HandshakeStarted {
		// First 0x5000
		h.client.Ctx.HandshakeStarted = true
		opts, err := packet.ReadByte()
		if err != nil {
			log.Error(err)
		}
		h.client.EncodingOptions = network2.NewEncodingOptions(opts)

		if h.client.EncodingOptions.Encryption {
			h.client.Ctx.HandshakeBlowfishKey, err = packet.ReadUInt64()
			if err != nil {
				log.Error(err)
			}
			h.client.Cipher, err = blowfish2.NewCipher(utils2.Uint64ToByteArray(h.client.Ctx.HandshakeBlowfishKey))
			if err != nil {
				log.Error(err)
			}
		}

		if h.client.EncodingOptions.EDC {
			sequenceSeed, err := packet.ReadUInt32()
			if err != nil {
				log.Error(err)
			}
			crcSeed, err := packet.ReadUInt32()
			if err != nil {
				log.Error(err)
			}
			h.client.Sequence = network2.NewMessageSequence(sequenceSeed)
			h.client.Crc = network2.NewMessageCRC(crcSeed)
		}

		if h.client.EncodingOptions.KeyExchange {
			h.setupHandshake(packet)
		}
	} else if h.client.Ctx.HandshakeStarted && h.client.Ctx.LocalSetupCompleted {
		encodingOption, err := packet.ReadByte()
		if err != nil {
			log.Error("Could not read encoding options")
		}
		if encodingOption&network2.EncodingKeyChallenge != 0 {
			h.validateServerSignature(packet)
			h.deriveFinalKey()
			packet := network2.EmptyClientPacket()
			packet.MessageID = 0x9000
			h.client.sendPacket(packet)
			h.sendModuleIdentifcation()
		}
	} else {
		log.Info("I don't know what to do")
	}
}

func (h *KeyExchangeHandler) setupHandshake(packet network2.Packet) {
	var err error
	h.client.Ctx.HandshakeKey, err = packet.ReadUInt64()
	if err != nil {
		log.Error("Failed to read HandshakeKey")
	}
	h.client.Ctx.Generator, err = packet.ReadUInt32()
	if err != nil {
		log.Error("Failed to read generator")
	}
	h.client.Ctx.Prime, err = packet.ReadUInt32()
	if err != nil {
		log.Error("Failed to read prime")
	}
	h.client.Ctx.RemotePublic, err = packet.ReadUInt32()
	if err != nil {
		log.Error("Failed to read server public")
	}
	h.client.Ctx.LocalPublic = security2.G_pow_X_mod_P(h.client.Ctx.Generator, h.client.Ctx.Private, h.client.Ctx.Prime)
	h.client.Ctx.CommonSecret = security2.G_pow_X_mod_P(h.client.Ctx.RemotePublic, h.client.Ctx.Private, h.client.Ctx.Prime)

	// initialize new blowfish Cipher
	newBlowfishKey := security2.CalculateKey(h.client.Ctx.CommonSecret, h.client.Ctx.RemotePublic, h.client.Ctx.LocalPublic)
	h.client.Cipher, err = blowfish2.NewCipher(newBlowfishKey)
	h.client.Ctx.LocalSignature = utils2.ByteArrayToUint64(security2.CalculateChallenge(h.client.Ctx.CommonSecret, h.client.Ctx.LocalPublic, h.client.Ctx.RemotePublic))

	signatureBytes := make([]byte, 8)
	h.client.Cipher.EncryptRev(signatureBytes, utils2.Uint64ToByteArray(h.client.Ctx.LocalSignature))
	h.client.Ctx.LocalSignature = utils2.ByteArrayToUint64(signatureBytes)

	// prepare outgoing packet
	packet2 := network2.EmptyClientPacket()
	packet2.MessageID = 0x5000
	packet2.WriteUInt32(h.client.Ctx.LocalPublic)
	packet2.WriteUInt64(h.client.Ctx.LocalSignature)

	h.client.Ctx.LocalSetupCompleted = true
	h.client.sendPacket(packet2)
}

func (h *KeyExchangeHandler) validateServerSignature(packet network2.Packet) {
	var err error
	h.client.Ctx.RemoteSignature, err = packet.ReadUInt64()
	newSignature := security2.CalculateChallenge(h.client.Ctx.CommonSecret, h.client.Ctx.LocalPublic, h.client.Ctx.RemotePublic)
	h.client.Cipher.EncryptRev(newSignature, newSignature)
	h.client.Ctx.RemoteSignature = utils2.ByteArrayToUint64(newSignature)
	if err != nil {
		log.Error("Failed to read server signature")
	}
	if h.client.Ctx.RemoteSignature != h.client.Ctx.LocalSignature {
		log.Error("Invalid Server Signature")
	}
}

func (h *KeyExchangeHandler) deriveFinalKey() {
	var err error
	finalKey := security2.KeyTransformValue(utils2.Uint64ToByteArray(h.client.Ctx.HandshakeKey), h.client.Ctx.CommonSecret, 3)
	h.client.Cipher, err = blowfish2.NewCipher(finalKey)
	if err != nil {
		log.Error(err)
	}
	h.client.Ctx.HandshakeBlowfishKey = utils2.ByteArrayToUint64(finalKey)
}

func (h *KeyExchangeHandler) sendModuleIdentifcation() {
	packet := network2.EmptyClientPacket()
	packet.MessageID = 0x2001
	packet.Encrypted = true
	packet.WriteString(h.client.ModuleID)
	packet.WriteByte(0)
	h.client.sendPacket(packet)
}
