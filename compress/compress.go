package compress

import (
	"bytes"
	"compress/zlib"
	"io"

	"github.com/Logiase/MiraiGo-Template/utils"
)

var logger = utils.GetModuleLogger("valbot.compress")

func DoCompress(src []byte) []byte {
	var in bytes.Buffer
	w, _ := zlib.NewWriterLevel(&in, zlib.BestCompression)
	defer w.Close()
	w.Write(src)
	if err := w.Flush(); err != nil {
		logger.Errorf("压缩失败: %v, 将返回原本的数据", err)
		return src
	}
	return in.Bytes()
}

func DoUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, err := zlib.NewReader(b)
	if err != nil {
		logger.Errorf("解压失败: %v, 将返回原本的数据", err)
		return compressSrc
	}
	defer r.Close()
	io.Copy(&out, r)
	return out.Bytes()
}
