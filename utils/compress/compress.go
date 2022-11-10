package compress

import (
	"github.com/Logiase/MiraiGo-Template/utils"
	"os"
)

var logger = utils.GetModuleLogger("valbot.compress")
var compresser Compression

type Compression interface {
	Compress(src []byte) []byte
	UnCompress(src []byte) []byte
}

var None = &NoCompression{}

var compressMap = map[string]Compression{
	"none": None,
	"zlib": &ZLibCompression{},
	"gzip": &GZipCompression{},
}

func DoCompress(src []byte) []byte {
	return compresser.Compress(src)
}

func DoUnCompress(compressSrc []byte) []byte {
	return compresser.UnCompress(compressSrc)
}

func SwitchType(t string) {
	if com, ok := compressMap[t]; ok {
		compresser = com
		logger.Infof("成功切換到 %s 壓縮模式。", t)
	} else {
		compresser = compressMap["none"]
		logger.Warnf("未知的壓縮類型: %s, 將使用無壓縮模式", t)
	}
}

func init() {
	t := os.Getenv("COMPRESS_TYPE")
	SwitchType(t)
}
