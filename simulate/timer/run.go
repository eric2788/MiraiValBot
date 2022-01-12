package main

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/modules/timer"
	"github.com/eric2788/MiraiValBot/simulate"
	"time"
)

func main() {

	timer.RegisterTimer("A", time.Second*5, func(bot *bot.Bot) error {
		fmt.Println("A for 5 seconds")
		return nil
	})

	timer.RegisterTimer("B", time.Second*3, func(bot *bot.Bot) error {
		fmt.Println("B for 3 seconds")
		return nil
	})

	simulate.RunBasic()

	<-time.After(time.Second * 10)
	bot.Stop()

}
