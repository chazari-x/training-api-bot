package model

import "github.com/bwmarrin/discordgo"

type Command struct {
	Help     *discordgo.ApplicationCommand
	Training Training
}

type Training struct {
	User   *discordgo.ApplicationCommand
	Admins *discordgo.ApplicationCommand
}
