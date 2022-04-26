package discord

import (
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/qq"
	"strconv"
	"strings"
)

func StartChatListen() {

	client.Identify.Intents = discordgo.IntentsGuildMessages
	client.AddHandler(messageCreate)

	err := client.Open()
	if err != nil {
		logger.Error("error opening discord connection: ", err)
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// 忽略非跨平台頻道
	if m.ChannelID != strconv.FormatInt(file.ApplicationYaml.Discord.CrossPlatChannel, 10) {
		return
	}

	author := m.Author.Username
	text := m.Content

	var imgs []string

	for _, attach := range m.Attachments {
		// 忽略所有不是圖片的附件
		if !strings.HasSuffix(attach.URL, ".png") && !strings.HasSuffix(attach.URL, ".jpg") {
			continue
		}
		imgs = append(imgs, attach.URL)
	}

	logger.Debugf("作者: %v", author)
	logger.Debugf("內容: %v", text)
	logger.Debugf("圖片: %v", strings.Join(imgs, ",\n"))

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("[Discord] %v: %v", author, text))
	for _, img := range imgs {
		imgElement, err := qq.NewImageByUrl(img)
		if err != nil {
			logger.Errorf("上傳 discord 圖片 (%v) 失敗: %v", img, err)
		}
		msg.Append(imgElement)
	}

	go qq.SendRiskyMessage(3, 5, func(try int) error {
		return qq.SendGroupMessage(msg)
	})

}
