package compress

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
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
	compressed := DoCompress(img)

	t.Logf("original: %s", original)
	t.Logf("compressed: %s(%d), uncompressed: %s(%d)", md5Str(compressed), binary.Size(compressed), md5Str(img), binary.Size(img))

	re1 := DoUnCompress(compressed)
	re2 := DoUnCompress(img)

	t.Logf("re1: %s(%d), re2: %s(%d)", md5Str(re1), len(re1), md5Str(re2), len(re2))

	assert.Equal(t, re2, img)
	assert.Equal(t, re1, img)
	//assert.True(t, binary.Size(compressed) < binary.Size(img))
}

func TestCompressFile(t *testing.T) {

	err := os.MkdirAll("data/", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	img, err := request.GetBytesByUrl(imgUrl)
	if err != nil {
		t.Fatal(err)
	}

	name := md5Str(img)

	compressed_name, uncompressed_name := fmt.Sprintf("data/%s_compressed", name), fmt.Sprintf("data/%s_uncompressed", name)

	compressed := DoCompress(img)

	err = os.WriteFile(uncompressed_name, img, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(compressed_name, compressed, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	unb, err := os.ReadFile(uncompressed_name)
	if err != nil {
		t.Fatal(err)
	}
	cb, err := os.ReadFile(compressed_name)
	if err != nil {
		t.Fatal(err)
	}

	uncb := DoUnCompress(cb)

	assert.Equal(t, uncb, unb)
}

func TestCompressString(t *testing.T) {
	plain := "this is a text, hello world"
	t.Logf("original: %s", plain)
	compressed := DoCompress([]byte(plain))
	t.Logf("compressed: %s", string(compressed))
	result := DoUnCompress(compressed)
	t.Logf("uncompressed: %s", string(result))

	assert.Equal(t, plain, string(result))
}

func md5Str(b []byte) string {
	m := md5.Sum(b)
	return hex.EncodeToString(m[:])
}
