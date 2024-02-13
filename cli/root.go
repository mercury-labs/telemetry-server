package cli

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

const (
	DefaultEnvironment = "local"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "server",
	Short: "telemetry server",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("loading common config")

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatal().Err(err).Msg("fatal error config file")
	}
	if viper.GetBool("debug") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	environment := getEnvironmentName()
	log.Info().Str("environment", environment).Msg("loading config")
	viper.SetConfigName(fmt.Sprintf("config.%s", environment))
	_ = viper.MergeInConfig()

	if err := RootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("root exec")
	}
}

func getEnvironmentName() string {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = DefaultEnvironment
	}
	return environment
}
