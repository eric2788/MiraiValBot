package aivoice

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	actors = map[string]int{
		"派蒙":    0,
		"凯亚":    1,
		"安柏":    2,
		"丽莎":    3,
		"琴":     4,
		"香菱":    5,
		"枫原万叶":  6,
		"迪卢克":   7,
		"温迪":    8,
		"可莉":    9,
		"早柚":    10,
		"托马":    11,
		"芭芭拉":   12,
		"优菈":    13,
		"云堇":    14,
		"钟离":    15,
		"魈":     16,
		"凝光":    17,
		"雷电将军":  18,
		"北斗":    19,
		"甘雨":    20,
		"七七":    21,
		"刻晴":    22,
		"神里绫华":  23,
		"戴因斯雷布": 24,
		"雷泽":    25,
		"神里绫人":  26,
		"罗莎莉亚":  27,
		"阿贝多":   28,
		"八重神子":  29,
		"宵宫":    30,
		"荒泷一斗":  31,
		"九条裟罗":  32,
		"夜兰":    33,
		"珊瑚宫心海": 34,
		"五郎":    35,
		"散兵":    36,
		"女士":    37,
		"达达利亚":  38,
		"莫娜":    39,
		"班尼特":   40,
		"申鹤":    41,
		"行秋":    42,
		"烟绯":    43,
		"久岐忍":   44,
		"辛焱":    45,
		"砂糖":    46,
		"胡桃":    47,
		"重云":    48,
		"菲谢尔":   49,
		"诺艾尔":   50,
		"迪奥娜":   51,
		"鹿野院平藏": 52,
	}
)

const (
	VoiceAPI = "https://genshin.azurewebsites.net/api/speak?format=wav&text=%s&id=%d"
)

var client = &http.Client{Timeout: time.Second * 300}

func GetGenshinVoice(msg, actor string) ([]byte, error) {
	if id, ok := actors[actor]; !ok {
		return nil, fmt.Errorf("未知的角色: %s", actor)
	} else {
		api := fmt.Sprintf(VoiceAPI, url.QueryEscape(msg), id)
		req, err := http.NewRequest(http.MethodGet, api, nil)
		if err != nil {
			return nil, fmt.Errorf("http_error: %v", err)
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		res, err := client.Do(req)
		if err != nil {
			if strings.Contains(err.Error(), "context deadline exceeded") {
				return nil, fmt.Errorf("http_error: %s", "API请求逾时")
			}
			return nil, fmt.Errorf("http_error: %v", err.Error())
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			return nil, fmt.Errorf("http_error: %v", res.Status)
		}
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return WavToAmr(b)
	}
}
