package core

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, "!online") {
		var onlinePlayers string
		var afkPlayers string

		CloseConnection()
		OpenConnection(AppContext.Config.Hostname, AppContext.Config.Port)

		response, err := SendAndRecvPacket(CmdOnlinePlayers, nil)
		if err != nil {
			fmt.Printf("Error sending packet: %s\n", err.Error())
			return
		}

		for _, player := range response.(OnlinePlayersPayload).Players {
			if player.Afk {
				afkPlayers += "- " + player.Name
			} else {
				onlinePlayers += "- " + player.Name
			}
		}

		embed := &discordgo.MessageEmbed{
			Title:  "Online player count",
			Color:  1752220, // #1ABC9C
			Image:  nil,
			Author: &discordgo.MessageEmbedAuthor{},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Online players",
					Value:  onlinePlayers,
					Inline: false,
				},
				{
					Name:   "AFK players",
					Value:  afkPlayers,
					Inline: false,
				},
			},
		}

		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}

	if strings.Contains(m.Content, "!info") {
		var tps string
		var ramUsage string
		var upTime string

		var overworldChunks string
		var netherChunks string
		var endChunks string
		var loadedChunks string

		CloseConnection()
		OpenConnection(AppContext.Config.Hostname, AppContext.Config.Port)

		response, err := SendAndRecvPacket(CmdServerInfo, nil)
		if err != nil {
			fmt.Printf("Error sending packet: %s\n", err.Error())
			return
		}

		tps = strconv.FormatFloat(response.(ServerInfoPayload).ServerInfo.Tps, 'f', 2, 64)
		ramUsage = strconv.FormatFloat(float64(response.(ServerInfoPayload).ServerInfo.RamUsage), 'f', 1, 64)
		upTime = response.(ServerInfoPayload).ServerInfo.Uptime

		overworldChunks = strconv.FormatInt(response.(ServerInfoPayload).ServerInfo.ChunksOverworld, 10)
		netherChunks = strconv.FormatInt(response.(ServerInfoPayload).ServerInfo.ChunksNether, 10)
		endChunks = strconv.FormatInt(response.(ServerInfoPayload).ServerInfo.ChunksEnd, 10)

		loadedChunks += "Overworld: " + overworldChunks
		loadedChunks += "\nNether: " + netherChunks
		loadedChunks += "\nEnd: " + endChunks

		embed := &discordgo.MessageEmbed{
			Title:  "Server Information",
			Color:  1752220, // #1ABC9C
			Image:  nil,
			Author: &discordgo.MessageEmbedAuthor{},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "TPS",
					Value:  tps,
					Inline: false,
				},
				{
					Name:   "Ram usage",
					Value:  ramUsage,
					Inline: false,
				},
				{
					Name:   "Uptime",
					Value:  upTime,
					Inline: false,
				},
				{
					Name:   "Loaded Chunks",
					Value:  loadedChunks,
					Inline: false,
				},
			},
		}

		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}
}
