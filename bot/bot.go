package bot

import (
	"github.com/m3dsh/catbot/config"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"regexp"
)

var BotID string
var CatBot *discordgo.Session

func Start() {
	CatBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := CatBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID

	CatBot.AddHandler(messageHandler)

	err = CatBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Bot is running! id:%q",BotID)
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == BotID {
		return
	}

	b, _:= regexp.MatchString(`.*\b[Cc]at\b.*`, m.Content)
	if b{
		_, _ = s.ChannelMessageSend(m.ChannelID, "Meoow!!")
	}
}
