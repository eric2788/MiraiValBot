package cosplayer

import (
	"errors"
	"github.com/eric2788/common-utils/request"
)

const ovooaUrl = "http://ovooa.com/API/cosplay/api.php"

type (
	OVOOA struct {
	}

	ovooaResp struct {
		Code string `json:"code"`
		Text string `json:"text"`

		Data struct {
			Title string   `json:"title"`
			Data  []string `json:"data"`
		} `json:"data"`
	}
)

func (o *OVOOA) GetImages() (*Data, error) {
	var resp ovooaResp
	err := request.Get(ovooaUrl, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != "1" {
		return nil, errors.New(resp.Text)
	}
	return &Data{
		Title: resp.Data.Title,
		Urls:  resp.Data.Data,
	}, nil
}
