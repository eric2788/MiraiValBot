package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/simulate"

	_ "github.com/eric2788/MiraiValBot/hooks/cmd"
)

func main() {

	simulate.RunBasic()

	fmt.Println("開始監聽指令輸入:")
	for {

		var line string

		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			line = scanner.Text()
		} else {
			err := scanner.Err()
			log.Fatal(err)
			return
		}

		if line == "exit" {
			break
		}

		res, err := command.InvokeCommand(line, true, nil)

		if err != nil {
			fmt.Printf("處理指令時出現錯誤: %v\n", err)
			continue
		}

		if res.ShowHelp {
			fmt.Printf("顯示幫助: %s\n", res.Content)
		} else {
			fmt.Printf("處理結果:\n略過: %v\n內容: %v\n", res.Ignore, res.Content)
		}
	}

	fmt.Println("退出指令輸入")
	bot.Stop()
}
