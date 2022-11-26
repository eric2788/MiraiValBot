package ai

import (
	"fmt"
	"github.com/corpix/uarand"
	"github.com/eric2788/MiraiValBot/services/huggingface"
	"github.com/eric2788/common-utils/request"
	"net/http"
	"net/url"
	"strings"
)

const (
	NovelAI8zywURL = "https://novelai.8zyw.cn/ajax.php?act=generate"

	NoR18 = `
		spread pussy, ass_visible_through_thighs, no_bra, naked_ribbon, naked_cape, naked_apron, naked, nude, bottomless, wardrobe_malfunction, nipples, nipple_slip, erect_nipples, areola, breast_grab, breast_hold, paizuri, bukkake, thigh_sex, buttjob, ass_grab, fellatio, bathing, vibrator, tentacles, sex, ass,oshiri,butt, dildo, pubic_hair, pee, pussy_juice, penis, cunnilingus, pussy,vulva, lactation, gangbang, uncensored, cum, fingering, futanari, extreme_content, censored, handjob, bestiality, masturbation, footjob, anal, cameltoe, bondage, enema, guro, nsfw,
	`

	WithR18    ExcludeType = `lowanderr`
	WithoutR18 ExcludeType = `r18`

	BestQualityTags = `best quality,masterpiece,`
)

type (
	ExcludeType string

	NovelAI8zywPayload struct {
		Desc    string
		UC      string
		Exclude ExcludeType
	}

	NovelAI8zywResp struct {
		Code        int                    `json:"code"`
		Msg         string                 `json:"msg"`
		Vip         map[string]interface{} `json:"vip"`
		Url         string                 `json:"url"`
		Seed        int64                  `json:"seed"`
		Probability float64                `json:"probability"`
		Icon        int                    `json:"icon"`
		Node        string                 `json:"node"`
	}
)

func (p *NovelAI8zywPayload) ToURLEncode() string {
	form := url.Values{
		"desc":          {p.Desc},
		"uc":            {p.UC},
		"exclude":       {string(p.Exclude)},
		"model":         {"nai-diffusion"},
		"resolution":    {"NORMALPortrait"},
		"sampler":       {"k_euler_ancestral"},
		"steps":         {"28"},
		"scale":         {"11"},
		"qualityToggle": {"on"},
		"seed":          {""},
		"type":          {"text"},
		"file":          {""},
		"strength":      {"0.7"},
		"noise":         {"0.2"},
		"apikey":        {""},
	}
	return form.Encode()
}

func GetNovelAI8zywImage(payload *NovelAI8zywPayload) (string, error) {

	body := strings.ReplaceAll(payload.ToURLEncode(), "%20", "+")
	req, err := http.NewRequest(http.MethodPost, NovelAI8zywURL, strings.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Referer", "https://novelai.8zyw.cn/")
	req.Header.Set("Origin", "https://novelai.8zyw.cn")
	req.Header.Set("User-Agent", uarand.GetRandom())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	resp := &NovelAI8zywResp{}
	err = request.Read(res, resp)
	if err != nil {
		return "", err
	} else if resp.Code != 200 {
		return "", fmt.Errorf("novelai.8zyw.cn: %s", resp.Msg)
	}
	return fmt.Sprintf("https://novelai.8zyw.cn/%s", resp.Url), nil
}

func New8zywPayload(tags string, r18 bool, badPrompt ...string) *NovelAI8zywPayload {
	var exclude ExcludeType
	var uc string

	if r18 {
		tags += ",uncensored," // 必須加這個否則沒有r18
		exclude = WithR18
		uc = huggingface.BadPrompt
	} else {
		exclude = WithoutR18
		uc = huggingface.BadPrompt + NoR18
	}
	return &NovelAI8zywPayload{
		Desc:    BestQualityTags + tags,
		UC:      uc + strings.Join(badPrompt, ", "),
		Exclude: exclude,
	}
}
