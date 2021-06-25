package bot

import (
	"github.com/m3dsh/catbot/config"
	"encoding/binary"
	"io"
	"os"
	"strings"
	"time"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"regexp"
	"github.com/fatih/color"
)

var BotID string
var CatBot *discordgo.Session
var buffer = make([][]byte, 0)
var quit = make(chan int)
var count = make(chan int)

func Start() {

	err := loadSound()
	
	if err != nil {
		fmt.Println("Error loading sound: ", err)
		return
	}
	//make a Bot
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
	
	//Register ready as a callback for the ready events
	CatBot.AddHandler(ready)
	//Register messageHandler as a callback for the messageCreate events.	
	CatBot.AddHandler(messageHandler)
	//Register guildCreate as a callback for the guildCreate events.
	CatBot.AddHandler(guildCreate)

	

	err = CatBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Bot is running! id:%q\n",BotID)
	ShowGuildMembers(CatBot,"Cat Empire")

}
func ShowGuildMembers(CatBot *discordgo.Session,GuildName string){
	var guild *discordgo.Guild
	for _, g:= range CatBot.State.Guilds{
		gg,_:=CatBot.Guild(g.ID)
		fmt.Println(gg.Name)
		if( GuildName == gg.Name){
			guild = gg
			break
		}

	}	
	if (guild == nil ){
		fmt.Println("Invalid guild name")
		return
	}
	members , _:=CatBot.GuildMembers(guild.ID, "", 100)
	
	for _, elem := range members {
		msg := fmt.Sprintf("%s: %s: ", elem.User.ID, elem.User.Username )
		fmt.Println(msg)
		for _, role := range elem.Roles {
			r, _:=CatBot.State.Role("851173362532089857", role)
			fmt.Print("\"",r.Name,"\"")
		}
		fmt.Print("\n------------------------------------------------------\n")
		//Meow(CatBot,elem.User.ID, msg) 
	}
}
func Meow(s *discordgo.Session, UserID string, msg string, quit chan int) {
	c := 0
	ch, err := s.UserChannelCreate(UserID)
	if err != nil {
		// If an error occurred, we failed to create the channel.
		//
		// Some common causes are:
		// 1. We don't share a server with the user (not possible here).
		// 2. We opened enough DM channels quickly enough for Discord to
		//    label us as abusing the endpoint, blocking us from opening
		//    new ones.
		fmt.Println("error creating channel:", err)
		return 
	}
	for  {
		select{
		case <-quit  :
			count <- c
			return
		default :
			_, err = s.ChannelMessageSend(ch.ID, msg)
			c += 1
			if err != nil {
				// If an error occurred, we failed to send the message.
				//
				// It may occur either when we do not share a server with the
				// user (highly unlikely as we just received a message) or
				// the user disabled DM in their settings (more likely).
				fmt.Println("error sending DM message:", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == BotID {
		return
	}
	
	if strings.HasPrefix(m.Content, "!Meow") {

		// Find the channel that the message came from.
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			// Could not find channel.
			return
		}
// Find the guild for that channel.
		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			// Could not find guild.
			return
		}

	for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				err = playSound(s, g.ID, vs.ChannelID)
				if err != nil {
					fmt.Println("Error playing sound:", err)
				}

				return
			}
		}

	}
	if(m.GuildID == ""){
		d := color.New(color.FgHiRed, color.Bold)
		whiteback := d.Add(color.BgHiCyan)
		loc, _ := time.LoadLocation("Asia/Tehran")
		t,_ := m.Timestamp.Parse()
		whiteback.Printf("(%s)\n %s: %s ",t.In(loc) ,m.Author.Username, m.Content)
		if(m.Content == "Meow"){
			_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("***I*** :heart: ***you %s***", m.Author.Username))
			quit <- 0
			fmt.Printf("\ncount Meow : %d\n",<- count)
		}
	}
		
	b, _:= regexp.MatchString(`.*\b[Cc]at\b.*`, m.Content)
	if b{
		_, _ = s.ChannelMessageSend(m.ChannelID, "Meoow!!")
	}
}

func ready(s *discordgo.Session, event *discordgo.Ready) {

	// Set the playing status.
	s.UpdateGameStatus(0, "!Meow")
	
}


func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, _ = s.ChannelMessageSend(channel.ID, "Meow is ready! Type !Meow while in a voice channel to play a sound.")
			return
		}
	}
}

func loadSound() error {

	file, err := os.Open("./Meow.dca")
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}

// playSound plays the current buffer to the provided channel.
func playSound(s *discordgo.Session, guildID, channelID string) (err error) {

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}
