package ai

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/corpix/uarand"
	"github.com/eric2788/MiraiValBot/services/huggingface"
	"github.com/eric2788/common-utils/request"
)

const (
	NovelAI8zywURL = "https://novelai.8zyw.cn/ajax.php?act=generate"

	NoR18 = `
		spread pussy, ass_visible_through_thighs, no_bra, naked_ribbon, naked_cape, naked_apron, naked, nude, bottomless, wardrobe_malfunction, nipples, nipple_slip, erect_nipples, areola, breast_grab, breast_hold, paizuri, bukkake, thigh_sex, buttjob, ass_grab, fellatio, bathing, vibrator, tentacles, sex, ass,oshiri,butt, dildo, pubic_hair, pee, pussy_juice, penis, cunnilingus, pussy,vulva, lactation, gangbang, uncensored, cum, fingering, futanari, extreme_content, censored, handjob, bestiality, masturbation, footjob, anal, cameltoe, bondage, enema, guro, nsfw,
	`

	WithR18    ExcludeType = `lowanderr` // totally r18
	WithNSFW   ExcludeType = `nsfw`      // with less r18 but no explict
	WithoutR18 ExcludeType = `r18`       // nealy no r18

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
		Data        []string               `json:"data,omitempty"`
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

// GetNovelAI8zywImage get image from 8zyw
// Deprecated: 免费生成已经转移到他的群内而非网站
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
		return "", fmt.Errorf("%s %s", resp.Msg, strings.Join(resp.Data, ", "))
	}
	return fmt.Sprintf("https://novelai.8zyw.cn/%s", resp.Url), nil
}

func New8zywPayload(tags string, exclude ExcludeType, badPrompt ...string) *NovelAI8zywPayload {
	var uc string

	if exclude == WithR18 {
		tags += ",uncensored," // 必須加這個否則沒有r18
		uc = huggingface.BadPrompt
	} else {
		uc = huggingface.BadPrompt + NoR18
	}
	return &NovelAI8zywPayload{
		Desc:    BestQualityTags + tags,
		UC:      uc + strings.Join(badPrompt, ", "),
		Exclude: exclude,
	}
}
