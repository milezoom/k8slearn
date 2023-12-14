package config

import (
	"strings"

	commonConfig "git.bluebird.id/mybb-ms/aphrodite/config"
)

var appConfig map[string]commonConfig.Value

var defaultConfig = map[string]interface{}{
	"app_name":      "printsvc",
	"grpc_port":     6002,
	"rest_port":     8002,
	"log_level":     "INFO",
	"log_directory": "",

	"pubsub_emulator_host_port": "",
	"pubsub_credential":         "",
	"pubsub_project_id":         "",

	"hellosvc_host": "",
	"hellosvc_port": "",
}

func LoadConfigMap() {
	appConfig = commonConfig.LoadConfig(defaultConfig)
}

func GetConfig(key string) (val commonConfig.Value) {
	if v, ok := appConfig[strings.ToLower(key)]; ok {
		val = v
	}
	return
}
