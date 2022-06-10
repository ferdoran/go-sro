package server

import (
	"github.com/ferdoran/go-sro/framework/network"
	log "github.com/sirupsen/logrus"
)

func StartSampleServer() {

	server := NewEngine(
		"127.0.0.1",
		15779,
		network.EncodingOptions{
			None:         false,
			Disabled:     false,
			Encryption:   true,
			EDC:          true,
			KeyExchange:  true,
			KeyChallenge: false,
		},
	)
	err := server.Start()

	if err != nil {
		log.Fatal(err)
	}
}
