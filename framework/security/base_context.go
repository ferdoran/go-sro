package security

type BaseSecurityContext struct {
	Private         uint32
	LocalPublic     uint32
	Prime           uint32
	RemotePublic    uint32
	Generator       uint32
	LocalSignature  uint64
	RemoteSignature uint64
	CommonSecret    uint32
}
