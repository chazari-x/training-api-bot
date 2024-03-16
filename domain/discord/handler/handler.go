package handler

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/chazari-x/training-api-bot/domain/discord/command"
	"github.com/chazari-x/training-api-bot/domain/discord/logger"
	"github.com/chazari-x/training-api-bot/domain/discord/model"
	model2 "github.com/chazari-x/training-api-bot/model"
	"github.com/chazari-x/training-api-bot/training"
)

type Handler struct {
	urls     model2.URLs
	bot      *discordgo.Session
	cmd      *command.Command
	cmds     *model.Command
	log      *logger.Logger
	training *training.Training
}

func NewHandler(urls model2.URLs, bot *discordgo.Session, cmd *command.Command, cmds *model.Command, log *logger.Logger, t *training.Training) *Handler {
	return &Handler{urls: urls, bot: bot, cmds: cmds, cmd: cmd, log: log, training: t}
}

func (h *Handler) InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		switch i.ApplicationCommandData().Name {
		case h.cmds.Help.Name:
			h.commandHelp(s, i)
		case h.cmds.Training.User.Name:
			h.commandApiTrainingUser(s, i)
		case h.cmds.Training.Admins.Name:
			h.commandApiTrainingAdmins(s, i)
		default:
			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "Неизвестная команда.",
				}}
			if err := s.InteractionRespond(i.Interaction, response); h.log.Error(err, "Неизвестная команда", func() uintptr {
				pc, _, _, _ := runtime.Caller(0)
				return pc
			}()) {
				h.sendCommandError(s, i)
				return
			}
			h.deleteAnswer(s, i)
		}
	default:
		if len(i.MessageComponentData().Values) > 0 {
			if strings.Contains(i.MessageComponentData().Values[0], "adminspage") {
				h.commandApiTrainingAdmins(s, i)
				return
			}
		}
	}
}

func (h *Handler) commandHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var embed = &discordgo.MessageEmbed{
		Color:  0xffffff,
		Fields: []*discordgo.MessageEmbedField{},
		Author: &discordgo.MessageEmbedAuthor{
			Name:    h.bot.State.User.Username,
			IconURL: h.bot.State.User.AvatarURL(""),
		}}
	var commands = []*discordgo.ApplicationCommand{
		h.cmds.Help,
		h.cmds.Training.User,
		h.cmds.Training.Admins,
	}

	if len(i.ApplicationCommandData().Options) > 0 {
		for _, applicationCommand := range commands {
			if i.ApplicationCommandData().Options[0].StringValue() == applicationCommand.Name {
				embed.Title = fmt.Sprintf("Команда /%s", applicationCommand.Name)
				embed.Description = fmt.Sprintf("%s%s%s", applicationCommand.Description, func() string {
					if applicationCommand.DefaultMemberPermissions != nil {
						return "\n\nПрисутствуют ограничения!"
					}

					return ""
				}(), func() string {
					if len(applicationCommand.Options) > 0 {
						return "\n\n**Список опций для команды:**"
					}
					return ""
				}())
				for _, option := range applicationCommand.Options {
					embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
						Name:  option.Name,
						Value: fmt.Sprintf("```%s```", option.Description),
					})
				}
			}
		}
	} else {
		embed.Title = "Список команд"
		embed.Description = fmt.Sprintf("Здравствуй, пользователь! Я - бот дискорд сообщества [.chazari](https://chazari.ru), мое имя %s. Я готов помочь тебе с информацией о доступных командах. Вот некоторые из них:", h.bot.State.User.Username)
		for _, applicationCommand := range commands {
			field := &discordgo.MessageEmbedField{
				Name: fmt.Sprintf("/%s", applicationCommand.Name),
				Value: fmt.Sprintf("```%s%s%s```", applicationCommand.Description, func() string {
					if applicationCommand.DefaultMemberPermissions != nil {
						return "\n\nПрисутствуют ограничения!"
					}

					return ""
				}(), func() string {
					if len(applicationCommand.Options) > 0 {
						return fmt.Sprintf("\n\nПодробнее: /%s %s:%s", h.cmds.Help.Name, h.cmds.Help.Options[0].Name, applicationCommand.Name)
					}

					return ""
				}()),
			}
			embed.Fields = append(embed.Fields, field)
		}
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		}}
	if err := s.InteractionRespond(i.Interaction, response); h.log.Error(err, "Отправка ответа на команду help", func() uintptr {
		pc, _, _, _ := runtime.Caller(0)
		return pc
	}()) {
		h.sendCommandError(s, i)
	}
}

func (h *Handler) commandApiTrainingUser(s *discordgo.Session, i *discordgo.InteractionCreate) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{}}
	switch len(i.ApplicationCommandData().Options) {
	case 1:
		stats, err := h.training.GetUser(i.ApplicationCommandData().Options[0].StringValue())
		if h.log.Error(err, "Ошибка при отправке запроса", func() uintptr {
			pc, _, _, _ := runtime.Caller(0)
			return pc
		}()) {
			h.sendCommandError(s, i)
			return
		}

		embed := &discordgo.MessageEmbed{
			Title:  fmt.Sprintf("Информация аккаунта %s #%d", stats.Data.Login, stats.Data.ID),
			Color:  0xffffff,
			Fields: []*discordgo.MessageEmbedField{},
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "TRAINING API",
				IconURL: "https://forum.training-server.com/assets/logo-yhzuhaii.png",
			}}

		if stats.Data.Moder > 0 || stats.Data.Verify > 0 {
			if stats.Data.Moder > 0 {
				embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
					Name:   "Статус аккаунта:",
					Value:  "```Модератор сервера```",
					Inline: true,
				})
			}

			if stats.Data.Verify > 0 {
				color := regexp.MustCompile(`\{([a-fA-F0-9]{6})}`).FindString(stats.Data.VerifyText)
				if len(color) >= 8 {
					if colorInt, err := strconv.ParseInt(color[1:7], 16, 64); !h.log.Error(err, "", func() uintptr {
						pc, _, _, _ := runtime.Caller(0)
						return pc
					}()) {
						embed.Color = int(colorInt)
					}
				}

				embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
					Name: "Подтвержденный аккаунт",
					Value: fmt.Sprintf("```%s```", func(s string) string {
						if s == "" {
							return " "
						}
						return s
					}(stats.Data.VerifyText)),
					Inline: true,
				})
			}
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "**\t\t\t\t\t\t\t\t\tАктивность аккаунта:**",
			Inline: false,
		}, &discordgo.MessageEmbedField{
			Name:   "Дата регистрации:",
			Value:  fmt.Sprintf("```%s```", stats.Data.Regdate),
			Inline: true,
		})

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: "Последняя авторизация:",
			Value: func() string {
				if stats.Data.Online > 0 {
					return fmt.Sprintf("```В сети. ID: %d```", stats.Data.Playerid)
				}
				return fmt.Sprintf("```%s```", stats.Data.Lastlogin)
			}(),
			Inline: true,
		})

		if len(stats.Data.Warn) > 0 {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "**\t\t\t\t\t\t\t\tСписок предупреждений:**",
				Inline: false,
			})
			for _, w := range stats.Data.Warn {
				embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
					Name:   fmt.Sprintf("%s - *%s*", w.Admin, w.Bantime),
					Value:  fmt.Sprintf("```%s```", w.Reason),
					Inline: false,
				})
			}
		}

		response.Data.Flags = discordgo.MessageFlagsHasThread
		response.Data.Embeds = []*discordgo.MessageEmbed{embed}
		if err := s.InteractionRespond(i.Interaction, response); h.log.Error(err, "Отправка ответа на команду user", func() uintptr {
			pc, _, _, _ := runtime.Caller(0)
			return pc
		}()) {
			h.sendCommandError(s, i)
			return
		}
	default:
		response.Data.Flags = discordgo.MessageFlagsEphemeral
		response.Data.Content = "Неверно введены опции для команды."
		if err := s.InteractionRespond(i.Interaction, response); h.log.Error(err, "Отправка ошибки на команду user", func() uintptr {
			pc, _, _, _ := runtime.Caller(0)
			return pc
		}()) {
			h.sendCommandError(s, i)
			return
		}
		h.deleteAnswer(s, i)
	}
}

func (h *Handler) commandApiTrainingAdmins(s *discordgo.Session, i *discordgo.InteractionCreate) {
	admins, err := h.training.GetAdmins()
	if h.log.Error(err, "Ошибка при отправке запроса", func() uintptr {
		pc, _, _, _ := runtime.Caller(0)
		return pc
	}()) {
		h.sendCommandError(s, i)
		return
	}

	if i.Type == discordgo.InteractionApplicationCommand {
		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsHasThread,
				Embeds: []*discordgo.MessageEmbed{{
					Title:  fmt.Sprintf("Список администраторов (Всего %d)", len(admins)),
					Color:  0xffffff,
					Fields: []*discordgo.MessageEmbedField{},
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "TRAINING API",
						IconURL: "https://forum.training-server.com/assets/logo-yhzuhaii.png",
					}}},
				Components: []discordgo.MessageComponent{},
			}}

		for i, admin := range admins {
			response.Data.Embeds[0].Fields = append(response.Data.Embeds[0].Fields, &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("**№%d.** %s #%d", i+1, admin.Login, admin.ID),
				Inline: false,
			}, &discordgo.MessageEmbedField{
				Name:   "Авторизация:",
				Value:  fmt.Sprintf("```%s```", time.Unix(int64(admin.LastLogin), 0).Format("2006-01-02 15:04:05")),
				Inline: true,
			}, &discordgo.MessageEmbedField{
				Name:   "Предупреждений:",
				Value:  fmt.Sprintf("```%d```", admin.Warn),
				Inline: true,
			})

			if (i+1)%5 == 0 {
				break
			}
		}

		var components []discordgo.MessageComponent

		selectMenu := &discordgo.SelectMenu{
			CustomID:    "menu",
			Placeholder: "Выберите опцию",
			Options: []discordgo.SelectMenuOption{
				{
					Default: true,
					Label:   "Страница 1",
					Value:   "01adminspage",
					Emoji: discordgo.ComponentEmoji{
						Name: "__button",
						ID:   "1218469987455340564",
					},
				},
			},
		}

		for i := 6; i <= len(admins); i += 5 {
			selectMenu.Options = append(selectMenu.Options, discordgo.SelectMenuOption{
				Label: "Страница " + strconv.Itoa((i+4)/5),
				Value: func(page int) string {
					if page < 10 {
						return "0" + strconv.Itoa(page)
					}
					return strconv.Itoa(page)
				}((i+4)/5) + "adminspage",
				Emoji: discordgo.ComponentEmoji{
					Name: "__button",
					ID:   "1218469987455340564",
				},
			})
		}

		selectMenu.Options = append(selectMenu.Options, discordgo.SelectMenuOption{
			Label: "Закрыть",
			Value: "00adminspage",
			Emoji: discordgo.ComponentEmoji{
				Name: "__button",
				ID:   "1218469987455340564",
			},
		})

		components = append(components, selectMenu)

		actions := discordgo.ActionsRow{
			Components: components,
		}
		response.Data.Components = append(response.Data.Components, actions)
		if err := s.InteractionRespond(i.Interaction, response); h.log.Error(err, "Отправка первой страницы списка admins", func() uintptr {
			pc, _, _, _ := runtime.Caller(0)
			return pc
		}()) {
			h.sendCommandError(s, i)
			return
		}
		return
	}

	message, err := s.ChannelMessage(i.ChannelID, i.Message.ID)
	if h.log.Error(err, "Получение ид отправителя команды admins", func() uintptr {
		pc, _, _, _ := runtime.Caller(0)
		return pc
	}()) {
		return
	}
	if message.Interaction.User.ID == i.Member.User.ID {
		page, err := strconv.Atoi(i.MessageComponentData().Values[0][:2])
		if h.log.Error(err, "Получение номера необходимой страницы", func() uintptr {
			pc, _, _, _ := runtime.Caller(0)
			return pc
		}()) {
			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: fmt.Sprintf("Ошибка получения информации о текущей странице: %s", err),
				}}
			if err := s.InteractionRespond(i.Interaction, response); h.log.Error(err, "Отправка ошибки для команды admins", func() uintptr {
				pc, _, _, _ := runtime.Caller(0)
				return pc
			}()) {
				h.sendCommandError(s, i)
				return
			}
			return
		}

		if page == 0 {
			if err := s.ChannelMessageDelete(i.ChannelID, i.Message.ID); h.log.Error(err, "Закрыть список admins", func() uintptr {
				pc, _, _, _ := runtime.Caller(0)
				return pc
			}()) {
				response := &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   discordgo.MessageFlagsHasThread,
						Content: "Ошибка закрытия списка.",
					}}
				if err := s.InteractionRespond(i.Interaction, response); h.log.Error(err, "Отправка информации, что команда обработана", func() uintptr {
					pc, _, _, _ := runtime.Caller(0)
					return pc
				}()) {
					h.sendCommandError(s, i)
					return
				}
				h.deleteAnswer(s, i)
				return
			}
			response := &discordgo.InteractionResponse{Type: discordgo.InteractionResponseDeferredMessageUpdate}
			h.log.Error(s.InteractionRespond(i.Interaction, response), "Отправка информации, что команда обработана", func() uintptr {
				pc, _, _, _ := runtime.Caller(0)
				return pc
			}())
			return
		}

		embeds := []*discordgo.MessageEmbed{{
			Title:  fmt.Sprintf("Список администраторов (Всего %d)", len(admins)),
			Color:  0xffffff,
			Fields: []*discordgo.MessageEmbedField{},
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "TRAINING API",
				IconURL: "https://forum.training-server.com/assets/logo-yhzuhaii.png",
			}}}

		for i, admin := range admins {
			if i+1 > page*5-5 {
				embeds[0].Fields = append(embeds[0].Fields, &discordgo.MessageEmbedField{
					Name:   fmt.Sprintf("**№%d.** %s #%d", i+1, admin.Login, admin.ID),
					Inline: false,
				}, &discordgo.MessageEmbedField{
					Name:   "Авторизация:",
					Value:  fmt.Sprintf("```%s```", time.Unix(int64(admin.LastLogin), 0).Format("2006-01-02 15:04:05")),
					Inline: true,
				}, &discordgo.MessageEmbedField{
					Name:   "Предупреждений:",
					Value:  fmt.Sprintf("```%d```", admin.Warn),
					Inline: true,
				})

				if (i+1)%5 == 0 {
					break
				}
			}
		}

		actions := discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{},
		}

		selectMenu := &discordgo.SelectMenu{
			CustomID:    "menu",
			Placeholder: "Выберите опцию",
			Options:     []discordgo.SelectMenuOption{},
		}

		for i := 1; i <= len(admins); i += 5 {
			selectMenu.Options = append(selectMenu.Options, discordgo.SelectMenuOption{
				Label: "Страница " + strconv.Itoa((i+4)/5),
				Default: func(page, i int) bool {
					return page == i
				}(page, (i+4)/5),
				Value: func(page int) string {
					if page < 10 {
						return "0" + strconv.Itoa(page)
					}
					return strconv.Itoa(page)
				}((i+4)/5) + "adminspage",
				Emoji: discordgo.ComponentEmoji{
					Name: "__button",
					ID:   "1218469987455340564",
				},
			})
		}

		selectMenu.Options = append(selectMenu.Options, discordgo.SelectMenuOption{
			Label: "Закрыть",
			Value: "00adminspage",
			Emoji: discordgo.ComponentEmoji{
				Name: "__button",
				ID:   "1218469987455340564",
			},
		})

		actions.Components = append(actions.Components, selectMenu)

		messageEdit := &discordgo.MessageEdit{
			Components: append([]discordgo.MessageComponent{}, actions),
			Embeds:     embeds,
			Flags:      discordgo.MessageFlagsHasThread,
			ID:         i.Message.ID,
			Channel:    i.ChannelID,
		}
		if _, err = s.ChannelMessageEditComplex(messageEdit); h.log.Error(err, "Перехода на другую страницу", func() uintptr {
			pc, _, _, _ := runtime.Caller(0)
			return pc
		}()) {
			h.sendCommandError(s, i)
			return
		}
		response := &discordgo.InteractionResponse{Type: discordgo.InteractionResponseDeferredMessageUpdate}
		h.log.Error(s.InteractionRespond(i.Interaction, response), "Отправка информации, что команда обработана", func() uintptr {
			pc, _, _, _ := runtime.Caller(0)
			return pc
		}())
	} else {
		response := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Вы не имеете доступа к взаимодействию с этим списком.\n\nОтправьте команду **/admins**, чтобы посмотреть весь список.",
			}}
		if err := s.InteractionRespond(i.Interaction, response); h.log.Error(err, "Отправка ошибки взаимодействия с чужим admins", func() uintptr {
			pc, _, _, _ := runtime.Caller(0)
			return pc
		}()) {
			h.sendCommandError(s, i)
			return
		}
		h.deleteAnswer(s, i)
	}

}

func (h *Handler) deleteAnswer(s *discordgo.Session, i *discordgo.InteractionCreate, seconds ...int) {
	defer func() {
		_ = s.InteractionResponseDelete(i.Interaction)
	}()

	if seconds != nil {
		time.Sleep(time.Second * time.Duration(seconds[0]))
	} else {
		time.Sleep(time.Second * 10)
	}
}

func (h *Handler) deleteMessage(s *discordgo.Session, i *discordgo.Message, seconds ...int) {
	defer func() {
		_ = s.ChannelMessageDelete(i.ChannelID, i.ID)
	}()

	if seconds != nil {
		time.Sleep(time.Second * time.Duration(seconds[0]))
	} else {
		time.Sleep(time.Second * 10)
	}
}

func (h *Handler) sendCommandError(s *discordgo.Session, i *discordgo.InteractionCreate) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{{
				Title:       "Error!",
				Description: "Произошла ошибка",
				Color:       0xff0000,
			}},
		}}
	if err := s.InteractionRespond(i.Interaction, response); !h.log.Error(err, "", func() uintptr {
		pc, _, _, _ := runtime.Caller(0)
		return pc
	}()) {
		h.deleteAnswer(s, i)
	}
}
