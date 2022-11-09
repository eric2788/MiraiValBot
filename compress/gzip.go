package compress

import (
	"bytes"
	"compress/gzip"
	"io"
)

type GZipCompression struct {
}

func (G *GZipCompression) Compress(src []byte) []byte {
	var in bytes.Buffer
	w, _ := gzip.NewWriterLevel(&in, gzip.BestCompression)
	defer w.Close()
	if _, err := w.Write(src); err != nil {
		logger.Errorf("压缩失败: %v, 将返回原本的数据", err)
		return src
	}
	if err := w.Flush(); err != nil {
		logger.Errorf("压缩失败: %v, 将返回原本的数据", err)
		return src
	}
	return in.Bytes()
}

func (G *GZipCompression) UnCompress(src []byte) []byte {
	b := bytes.NewReader(src)
	var out bytes.Buffer
	r, err := gzip.NewReader(b)
	if err != nil {
		logger.Errorf("解压失败: %v, 将返回原本的数据", err)
		return src
	}
	defer r.Close()
	_, _ = io.Copy(&out, r)
	return out.Bytes()
}
