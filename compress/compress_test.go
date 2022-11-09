package compress

import (
	"crypto/md5"
	"encoding/hex"
	"testing"

	"github.com/eric2788/common-utils/request"
	"github.com/stretchr/testify/assert"
)

const imgUrl = "https://media.valorant-api.com/competitivetiers/564d8e28-c226-3180-6285-e48a390db8b1/3/ranktriangleupicon.png"

func TestCompressAndUnCompress(t *testing.T) {
	img, err := request.GetBytesByUrl(imgUrl)
	if err != nil {
		t.Fatal(err)
	}

	original := md5Str(img)
	compressed := DoZlibCompress(img)

	t.Logf("original: %s", original)

	t.Logf("compressed: %s(%d), uncompressed: %s(%d)", md5Str(compressed), len(compressed), md5Str(img), len(img))

	re1 := DoZlibUnCompress(compressed)
	re2 := DoZlibUnCompress(img)

	t.Logf("re1: %s, re2: %s", md5Str(re1), md5Str(re2))

	assert.Equal(t, re2, img)
	assert.True(t, len(compressed) < len(img))
}

func md5Str(b []byte) string {
	m := md5.Sum(b)
	return hex.EncodeToString(m[:])
}
