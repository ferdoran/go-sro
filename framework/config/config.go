package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	DatabaseAccountDriver   = "db.account.driver"
	DatabaseAccountHost     = "db.account.host"
	DatabaseAccountPort     = "db.account.port"
	DatabaseAccountUser     = "db.account.user"
	DatabaseAccountPassword = "db.account.pw"
	DatabaseAccountDatabase = "db.account.database"
	DatabaseAccountParams   = "db.account.params"

	DatabaseShardDriver   = "db.shard.driver"
	DatabaseShardHost     = "db.shard.host"
	DatabaseShardPort     = "db.shard.port"
	DatabaseShardUser     = "db.shard.user"
	DatabaseShardPassword = "db.shard.pw"
	DatabaseShardDatabase = "db.shard.database"
	DatabaseShardParams   = "db.shard.params"

	LogLevel                = "log.level"
	LogFileUnhandledPackets = "log.unhandled_packets.file"
)

func Initialize() {
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/agent-server")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	setDefaultValues()
	bindEnvAliases()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Info("config file not found - using env config or defaults")
		} else {
			logrus.Error(err)
		}
	}

	logrus.Info("config initialized")
}

func bindEnvAliases() {

	viper.BindEnv(LogLevel, "LOG_LEVEL")
	viper.BindEnv(LogFileUnhandledPackets, "LOG_UNHANDLED_PACKETS_FILE")

	viper.BindEnv("db.account.driver", "DATABASE_ACCOUNT_DRIVER")
	viper.BindEnv("db.account.host", "DB_ACCOUNT_HOST")
	viper.BindEnv("db.account.port", "DB_ACCOUNT_PORT")
	viper.BindEnv("db.account.user", "DB_ACCOUNT_USER")
	viper.BindEnv("db.account.pw", "DB_ACCOUNT_PW")
	viper.BindEnv("db.account.database", "DB_ACCOUNT_DATABASE")
	viper.BindEnv("db.account.params", "DB_ACCOUNT_PARAMS")

	viper.BindEnv("db.shard.driver", "DATABASE_SHARD_DRIVER")
	viper.BindEnv("db.shard.host", "DB_SHARD_HOST")
	viper.BindEnv("db.shard.port", "DB_SHARD_PORT")
	viper.BindEnv("db.shard.user", "DB_SHARD_USER")
	viper.BindEnv("db.shard.pw", "DB_SHARD_PW")
	viper.BindEnv("db.shard.database", "DB_SHARD_DATABASE")
	viper.BindEnv("db.shard.params", "DB_SHARD_PARAMS")
}

func setDefaultValues() {
	viper.SetDefault(LogLevel, "info")
	viper.SetDefault(LogFileUnhandledPackets, "unhandled_packets.log")

	viper.SetDefault(DatabaseAccountDriver, "mysql")
	viper.SetDefault(DatabaseAccountHost, "127.0.0.1")
	viper.SetDefault(DatabaseAccountPort, 3306)
	viper.SetDefault(DatabaseAccountUser, "sro")
	viper.SetDefault(DatabaseAccountPassword, "1234")
	viper.SetDefault(DatabaseAccountDatabase, "SRO_ACCOUNT")
	viper.SetDefault(DatabaseAccountParams, "?parseTime=true")

	viper.SetDefault(DatabaseShardDriver, "mysql")
	viper.SetDefault(DatabaseShardHost, "127.0.0.1")
	viper.SetDefault(DatabaseShardPort, 3306)
	viper.SetDefault(DatabaseShardUser, "sro")
	viper.SetDefault(DatabaseShardPassword, "1234")
	viper.SetDefault(DatabaseShardDatabase, "SRO_SHARD")
	viper.SetDefault(DatabaseShardParams, "?parseTime=true")
}

// "user:password@tcp(127.0.0.1:3306)/hello?params"

func ConnStringAccount() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s",
		viper.GetString(DatabaseAccountUser),
		viper.GetString(DatabaseAccountPassword),
		viper.GetString(DatabaseAccountHost),
		viper.GetString(DatabaseAccountPort),
		viper.GetString(DatabaseAccountDatabase),
		viper.GetString(DatabaseAccountParams))
}

func ConnStringShard() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s",
		viper.GetString(DatabaseShardUser),
		viper.GetString(DatabaseShardPassword),
		viper.GetString(DatabaseShardHost),
		viper.GetString(DatabaseShardPort),
		viper.GetString(DatabaseShardDatabase),
		viper.GetString(DatabaseShardParams))
}
