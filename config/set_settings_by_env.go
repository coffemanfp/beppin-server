package config

import (
	"github.com/spf13/viper"
)

// SetSettingsByEnv sets the env settings by environment vars.
func SetSettingsByEnv() {
	viper.SetEnvPrefix("beppin")
	viper.AutomaticEnv()

	viper.BindEnv(
		"port",
		"host",
		"assets",
		"logsFile",
		"secretKey",
		"maxElementsPerPagination",
		"maxImageSize",

		// Database
		"db_port",
		"db_name",
		"db_user",
		"db_password",
		"db_host",
		"db_sslMode",
		"db_url",
	)
}
