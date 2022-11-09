package compress

type Compression interface {
	Compress(src []byte) []byte
	UnCompress(src []byte) []byte
}

var None = &NoCompression{}

var compressMap = map[string]Compression{
	"":     None,
	"none": None,
	"zlib": &ZLibCompression{},
	"gzip": &GZipCompression{},
}

type NoCompression struct {
}

func (n *NoCompression) Compress(src []byte) []byte {
	return src
}

func (n *NoCompression) UnCompress(src []byte) []byte {
	return src
}
