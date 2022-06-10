package network

const (
	EncodingNone         byte = 0
	EncodingDisabled     byte = 1
	EncodingEncryption   byte = 2
	EncodingEDC          byte = 4
	EncodingKeyExchange  byte = 8
	EncodingKeyChallenge byte = 16
)

type EncodingOptions struct {
	None         bool
	Disabled     bool
	Encryption   bool
	EDC          bool
	KeyExchange  bool
	KeyChallenge bool
}

func (e *EncodingOptions) GetEncodingOptionsByte() byte {
	opts := byte(0)

	if e.Disabled {
		opts |= EncodingDisabled
	}
	if e.Encryption {
		opts |= EncodingEncryption
	}
	if e.EDC {
		opts |= EncodingEDC
	}
	if e.KeyExchange {
		opts |= EncodingKeyExchange
	}
	if e.KeyChallenge {
		opts |= EncodingKeyChallenge
	}

	return opts
}

func NewEncodingOptions(options byte) EncodingOptions {
	opts := EncodingOptions{}
	if options == EncodingNone {
		opts.None = true
	}
	if options&EncodingEncryption != 0 {
		opts.Encryption = true
	}
	if options&EncodingEDC != 0 {
		opts.EDC = true
	}
	if options&EncodingKeyExchange != 0 {
		opts.KeyExchange = true
	}
	if options&EncodingKeyChallenge != 0 {
		opts.KeyChallenge = true
	}

	return opts
}
