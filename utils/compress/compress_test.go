package compress

import (
	"crypto/md5"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompressString(t *testing.T) {
	plain := []byte("this is a text, hello world\n")
	t.Logf("original: %s(%d)", plain, len(plain))
	compressed := DoCompress(plain)
	t.Logf("compressed: %s(%d)", compressed, len(compressed))
	result := DoUnCompress(compressed)
	t.Logf("uncompressed: %s(%d)", result, len(result))
	// enlarged ???
	//assert.Truef(t, len(compressed) <= len(plain), "got %d <= %d", len(compressed), len(plain))
	assert.Equal(t, string(plain), string(result))
}

func md5Str(b []byte) string {
	m := md5.Sum(b)
	return hex.EncodeToString(m[:])
}

func init() {
	SwitchType("zlib")
}
