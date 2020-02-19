package env

import "github.com/spf13/viper"

func ReadOS() Config {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")

	viper.SetDefault("PRETTY_LOG_OUTPUT", true)
	viper.SetDefault("LOG_LEVEL", "DEBUG")
	//viper.SetDefault("SRC_FILE", "")
	//viper.SetDefault("OUTPUT_MODE", "")

	return Config{
		PrettyLogOutput: viper.GetBool("PRETTY_LOG_OUTPUT"),
		LogLevel:        viper.GetString("LOG_LEVEL"),
		SourceFile:      viper.GetString("SRC_FILE"),
		OutputMode:      viper.GetString("OUTPUT_MODE"),
	}
}
