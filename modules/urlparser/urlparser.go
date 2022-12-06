package urlparser

import (
	"net/http"
	"regexp"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/corpix/uarand"
	"github.com/eric2788/MiraiValBot/internal/eventhook"
)

const Tag = "valbot.urlparser"

var (
	linkPattern = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
	logger      = utils.GetModuleLogger(Tag)
	instance    = &urlParser{}
	providers   = []InfoProvider{
		&bilibili{},
		&common{}, // 最底部
	}
)

type (
	InfoProvider interface {
		ParseURL(url string) Broadcaster
	}

	Broadcaster func(bot *client.QQClient, event *message.GroupMessage) error

	urlParser struct {
	}
)

func (u *urlParser) HookEvent(bot *bot.Bot) {
	bot.GroupMessageEvent.Subscribe(func(client *client.QQClient, event *message.GroupMessage) {

		if event.Sender.Uin == client.Uin {
			return
		}

		msg := event.ToString()
		if !linkPattern.MatchString(msg) {
			return
		}

		links := linkPattern.FindAllString(msg, -1)

		for _, link := range links {
			for _, provider := range providers {
				if do := provider.ParseURL(link); do == nil {
					continue
				} else if err := do(client, event); err != nil {
					logger.Error(err)
				} else {
					break
				}
			}
		}
	})
}

func parsePattern(url string, pattern *regexp.Regexp) []string {

	if !pattern.MatchString(url) {
		return nil
	}

	match := pattern.FindStringSubmatch(url)

	if len(match) < 2 {
		return nil
	}

	return match[1:]
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

func init() {
	eventhook.RegisterAsModule(instance, "URL解析", Tag, logger)
}
