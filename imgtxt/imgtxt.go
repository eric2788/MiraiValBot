package imgtxt

import (
	"bytes"
	"image"
	"io/ioutil"
	"net/http"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/hqbobo/text2pic"
)

type TextPrepend struct {
	prepend *text2pic.TextPicture
	font    *truetype.Font
}

func GetFontFromOwner() (*truetype.Font, error) {
	resp, err := http.Get("https://github.com/hqbobo/text2pic/blob/master/example/FZHTJW.TTF?raw=true")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return freetype.ParseFont(b)
}

func NewPrependMessage() (*TextPrepend, error) {

	f, err := GetFontFromOwner()
	if err != nil {
		return nil, err
	}

	return &TextPrepend{
		prepend: text2pic.NewTextPicture(text2pic.Configure{
			Width:   1920,
			BgColor: image.White,
		}),
		font: f,
	}, nil
}

func NewPrependMessageWithFont(f *truetype.Font) *TextPrepend {
	return &TextPrepend{
		prepend: text2pic.NewTextPicture(text2pic.Configure{
			Width:   1920,
			BgColor: image.White,
		}),
		font: f,
	}
}

func (prepend *TextPrepend) Append(element *message.TextElement) *TextPrepend {
	prepend.prepend.AddTextLine(element.Content, 12, prepend.font, text2pic.ColorBlack, text2pic.Padding{})
	return prepend
}

func (prepend *TextPrepend) GenerateImage() ([]byte, error) {
	var b []byte
	buffer := bytes.NewBuffer(b)
	err := prepend.prepend.Draw(buffer, text2pic.TypeJpeg)
	return buffer.Bytes(), err
}
