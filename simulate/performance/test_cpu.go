package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/simulate"

	// 最吃CPU的是 broadcaster module
	_ "github.com/eric2788/MiraiValBot/modules/broadcaster"

	// 所有廣播訂閱平台
	_ "github.com/eric2788/MiraiValBot/hooks/sites/bilibili"
	_ "github.com/eric2788/MiraiValBot/hooks/sites/twitter"
	_ "github.com/eric2788/MiraiValBot/hooks/sites/youtube"
)

func main() {

	simulate.EnableDebug()
	simulate.RunBasic()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	<-time.After(time.Second * 30)
	bot.Stop()
}
