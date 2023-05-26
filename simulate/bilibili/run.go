package main

import (
	"fmt"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/hooks/sites/bilibili"
	"github.com/eric2788/MiraiValBot/simulate"

	_ "github.com/eric2788/MiraiValBot/hooks/handlers"
)

var bilibiliRooms = []int64{

	22047448,
	22571958,
	6632844,
	22853788,
	22920508,
	21320551,
	1321846,
	21402309,
	8725120,
	21013446,
	22359795,
	21685677,
	23733603,
	6632844,
	22671795,
	3822389,
	14327465,
	21342742,
	23089686,
	22347054,
	4895312,
	255,
}

func main() {

	simulate.RunBasic()

	for _, room := range bilibiliRooms {

		_, err := bilibili.StartListen(room)

		if err != nil {
			fmt.Printf("啟動監聽房間 %v 時出現錯誤: %v\n", room, err)
		}
	}

	<-time.After(time.Second * 15) // 測試
	fmt.Println("正在停止...")

	for _, room := range bilibiliRooms {

		_, err := bilibili.StopListen(room) // 使用 StopListen 來刪除離線數據

		if err != nil {
			fmt.Printf("中止監聽房間 %v 時出現錯誤: %v\n", room, err)
		}

	}

	bot.Stop()
}
