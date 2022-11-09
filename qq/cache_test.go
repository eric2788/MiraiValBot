package qq

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"testing"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/compress"
)

const imgUrl = "https://media.valorant-api.com/competitivetiers/564d8e28-c226-3180-6285-e48a390db8b1/3/ranktriangleupicon.png"

func TestSaveImage(t *testing.T) {
	gpImage := &message.GroupImageElement{}
	hash := md5.Sum([]byte(imgUrl))
	gpImage.Md5 = hash[:]
	gpImage.ImageId = hex.EncodeToString(hash[:]) + ".png"
	gpImage.Url = imgUrl

	groupMessage := &message.GroupMessage{}
	groupMessage.Id = 1
	groupMessage.Elements = []message.IMessageElement{gpImage}

	saveGroupImages(groupMessage)
}

func TestGetImage(t *testing.T) {
	hash := md5.Sum([]byte(imgUrl))
	fileName := hex.EncodeToString(hash[:])
	b, err := os.ReadFile(cacheDirPath + imagePath + fileName)
	if err != nil {
		t.Log(err)
	} else {
		t.Logf("file size: %d", len(b))
		b = compress.DoUnCompress(b)
		t.Logf("file size (uncompressed): %d", len(b))
	}
}
