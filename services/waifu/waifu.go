package waifu

import (
	"fmt"
	"golang.org/x/exp/maps"
	"math/rand"
	"time"
)

type (
	ImageApi interface {
		GetImages(option *SearchOptions) ([]ImageData, error)
	}

	ImageData struct {
		Pid    uint64
		Uid    uint64
		Title  string
		R18    bool
		Author string
		Url    string
		Tags   []string
		Image  []byte // 與 Url 二選一
	}

	SearchOptions struct {
		Keyword string
		Tags    []string
		Amount  int
		R18     bool
	}

	Searcher func(option *SearchOptions)
)

var providers = map[string]ImageApi{
	"lolicron": &Lolicron{},
	"pixivmoe": &PixelMoe{},
}

func NewOptions(searcher ...Searcher) *SearchOptions {
	defaultOpt := &SearchOptions{
		Amount: 5,
		R18:    false,
	}
	for _, s := range searcher {
		s(defaultOpt)
	}
	return defaultOpt
}

func WithKeyword(keyword string) Searcher {
	return func(option *SearchOptions) {
		option.Keyword = keyword
	}
}

func WithTags(tags ...string) Searcher {
	return func(option *SearchOptions) {
		option.Tags = tags
	}
}

func WithAmount(amount int) Searcher {
	return func(option *SearchOptions) {
		option.Amount = amount
	}
}

func WithR18(r18 bool) Searcher {
	return func(option *SearchOptions) {
		option.R18 = r18
	}
}

func GetRandomImages(option *SearchOptions) ([]ImageData, error) {
	rand.Seed(time.Now().UnixNano())
	keys := maps.Keys(providers)
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})
	for _, key := range keys {
		provider := providers[key]
		images, err := provider.GetImages(option)
		if err != nil {
			logger.Errorf("使用 %s 獲取圖片失敗: %v, 將使用下一個API", key, err)
			continue
		}
		return images, nil
	}
	return nil, fmt.Errorf("所有API都無法獲取圖片")
}
