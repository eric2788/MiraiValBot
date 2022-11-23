package imgtag

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/corpix/uarand"
	"github.com/eric2788/common-utils/request"
)

const searchTagURL = "https://api.cerfai.com/search_tags"

type CerftAiResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Data []struct {
		Remarks string `json:"remarks"`
		IsNsfw  int    `json:"is_nsfw"`
		Name    string `json:"name"`
		TName   string `json:"t_name"`
	} `json:"data"`
}

func SearchTags(keyword string) (map[string]string, error) {
	b, err := json.Marshal(map[string]string{
		"keyword": keyword,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, searchTagURL, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", uarand.GetRandom())
	req.Header.Set("Origin", "https://www.cerfai.com")
	req.Header.Set("Referer", "https://www.cerfai.com/")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	} else if res.StatusCode != 200 {
		defer res.Body.Close()
		if b, berr := io.ReadAll(res.Body); berr == nil {
			return nil, errors.New(string(b))
		} else {
			return nil, errors.New(res.Status)
		}
	}

	var resp CerftAiResp
	err = request.Read(res, &resp)
	if err != nil {
		return nil, err
	} else if resp.Code != 200 {
		return nil, errors.New(resp.Msg)
	}

	tags := make(map[string]string)
	for _, tag := range resp.Data {
		tags[tag.Name] = tag.TName
	}
	return tags, nil
}
