package huggingface

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/eric2788/MiraiValBot/utils/misc"
)

func NewSpaceApi(subdomain string, inputs ...interface{}) *SpaceApi {
	return &SpaceApi{
		Id:       subdomain,
		Data:     inputs,
		endpoint: "run/predict",
		Hash:     generateSessionHash(),
		handler:  &httpRequestHandler{}, // default http
	}
}

func (s *SpaceApi) EndPoint(endpoint string) *SpaceApi {
	s.endpoint = endpoint
	return s
}

func (s *SpaceApi) UseWebsocketHandler() *SpaceApi {
	s.handler = &websocketHandler{}
	return s
}

func (s *SpaceApi) GetResultImages() ([][]byte, error) {
	resp, err := s.handler.Handle(s)
	if err != nil {
		return nil, err
	}
	var list [][]byte
	for _, d := range resp.Data {

		// should be error message
		if strings.HasPrefix(d, "<h4>Error</h4>") {
			return nil, fmt.Errorf(d)
		}

		if !strings.HasPrefix(d, "data:image/") {
			continue
		}
		b64 := misc.TrimPrefixes(d, "data:image/png;base64,", "data:image/jpeg;base64,")
		b, err := base64.StdEncoding.DecodeString(b64)
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
	resp, err := s.handler.Handle(s)
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
