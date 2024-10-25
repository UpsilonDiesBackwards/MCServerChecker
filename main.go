package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"mcserverchecker/core"
	"os"
	"os/signal"
)

func main() {
	core.AppContext = &core.Context{}

	err := core.LoadConfig("config.json")
	if err != nil {
		fmt.Println("Failed to load config.json: ", err.Error())
		return
	}

	core.AppContext.Discord, err = discordgo.New("Bot " + core.AppContext.Config.Token)
	if err != nil {
		fmt.Println("Failed to create Discord session", err)
		return
	}

	core.AppContext.Discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println("Successfully connected to Discord!")
	})

	err = core.AppContext.Discord.Open()
	if err != nil {
		fmt.Printf("Error! Failed to connect to Discord: %s\n", err.Error())
		return
	}

	err = core.AppContext.Discord.UpdateGameStatus(0, core.AppContext.Config.Status)
	if err != nil {
		fmt.Printf("Error! Failed to update Discord state: %s\n", err.Error())
	}

	core.AppContext.Discord.AddHandler(core.MessageCreateHandler)

	defer func(Discord *discordgo.Session) {
		_ = Discord.Close()
	}(core.AppContext.Discord)

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("\nStopping bot...")
}
