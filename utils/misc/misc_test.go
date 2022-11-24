package misc

import "testing"

func TestReadURLToSrcData(t *testing.T) {
	url := "https://gchat.qpic.cn/gchatpic_new/2899929243/3956353663-2897479574-4DBAF36D8DB12F0045F779E4F3F8B12A/0?term=3&is_origin=0"
	data, ty, err := ReadURLToSrcData(url)
	if err != nil {
		t.Skip(err)
	}
	t.Log(data, ty)
}
