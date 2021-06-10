package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

const token string = "ODUyMjE5NjIyMTgxOTYxNzY2.YMDpXw.1wapy2cEGr40CKO0P28glKxQZYA"

var BotID string

func main() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	u, err := dg.User("@me")
	if err != nil {
		fmt.Println(err.Error())
	}

	BotID := u.ID
	err = dg.Open()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Bot is runnig", BotID)
	<- make(chan struct{})

}
