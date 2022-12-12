package responses

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/corpix/uarand"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/response"
)

const apiURL = "https://lab.magiconch.com/api/nbnhhsh/guess"

var pattern = regexp.MustCompile(`^[a-zA-Z]+$`)

type (
	nbnhhsh struct {
	}

	nbnhhshResp struct {
		Name  string   `json:"name"`
		Trans []string `json:"trans"`
	}
)

func (n *nbnhhsh) ShouldHandle(msg *message.GroupMessage) bool {
	return pattern.MatchString(msg.ToString())
}

func (n *nbnhhsh) Handle(c *client.QQClient, msg *message.GroupMessage) error {

	text := pattern.FindString(msg.ToString())

	if text == "" {
		return fmt.Errorf("pattern extracted string is empty")
	}

	b, err := json.Marshal(map[string]string{
		"text": text,
	})

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(b))
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", uarand.GetRandom())
	req.Header.Set("Referer", "https://lab.magiconch.com/nbnhhsh/")
	req.Header.Set("Origin", "https://lab.magiconch.com")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var result []nbnhhshResp

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	if len(result) == 0 {
		return nil
	}

	if result[0].Trans == nil {
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	random := result[0].Trans[rand.Intn(len(result[0].Trans))]
	return qq.SendGroupMessageByGroup(msg.GroupCode, qq.CreateReply(msg).Append(message.NewText(random)))
}

func init() {
	response.AddHandle(&nbnhhsh{})
}
