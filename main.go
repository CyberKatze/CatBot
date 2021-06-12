package main

import (
	"github.com/m3dsh/catbot/bot"
	"github.com/m3dsh/catbot/config"
	"fmt"
)

func main() {
	err := config.ReadConfig()
	
	if err !=nil {
		fmt.Println(err.Error())
		return
	}
	bot.Start()
	<- make(chan struct{})
	return


}
