package huggingface

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/corpix/uarand"
)

const url = "https://api-inference.huggingface.co/models/%s"

type (
	FaceParam struct {
		Inputs  interface{}  `json:"inputs"`
		Options *FaceOptions `json:"options"`
	}

	FaceOptions struct {
		WaitForModel bool `json:"wait_for_model"`
		UseCache     bool `json:"use_cache"`
	}

	Option func(*FaceParam)
)

func doRequest(model string, param *FaceParam) (res *http.Response, err error) {

	token := os.Getenv("HUGGING_FACE_TOKEN")

	if token == "" {
		return nil, fmt.Errorf("HUGGING_FACE_TOKEN is not set")
	}

	body, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(url, model), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", uarand.GetRandom())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err = http.DefaultClient.Do(req)
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

func NewParam(options ...Option) *FaceParam {

	def := &FaceParam{
		Inputs: "",
		Options: &FaceOptions{
			WaitForModel: true,
			UseCache:     false,
		},
	}

	for _, o := range options {
		o(def)
	}

	return def
}

func WaitForModel(wait bool) Option {
	return func(fp *FaceParam) {
		fp.Options.WaitForModel = wait
	}
}

func UseCache(use bool) Option {
	return func(fp *FaceParam) {
		fp.Options.UseCache = use
	}
}

func Input(input interface{}) Option {
	return func(fp *FaceParam) {
		fp.Inputs = input
	}
}

func GetResultImage(model string, param *FaceParam) ([]byte, error) {
	res, err := doRequest(model, param)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}
