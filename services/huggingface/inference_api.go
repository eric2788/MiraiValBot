package huggingface

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/corpix/uarand"
	"github.com/eric2788/common-utils/request"
	"io"
	"net/http"
	"os"
	"strings"
)

const url = "https://api-inference.huggingface.co/models/%s"

func NewInferenceApi(model string, options ...Option) *InferenceApi {

	opt := &FaceParam{
		Inputs: "",
		Options: &FaceOptions{
			WaitForModel: true,
			UseCache:     false,
		},
	}

	for _, o := range options {
		o(opt)
	}

	return &InferenceApi{
		model: model,
		param: opt,
	}
}

func (in *InferenceApi) ChangeParams(options ...Option) {
	for _, change := range options {
		change(in.param)
	}
}

func (in *InferenceApi) doRequest() (res *http.Response, err error) {

	token := os.Getenv("HUGGING_FACE_TOKEN")

	if token == "" {
		return nil, fmt.Errorf("HUGGING_FACE_TOKEN is not set")
	}

	body, err := json.Marshal(in.param)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(url, in.model), bytes.NewReader(body))
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

func (in *InferenceApi) GetResultImage() (img []byte, err error) {
	res, err := in.doRequest()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	img, err = io.ReadAll(res.Body)
	// nsfw filtered
	if err == nil && IsImageBlocked(img) {
		err = fmt.Errorf("图像被NSFW过滤屏蔽")
	}
	return
}

func (in *InferenceApi) GetGeneratedText() (string, error) {
	res, err := in.doRequest()
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	var resp []map[string]string
	err = request.Read(res, &resp)
	if err != nil {
		return "", err
	} else if len(resp) < 1 {
		return "", fmt.Errorf("no result: %v", resp)
	}
	return resp[0]["generated_text"], nil
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

func InputWithoutBracket(input string) Option {
	input = strings.ReplaceAll(input, "{", "")
	input = strings.ReplaceAll(input, "}", "")
	input = strings.ReplaceAll(input, "[", "")
	input = strings.ReplaceAll(input, "]", "")
	return func(fp *FaceParam) {
		fp.Inputs = input
	}
}

func Input(input interface{}) Option {
	return func(fp *FaceParam) {
		fp.Inputs = input
	}
}
