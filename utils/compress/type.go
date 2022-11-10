package compress


type NoCompression struct {
}

func (n *NoCompression) Compress(src []byte) []byte {
	return src
}

func (n *NoCompression) UnCompress(src []byte) []byte {
	return src
}
