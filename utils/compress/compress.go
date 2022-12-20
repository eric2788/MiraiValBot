package compress

import (
	"os"
	"strings"

	"github.com/Logiase/MiraiGo-Template/utils"
)

var logger = utils.GetModuleLogger("valbot.compress")
var compresser Compression

type Compression interface {
	Compress(src []byte) []byte
	UnCompress(src []byte) []byte
}

func DoCompress(src []byte) []byte {
	return compresser.Compress(src)
}

func DoUnCompress(compressSrc []byte) []byte {
	return compresser.UnCompress(compressSrc)
}

func SwitchType(t string) {
	switch strings.ToLower(t) {
	case "gzip":
		compresser = &GZipCompression{}
		logger.Infof("已切換為 Gzip 壓縮")
	case "zlib":
		compresser = &ZLibCompression{}
		logger.Infof("已切換為 Zlib 壓縮")
	default:
		compresser = &NoCompression{}
		logger.Infof("未知的压缩类型 %v, 已切換為無壓縮", t)
	}
}

func init() {
	t := os.Getenv("COMPRESS_TYPE")
	SwitchType(t)
}
