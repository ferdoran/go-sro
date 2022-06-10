package client

import (
	security2 "github.com/ferdoran/go-sro/framework/security"
)

type SecurityContext struct {
	security2.BaseSecurityContext
	HandshakeStarted     bool
	HandshakeCompleted   bool
	HandshakeBlowfishKey uint64
	HandshakeKey         uint64
	CommonSecret         uint32
	LocalSetupCompleted  bool
}
