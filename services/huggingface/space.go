package huggingface

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/eric2788/MiraiValBot/utils/misc"
	"github.com/eric2788/common-utils/request"
)

const SpaceApiUrl = "https://%s.hf.space/api/predict"

type (
	SpaceApi struct {
		id   string
		data []interface{}
	}

	// Common Resp
	SpaceResp struct {
		Data []string `json:"data"`
	}

	SpaceImgResp struct {
		SpaceResp
		Durations    []float64 `json:"durations"`
		AvgDurations []float64 `json:"avg_durations"`
	}

	SpaceTxtResp struct {
		SpaceResp
		Duration        float64 `json:"duration"`
		AverageDuration float64 `json:"average_duration"`
	}
)

func NewSpaceApi(subdomain string, inputs ...interface{}) *SpaceApi {
	return &SpaceApi{
		id:   subdomain,
		data: inputs,
	}
}

func (s *SpaceApi) doRequest() (res *http.Response, err error) {
	body, err := json.Marshal(map[string]interface{}{
		"data": s.data,
	})
	if err != nil {
		return nil, err
	}

	res, err = http.Post(fmt.Sprintf(SpaceApiUrl, s.id), "application/json", bytes.NewReader(body))
	if err == nil && res.StatusCode != 200 {
		defer res.Body.Close() // only close when no err but non 200 code
		if b, berr := io.ReadAll(res.Body); berr == nil {
			err = fmt.Errorf(string(b))
		} else {
			err = fmt.Errorf(res.Status)
		}
	}
	return
}

func (s *SpaceApi) GetResultImages() ([][]byte, error) {
	res, err := s.doRequest()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var resp SpaceImgResp
	err = request.Read(res, &resp)
	if err != nil {
		return nil, err
	}
	var list [][]byte
	for _, d := range resp.Data {
		if !strings.HasPrefix(d, "data:image/") {
			continue
		}

		b64 := misc.TrimPrefixes(d, "data:image/png;base64,", "data:image/jpeg;base64,")
		b, err := base64.RawStdEncoding.DecodeString(b64)
		if err != nil {
			logger.Errorf("error while parsing base64: %v, source: %s", err, b64)
		} else if IsImageBlocked(b) {
			logger.Warnf("图像被NSFW过滤屏蔽")
		} else {
			list = append(list, b)
		}
	}
	return list, nil
}

func (s *SpaceApi) GetGeneratedTexts() ([]string, error) {
	res, err := s.doRequest()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var resp SpaceTxtResp
	err = request.Read(res, &resp)
	if err != nil {
		return nil, err
	}
	var list []string
	for _, d := range resp.Data {
		txts := strings.Split(d, "\n\n")
		list = append(list, txts...)
	}
	return list, nil
}
