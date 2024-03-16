package logger

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/chazari-x/training-api-bot/model"
	log "github.com/sirupsen/logrus"
)

type Logger struct {
	cfg model.Discord
	bot *discordgo.Session
}

func NewLogger(cfg model.Discord, bot *discordgo.Session) *Logger {
	return &Logger{cfg: cfg, bot: bot}
}

func (l *Logger) Embed(embed *discordgo.MessageEmbed) {
	var tag string
	if strings.Contains(strings.ToLower(embed.Title), "ошибка") {
		tag = "1180958677721686037"
	} else if strings.Contains(strings.ToLower(embed.Title), "ошибка") {
		tag = "1180958763063189637"
	}

	_, err := l.bot.ForumThreadStartComplex(l.cfg.Channel.Log, &discordgo.ThreadStart{
		Name:        embed.Title,
		AppliedTags: []string{tag},
	}, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		pc, _, _, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		_, line := f.FileLine(pc)
		str := fmt.Sprintf("%s line %d", f.Name(), line)
		log.Errorf("%s: %s", str, err.Error())
	}
}

func (l *Logger) Error(err error, str string, pc ...uintptr) bool {
	if err == nil {
		return false
	}

	if len(pc) > 0 {
		f := runtime.FuncForPC(pc[0])
		_, line := f.FileLine(pc[0])
		str = fmt.Sprintf("%s line %d", f.Name(), line)
	}

	log.Errorf("%s: %s", str, err.Error())

	_, err = l.bot.ForumThreadStartComplex(l.cfg.Channel.Log, &discordgo.ThreadStart{
		Name:        "Ошибка",
		AppliedTags: []string{"1180958677721686037"},
	}, &discordgo.MessageSend{
		Content: fmt.Sprintf("```%s:``` ```%s```", str, err.Error()),
	})
	if err != nil {
		pc, _, _, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		_, line := f.FileLine(pc)
		str := fmt.Sprintf("%s line %d", f.Name(), line)
		log.Errorf("%s: %s", str, err.Error())
	}

	return true
}

func (l *Logger) EmbedToForum(embed *discordgo.MessageEmbed, channel string) {
	_, err := l.bot.ForumThreadStartEmbed(channel, embed.Title, 0, embed)
	if err != nil {
		pc, _, _, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		_, line := f.FileLine(pc)
		str := fmt.Sprintf("%s line %d", f.Name(), line)
		log.Errorf("%s: %s", str, err.Error())
	}
}

func (l *Logger) ErrorToForum(err error, str string, channel string, pc ...uintptr) bool {
	if err == nil {
		return false
	}

	if len(pc) > 0 {
		f := runtime.FuncForPC(pc[0])
		_, line := f.FileLine(pc[0])
		str = fmt.Sprintf("%s line %d", f.Name(), line)
	}

	log.Errorf("%s: %s", str, err.Error())

	start, err := l.bot.GuildThreadsActive("891955325391998976")
	if err != nil {
		pc, _, _, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		_, line := f.FileLine(pc)
		str := fmt.Sprintf("%s line %d", f.Name(), line)
		log.Errorf("%s: %s", str, err.Error())
	}

	if start != nil {
		for _, s := range start.Threads {
			if len(s.AppliedTags) > 0 {
				log.Info(s.AppliedTags)
			}
		}
	}

	_, err = l.bot.ForumThreadStartComplex(channel, &discordgo.ThreadStart{
		Name: "Ошибка",
	}, &discordgo.MessageSend{
		Content: fmt.Sprintf("```%s:``` ```%s```", str, err.Error()),
	})
	if err != nil {
		pc, _, _, _ := runtime.Caller(0)
		f := runtime.FuncForPC(pc)
		_, line := f.FileLine(pc)
		str := fmt.Sprintf("%s line %d", f.Name(), line)
		log.Errorf("%s: %s", str, err.Error())
	}

	return true
}
