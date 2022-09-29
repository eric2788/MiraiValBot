package imgtxt

import (
	"bytes"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/qq"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/hqbobo/text2pic"
	"image"
	"io/ioutil"
	"net/http"
)

type (
	TextImage struct {
		prepend *text2pic.TextPicture
		font    *truetype.Font
	}

	Options struct {
		FontUrl string
		Width   int
	}
)

const (
	DefaultWidth = 1200
	ownerFont    = "https://github.com/hqbobo/text2pic/blob/master/example/FZHTJW.TTF?raw=true"
	DefaultFont  = "https://github.com/bingwen/befit/raw/master/static/resources/%E5%AD%97%E4%BD%93%E5%8C%85/simhei_0.ttf"
)

func GetFontByURL(url string) (*truetype.Font, error) {
	resp, err := http.Get(url)
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

func NewPrependMessage() (*TextImage, error) {
	return NewPrependMessageWithOptions(nil)
}

func NewPrependMessageWithOptions(options *Options) (*TextImage, error) {
	if options == nil {
		options = &Options{}
	}
	if options.FontUrl == "" {
		options.FontUrl = DefaultFont
	}
	if options.Width == 0 {
		options.Width = DefaultWidth
	}
	f, err := GetFontByURL(options.FontUrl)
	if err != nil {
		return nil, err
	}
	return &TextImage{
		prepend: text2pic.NewTextPicture(text2pic.Configure{
			Width:   options.Width,
			BgColor: image.White,
		}),
		font: f,
	}, nil
}

func (prepend *TextImage) Append(element *message.TextElement) *TextImage {
	prepend.prepend.AddTextLine(element.Content, 10, prepend.font, text2pic.ColorBlack, text2pic.Padding{})
	return prepend
}

func (prepend *TextImage) GenerateImage() ([]byte, error) {
	var b []byte
	buffer := bytes.NewBuffer(b)
	err := prepend.prepend.Draw(buffer, text2pic.TypePng)
	return buffer.Bytes(), err
}

func (prepend *TextImage) ToGroupImageElement() (*message.GroupImageElement, error) {
	b, err := prepend.GenerateImage()
	if err != nil {
		return nil, err
	}
	return qq.NewImageByByte(b)
}

func (prepend *TextImage) ToPrivateImageElement(uid int64) (*message.FriendImageElement, error) {
	b, err := prepend.GenerateImage()
	if err != nil {
		return nil, err
	}
	return qq.NewImagesByByteWithPrivate(uid, b)
}
