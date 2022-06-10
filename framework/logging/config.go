package logging

import (
	"github.com/ferdoran/go-sro/framework/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

const (
	UnhandledPacketsLogfile = "unhandled_packets.txt"
)

func Init() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:               true,
		EnvironmentOverrideColors: true,
		FullTimestamp:             true,
		TimestampFormat:           "2006-01-01 15:04:05.000",
	})
	logLevel, err := log.ParseLevel(viper.GetString(config.LogLevel))

	if err != nil {
		log.Warnf("failed to parse log level: [%s]. setting to info", logLevel)
		logLevel = log.InfoLevel
	}

	log.SetLevel(logLevel)
}

func UnhandledPacketLogger() *log.Logger {
	unhandledPacketsLog, err := os.Create(UnhandledPacketsLogfile)
	if err != nil {
		log.Error(err)
	}
	logger := log.New()
	logger.SetOutput(unhandledPacketsLog)
	logger.SetFormatter(&log.TextFormatter{
		ForceColors:               true,
		EnvironmentOverrideColors: true,
		FullTimestamp:             true,
		TimestampFormat:           "2006-01-01 15:04:05.000",
	})
	return logger
}
