package discord

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/chazari-x/training-api-bot/domain/discord/command"
	"github.com/chazari-x/training-api-bot/domain/discord/handler"
	"github.com/chazari-x/training-api-bot/domain/discord/logger"
	"github.com/chazari-x/training-api-bot/model"
	"github.com/chazari-x/training-api-bot/training"
	log "github.com/sirupsen/logrus"
)

func StartDiscord(cfg model.Discord, urls model.URLs, t *training.Training) (err error) {
	bot, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return
	}

	l := logger.NewLogger(cfg, bot)

	c, cmds := command.NewCommandsList(cfg, bot, l)

	bot.Identify.Intents = discordgo.IntentsAll

	if err = bot.Open(); err != nil {
		return fmt.Errorf("bot open err: %s", err)
	}
	defer func() {
		_ = bot.Close()
	}()

	bot.AddHandler(handler.NewHandler(urls, bot, c, cmds, l, t).InteractionCreate)

	if err = c.RegisterCommands(); l.Error(err, "register commands", func() uintptr {
		pc, _, _, _ := runtime.Caller(0)
		return pc
	}()) {
		return fmt.Errorf("register commands err: %s", err)
	}

	log.Info("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	return
}
