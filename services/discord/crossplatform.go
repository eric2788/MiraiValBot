package discord

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/internal/qq"
)

func StartChatListen() {
	client.AddHandler(messageCreate)
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

	var author string
	if m.Member.Nick == "" {
		author = m.Author.Username
	} else {
		author = fmt.Sprintf("%v (%v)", m.Member.Nick, m.Author.Username)
	}

	text := m.Content

	// 替換 <@uid> 為 @username
	for _, mention := range m.Mentions {
		text = strings.Replace(text, fmt.Sprintf("<@%v>", mention.ID), fmt.Sprintf("@%s", mention.Username), 1)
	}

	g, err := s.Guild(m.GuildID)

	if err != nil {
		logger.Errorf("嘗試獲取 discord guild 時失敗: %v", err)
	} else {
		// 替換 <@&role_id> 為 @role
		for _, mention := range m.MentionRoles {
			for _, role := range g.Roles {
				if role.ID == mention {
					text = strings.Replace(text, fmt.Sprintf("<@&%v>", mention), fmt.Sprintf("@%s", role.Name), 1)
				}
			}
		}
	}

	var imgs []string

	for _, attach := range m.Attachments {
		// 忽略所有不是圖片的附件
		if !strings.HasSuffix(attach.URL, ".png") &&
			!strings.HasSuffix(attach.URL, ".jpg") &&
			!strings.HasSuffix(attach.URL, ".jpeg") &&
			!strings.HasSuffix(attach.URL, ".gif") {
			continue
		}
		imgs = append(imgs, attach.URL)
	}

	logger.Debugf("作者: %v", author)
	logger.Debugf("內容: %v", text)
	logger.Debugf("圖片: %v", strings.Join(imgs, ",\n"))

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextLn("[來自Discord]"))
	msg.Append(qq.NewTextfLn("%v: %v", author, text))
	for _, img := range imgs {
		imgElement, err := qq.NewImageByUrl(img)
		if err != nil {
			logger.Errorf("上傳 discord 圖片 (%v) 失敗: %v", img, err)
			continue
		}
		msg.Append(imgElement)
	}

	_ = qq.SendWithRandomRiskyStrategy(msg)
}
