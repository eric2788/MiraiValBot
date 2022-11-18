package waifu

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/eric2788/common-utils/request"
)

const lolicronApi = "https://api.lolicon.app/setu/v2?tag=%s&r18=%d&num=%d&size=original&keyword=%s"

type (
	Lolicron struct {
	}

	LolicronResp struct {
		Error string `json:"error"`
		Data  []struct {
			Pid        uint64            `json:"pid"`
			Uid        uint64            `json:"uid"`
			Title      string            `json:"title"`
			Author     string            `json:"author"`
			R18        bool              `json:"r18"`
			Width      int               `json:"width"`
			Height     int               `json:"height"`
			Tags       []string          `json:"tags"`
			Ext        string            `json:"ext"`
			AiType     int               `json:"aiType"`
			Urls       map[string]string `json:"urls"`
			UploadDate int64             `json:"uploadDate"`
		} `json:"data"`
	}
)

func (l *Lolicron) GetImages(option *SearchOptions) ([]*ImageData, error) {
	var resp LolicronResp
	r18 := 0
	if option.R18 {
		r18 = 1
	}
	err := request.Get(fmt.Sprintf(lolicronApi, strings.Join(option.Tags, ","), r18, option.Amount, url.QueryEscape(option.Keyword)), &resp)
	if err != nil {
		return nil, err
	} else if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}

	var results []*ImageData
	for _, data := range resp.Data {
		results = append(results, &ImageData{
			Pid:    data.Pid,
			Uid:    data.Uid,
			Author: data.Author,
			R18:    data.R18,
			Title:  data.Title,
			Tags:   data.Tags,
			Url:    data.Urls["original"],
		})
	}

	return results, nil
}
