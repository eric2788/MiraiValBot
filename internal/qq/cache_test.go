package qq

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"github.com/eric2788/MiraiValBot/utils/cache"
	"github.com/eric2788/MiraiValBot/utils/test"
	"testing"

	"github.com/eric2788/MiraiValBot/utils/compress"
	"github.com/eric2788/common-utils/request"
	"github.com/stretchr/testify/assert"

	"github.com/Mrs4s/MiraiGo/message"
)

const imgUrl = "https://media.valorant-api.com/competitivetiers/564d8e28-c226-3180-6285-e48a390db8b1/3/ranktriangleupicon.png"

func TestSaveAndGetImage(t *testing.T) {
	gpImage := &message.GroupImageElement{}
	hash := md5.Sum([]byte(imgUrl))
	gpImage.Md5 = hash[:]
	gpImage.ImageId = hex.EncodeToString(hash[:]) + ".png"
	gpImage.Url = imgUrl

	groupMessage := &message.GroupMessage{}
	groupMessage.Id = 1
	groupMessage.Elements = []message.IMessageElement{gpImage}
	saveGroupImages(groupMessage)

	fileName := hex.EncodeToString(hash[:])
	b, err := GetCacheImage(fileName)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("file size: %d", len(b))
	}

	bb, _ := request.GetBytesByUrl(imgUrl)

	assert.Equal(t, bb, b)

	imgs := GetImageList()
	t.Logf("image list: %d", len(imgs))
	assert.Equal(t, len(imgs), 1)
}

func TestSaveMessage(t *testing.T) {
	gpImage := &message.GroupImageElement{}
	hash := md5.Sum([]byte(imgUrl))
	gpImage.Md5 = hash[:]
	gpImage.ImageId = hex.EncodeToString(hash[:]) + ".png"
	gpImage.Url = imgUrl

	groupMessage := &message.GroupMessage{}
	groupMessage.Id = 1
	groupMessage.Elements = []message.IMessageElement{gpImage}

	saveGroupEssence(groupMessage)

	compressed, err := essenceCache.Get("1")

	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("file size: %d", len(compressed))
	}
	b := compress.DoUnCompress(compressed)
	t.Logf("uncompressed size: %d", len(b))
	persit := &PersistentGroupMessage{}
	buffer := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(persit)
	if err != nil {
		t.Fatal(err)
	}
	msg, err := persit.ToGroupMessage()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, msg.Id, groupMessage.Id)
	img := msg.Elements[0].(*message.GroupImageElement)
	assert.Equal(t, img.Md5, hash[:])
	assert.Equal(t, img.ImageId, hex.EncodeToString(hash[:])+".png")
}

func TestCompressUnCompressMessage(t *testing.T) {

	for name, compresser := range map[string]compress.Compression{
		"zlib": &compress.ZLibCompression{},
		"gzip": &compress.GZipCompression{},
		"none": compress.None,
	} {

		t.Logf("using compresser: %s", name)

		var persist = &PersistentGroupMessage{}

		groupMessage := &message.GroupMessage{}
		groupMessage.Id = 1
		groupMessage.GroupName = "test"
		groupMessage.GroupCode = 123456
		groupMessage.Sender = &message.Sender{
			Uin:      123456,
			Nickname: "tester",
		}

		gpImage := &message.GroupImageElement{}
		hash := md5.Sum([]byte(imgUrl))
		gpImage.Md5 = hash[:]
		gpImage.ImageId = "test.png"
		gpImage.Url = imgUrl

		groupMessage.Elements = []message.IMessageElement{gpImage}

		err := persist.Parse(groupMessage)

		if err != nil {
			t.Fatal(err)
		}

		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		err = enc.Encode(persist)

		if err != nil {
			t.Fatal(err)
		}

		content := buffer.Bytes()

		original := md5Str(content)
		compressed := compresser.Compress(content)

		t.Logf("original: %s", original)
		t.Logf("compressed: %s(%d), uncompressed: %s(%d)", md5Str(compressed), binary.Size(compressed), md5Str(content), binary.Size(content))

		re1 := compresser.UnCompress(compressed)
		re2 := compresser.UnCompress(content)

		t.Logf("re1: %s(%d), re2: %s(%d)", md5Str(re1), len(re1), md5Str(re2), len(re2))

		assert.Equal(t, re2, content)
		assert.Equal(t, re1, content)
		assert.Truef(t, binary.Size(compressed) <= binary.Size(content), "got %d <= %d", binary.Size(compressed), binary.Size(content))

	}
}

func md5Str(b []byte) string {
	m := md5.Sum(b)
	return hex.EncodeToString(m[:])
}

func init() {
	test.InitTesting()
	compress.SwitchType("zlib")
	imgCache = cache.NewCache(imagePath)
	essenceCache = cache.NewCache(essencePath)
}
