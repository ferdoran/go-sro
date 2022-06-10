package server

import (
	"encoding/binary"
	"github.com/ferdoran/go-sro/framework/security"
	"math/rand"
)

type SecurityContext struct {
	security.BaseSecurityContext
	InitialBlowfishKey   uint64
	HandshakeBlowfishKey uint64
	StartedHandshake     bool
	CompletedLocalSetup  bool
	CompletedRemoteSetup bool
	FinishedHandshake    bool
}

func NewContext() SecurityContext {
	prime := rand.Uint32() & 0x7FFFFFFF
	generator := rand.Uint32() & 0x7FFFFFFF
	private := rand.Uint32() & 0x7FFFFFFF
	key := make([]byte, 8)
	rand.Read(key)

	handshakeBlowfishKey := make([]byte, 8)
	rand.Read(handshakeBlowfishKey)

	ctx := SecurityContext{}
	ctx.Prime = prime
	ctx.Generator = generator
	ctx.Private = private
	ctx.InitialBlowfishKey = binary.LittleEndian.Uint64(key)
	ctx.HandshakeBlowfishKey = binary.LittleEndian.Uint64(handshakeBlowfishKey)
	ctx.LocalPublic = security.G_pow_X_mod_P(ctx.Generator, ctx.Private, ctx.Prime)
	return ctx
}
