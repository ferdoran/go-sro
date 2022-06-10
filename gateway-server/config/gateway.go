package config

import (
	"github.com/ferdoran/go-sro/framework/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	AgentHost     = "agent.host"
	AgentPublicIp = "agent.public_ip"
	AgentPort     = "agent.port"
	AgentSecret   = "agent.secret"
	AgentModuleId = "agent.module_id"

	GatewayHost     = "gateway.ip"
	GatewayPort     = "gateway.Port"
	GatewayModuleId = "gateway.module_id"
	GatewaySecret   = "gateway.secret"
)

func Initialize() {
	config.Initialize()

	setDefaultValues()
	bindEnvAliases()

	logrus.Info("gateway config initialized")
}

func bindEnvAliases() {
	viper.BindEnv(AgentHost, "AGENT_HOST")
	viper.BindEnv(AgentPublicIp, "AGENT_PUBLIC_IP")
	viper.BindEnv(AgentPort, "AGENT_PORT")
	viper.BindEnv(AgentSecret, "AGENT_SECRET")
	viper.BindEnv(AgentModuleId, "AGENT_MODULE_ID")

	viper.BindEnv(GatewayHost, "GATEWAY_HOST")
	viper.BindEnv(GatewayPort, "GATEWAY_PORT")
	viper.BindEnv(GatewaySecret, "GATEWAY_SECRET")
	viper.BindEnv(GatewayModuleId, "GATEWAY_MODULE_ID")
}

func setDefaultValues() {
	viper.SetDefault(AgentHost, "127.0.0.1")
	viper.SetDefault(AgentPublicIp, "127.0.0.1")
	viper.SetDefault(AgentPort, 15882)
	viper.SetDefault(AgentSecret, "agent-server")
	viper.SetDefault(AgentModuleId, "AgentServer")

	viper.SetDefault(GatewayHost, "127.0.0.1")
	viper.SetDefault(GatewayPort, 15779)
	viper.SetDefault(GatewaySecret, "gateway-server")
	viper.SetDefault(GatewayModuleId, "GatewayServer")
}
