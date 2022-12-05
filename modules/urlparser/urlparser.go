package urlparser

import (
	"net/http"
	"regexp"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/corpix/uarand"
)

const Tag = "valbot.urlparser"

var (
	logger    = utils.GetModuleLogger(Tag)
	instance  = &urlParser{}
	providers = []InfoProvider{
		&bilibili{},
	}
)

type (
	InfoProvider interface {
		ParseURL(content string) Broadcaster
	}

	Broadcaster func(bot *client.QQClient, event *message.GroupMessage)

	urlParser struct {
	}
)

func (u *urlParser) HookEvent(bot *bot.Bot) {
	bot.GroupMessageEvent.Subscribe(func(client *client.QQClient, event *message.GroupMessage) {
		msg := event.ToString()
		for _, provider := range providers {
			if do := provider.ParseURL(msg); do != nil {
				do(client, event)
			}
		}
	})
}

func parsePattern(text string, pattern *regexp.Regexp) [][]string {

	if !pattern.MatchString(text) {
		return nil
	}

	results := make([][]string, 0)

	matches := pattern.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		results = append(results, match[1:])
	}

	return results
}

func getRedirectLink(shortURL string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, shortURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", uarand.GetRandom())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	return res.Request.URL.String(), nil
}
