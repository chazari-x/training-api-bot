package command

import (
	"fmt"
	"runtime"

	"github.com/bwmarrin/discordgo"
	"github.com/chazari-x/training-api-bot/domain/discord/logger"
	"github.com/chazari-x/training-api-bot/domain/discord/model"
	model2 "github.com/chazari-x/training-api-bot/model"
	log "github.com/sirupsen/logrus"
)

type Command struct {
	bot *discordgo.Session
	cmd model.Command
	log *logger.Logger
}

func NewCommandsList(cfg model2.Discord, bot *discordgo.Session, log *logger.Logger) (*Command, *model.Command) {
	var (
		helpCmd = &discordgo.ApplicationCommand{
			Version:     "1",
			Name:        fmt.Sprintf("%shelp", cfg.Prefix),
			Description: "Получить помощь по боту.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "command",
					Description: "Получить помощь по команде.",
				},
			},
		}

		trainingUserCmd = &discordgo.ApplicationCommand{
			Version:     "1",
			Name:        fmt.Sprintf("%suser", cfg.Prefix),
			Description: "Получить информацию об игроке TRAINING'а.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "nickname",
					Description: "NickName игрока.",
					Required:    true,
				},
			},
		}

		trainingAdminsCmd = &discordgo.ApplicationCommand{
			Version:     "1",
			Name:        fmt.Sprintf("%sadmins", cfg.Prefix),
			Description: "Получить список модераторов TRAINING'а.",
		}
	)

	var command = model.Command{
		Help: helpCmd,
		Training: model.Training{
			User:   trainingUserCmd,
			Admins: trainingAdminsCmd,
		},
	}

	return &Command{bot: bot, cmd: command, log: log}, &command
}

func (c *Command) RegisterCommands() error {
	registeredCommands, err := c.bot.ApplicationCommands(c.bot.State.User.ID, "")
	if c.log.Error(err, "Получение списка зарегистрированных команд") {
		return err
	}

	var theNecessaryCommands = []*discordgo.ApplicationCommand{
		c.cmd.Help,
		c.cmd.Training.User,
		c.cmd.Training.Admins,
	}
	var incorrectlyCommands []*discordgo.ApplicationCommand
	for _, rc := range registeredCommands {
		if func() bool {
			for i, nc := range theNecessaryCommands {
				if rc.Name == nc.Name && rc.Description == nc.Description && len(rc.Options) == len(nc.Options) {
					for o, option := range nc.Options {
						if option.Type != rc.Options[o].Type || option.Description != rc.Options[o].Description ||
							option.Name != rc.Options[o].Name || option.MaxLength != rc.Options[o].MaxLength ||
							len(option.Choices) != len(rc.Options[o].Choices) || option.Required != rc.Options[o].Required {
							log.Info("Найдена некорректная команда: ", rc.Name)
							return false
						} else {
							for c, choice := range option.Choices {
								if choice.Value != rc.Options[o].Choices[c].Value || choice.Name != rc.Options[o].Choices[c].Name {
									log.Info("Найдена некорректная команда: ", rc.Name)
									return false
								}
							}
						}
					}
					var newSlice []*discordgo.ApplicationCommand
					newSlice = append(newSlice, theNecessaryCommands[:i]...)
					newSlice = append(newSlice, theNecessaryCommands[i+1:]...)
					theNecessaryCommands = newSlice
					return true
				}
			}
			log.Info("Найдена некорректная команда: ", rc.Name)
			return false
		}() {
			incorrectlyCommands = append(incorrectlyCommands, rc)
		}
	}

	_, err = c.bot.ApplicationCommandBulkOverwrite(c.bot.State.User.ID, "", incorrectlyCommands)
	if c.log.Error(err, "Удаление команд, которых не должно быть", func() uintptr {
		pc, _, _, _ := runtime.Caller(0)
		return pc
	}()) {
		return err
	}

	for _, cmd := range theNecessaryCommands {
		if _, err = c.bot.ApplicationCommandCreate(c.bot.State.User.ID, "", cmd); c.log.Error(err, "Регистрация команды", func() uintptr {
			pc, _, _, _ := runtime.Caller(0)
			return pc
		}()) {
			return fmt.Errorf("%s: %s", cmd.Name, err)
		}
	}

	return nil
}
