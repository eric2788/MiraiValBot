package compress

type Compresser interface {
	Compress(src []byte) []byte
	UnCompress(src []byte) []byte
}

var compressMap = map[string]Compresser {
}

