package imgtxt

import (
	"bytes"
	"github.com/eric2788/MiraiValBot/qq"
	"image"
	"io/ioutil"
	"net/http"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/hqbobo/text2pic"
)

type TextImage struct {
	prepend *text2pic.TextPicture
	font    *truetype.Font
}

const (
	Width       = 1000
	ownerFont   = "https://github.com/hqbobo/text2pic/blob/master/example/FZHTJW.TTF?raw=true"
	DefaultFont = "https://github.com/bingwen/befit/raw/master/static/resources/%E5%AD%97%E4%BD%93%E5%8C%85/simhei_0.ttf"
)

func GetDefaultFont() (*truetype.Font, error) {
	resp, err := http.Get(DefaultFont)
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

	f, err := GetDefaultFont()
	if err != nil {
		return nil, err
	}

	return &TextImage{
		prepend: text2pic.NewTextPicture(text2pic.Configure{
			Width:   Width,
			BgColor: image.White,
		}),
		font: f,
	}, nil
}

func NewPrependMessageWithFont(f *truetype.Font) *TextImage {
	return &TextImage{
		prepend: text2pic.NewTextPicture(text2pic.Configure{
			Width:   Width,
			BgColor: image.White,
		}),
		font: f,
	}
}

func (prepend *TextImage) Append(element *message.TextElement) *TextImage {
	prepend.prepend.AddTextLine(element.Content, 12, prepend.font, text2pic.ColorBlack, text2pic.Padding{})
	return prepend
}

func (prepend *TextImage) GenerateImage() ([]byte, error) {
	var b []byte
	buffer := bytes.NewBuffer(b)
	err := prepend.prepend.Draw(buffer, text2pic.TypeJpeg)
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
