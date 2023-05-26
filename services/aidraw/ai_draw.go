package aidraw

import (
	"fmt"
	"github.com/eric2788/common-utils/stream"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
)

const (
	badPrompt    = "text font ui, error, messy drawing, blurred, lowres, low res, bad shadow, text, ui, cropped, watermark, worst quality, bad anatomy, disfigured, mutation, mutated,  liquid body, disfigured, malformed, mutated, anatomical nonsense, bad anatomy, bad proportions, uncoordinated body, malformed feet, extra feet, bad feet, poorly drawn feet, fused feet, missing feet, extra shoes, bad shoes, fused shoes, more than two shoes, poorly drawn shoes,multiple breasts, fused breasts, bad breasts, poorly drawn breasts, extra breasts, liquid breasts, missing breasts, missing breasts, more than 2 nipples, missing nipples, different nipples, fused nipples, bad nipples, poorly drawn nipples, black nipples, colorful nipples, missing clit, bad clit, fused clit, colorful clit, black clit, liquid clit, bad camel toe, colorful camel toe, bad asshole, poorly drawn asshole, fused asshole, missing asshole, bad anus, bad pussy, bad crotch, bad crotch seam, fused anus, fused pussy, fused anus, fused crotch,  black-white"
	prefixPrompt = "masterpiece, best quality, "
)

type (
	ModelType string

	Drawer func(payload Payload) (*Response, error)

	Payload struct {
		Prompt string
		Model  string
	}

	Response struct {
		ImgUrl  string
		ImgData []byte
		Source  string
	}
)

var (
	drawableSources = make(map[string]Drawer)
	logger          = logrus.WithField("service", "aidraw")
)

func Draw(payload Payload) (res *Response, err error) {
	selected := stream.From(maps.Keys(drawableSources)).Shuffle().MustFirst()
	res, err = drawableSources[selected](payload)
	if res != nil {
		// assign which Source has been used
		res.Source = selected
	} else if err != nil {
		err = fmt.Errorf("%v (使用 %v)", err, selected)
	}
	return
}
