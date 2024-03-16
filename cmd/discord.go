package cmd

import (
	"github.com/chazari-x/training-api-bot/domain/discord"
	"github.com/chazari-x/training-api-bot/training"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "discord",
		Short: "discord",
		Long:  "discord",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := getConfig(cmd)

			log.SetReportCaller(false)

			log.Trace("discord bot starting..")
			defer log.Trace("discord bot stopped")

			if err := discord.StartDiscord(cfg.Discord, cfg.URLs, training.NewTraining(cfg.URLs)); err != nil {
				log.Fatalln(err)
			}
		},
	}
	cmd.PersistentFlags().String("config", "", "dev")
	cmd.PersistentFlags().String("token", "", "token")
	rootCmd.AddCommand(cmd)
}
