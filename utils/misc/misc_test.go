package misc

import (
	"reflect"
	"testing"
)

func TestReadURLToSrcData(t *testing.T) {
	url := "https://gchat.qpic.cn/gchatpic_new/2899929243/3956353663-2897479574-4DBAF36D8DB12F0045F779E4F3F8B12A/0?term=3&is_origin=0"
	data, ty, err := ReadURLToSrcData(url)
	if err != nil {
		t.Skip(err)
	}
	t.Log(data, ty)
}


func TestXMLEscape(t *testing.T) {
	a := "<hello world&>"
	t.Log(XmlEscape(a))
}

type (
	common interface {
		Foo() string
	}
	foo struct {}
	bar struct {}
)

func (f foo) Foo() string {
	return "foo"
}

func (b bar) Foo() string {
	return "bar"
}

func TestGetTypeName(t *testing.T){
	var c common
	c = foo{}
	t.Log(reflect.TypeOf(c).Name())
	c = bar{}
	t.Log(reflect.TypeOf(c).Name())
}