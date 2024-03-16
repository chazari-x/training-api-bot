package cmd

import (
	"fmt"
	"os"

	"github.com/chazari-x/training-api-bot/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configFile = "etc/"

type Config struct {
	Log     model.Log     `yaml:"log"`
	Discord model.Discord `yaml:"bot"`
	URLs    model.URLs    `yaml:"urls"`
}

func getConfig(cmd *cobra.Command) *Config {
	var cfg Config

	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:             true,
		TimestampFormat:           "2006-01-02 15:04:05",
		ForceColors:               true,
		PadLevelText:              true,
		EnvironmentOverrideColors: true,
	})

	file, err := cmd.Flags().GetString("config")
	if err != nil {
		log.Fatalf("get flag err: %s", err)
	} else if file != "" {
		file += "."
	}

	configFile += fmt.Sprintf("config.%syaml", file)

	f, err := os.Open(configFile)
	if err != nil {
		log.Fatalf("open config file \"%s\": %s", configFile, err)
	}

	if err = yaml.NewDecoder(f).Decode(&cfg); err != nil {
		log.Fatalf("decode config file: %s", err)
	}

	level, err := log.ParseLevel(cfg.Log.Level)
	if err != nil {
		log.Fatalf("parse level err: %s", err)
	}

	if cfg.Log.Level == "" {
		cfg.Log.Level = "trace"
	}
	log.SetLevel(level)

	token, err := cmd.Flags().GetString("token")
	if err != nil || token == "" {
		log.Fatalf("get token err: %s", err)
	}

	cfg.Discord.Token = token

	return &cfg
}
