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

func NewPrependMessage() (*TextImage, error) {

	f, err := GetFontFromOwner()
	if err != nil {
		return nil, err
	}

	return &TextImage{
		prepend: text2pic.NewTextPicture(text2pic.Configure{
			Width:   1920,
			BgColor: image.White,
		}),
		font: f,
	}, nil
}

func NewPrependMessageWithFont(f *truetype.Font) *TextImage {
	return &TextImage{
		prepend: text2pic.NewTextPicture(text2pic.Configure{
			Width:   1920,
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
