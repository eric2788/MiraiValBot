package redis

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	rgo "github.com/go-redis/redis/v8"
	"strconv"
	"testing"
	"time"
)

type Class struct {
	Name    string
	Teacher Teacher
	Student []Student
}

type Student struct {
	Name string
	SID  int
	Age  int
}

type Teacher struct {
	Name string
	Age  int
}

var testData = Class{
	Name: "Test Class",
	Teacher: Teacher{
		Name: "Teacher T",
		Age:  42,
	},
	Student: []Student{
		{
			Name: "Student A",
			SID:  1,
			Age:  16,
		},
		{
			Name: "Student B",
			SID:  2,
			Age:  17,
		},
		{
			Name: "Student 3",
			SID:  3,
			Age:  18,
		},
	},
}

var testOption = &rgo.Options{
	Addr:     "192.168.0.127:6379",
	Password: "",
	DB:       1,
}

func TestStoreBinary(t *testing.T) {

	key := "test:test_binary"

	rdb = rgo.NewClient(testOption)

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	if err := enc.Encode(testData); err != nil {
		t.Fatal(err)
	}

	if err := rdb.Set(ctx, key, buffer.Bytes(), time.Minute*30).Err(); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("saved %+v\n", buffer)

	b, err := rdb.Get(ctx, key).Bytes()

	if err != nil {
		t.Fatal(err)
	}

	data := &Class{}
	dBuffer := bytes.NewBuffer(b)
	dec := gob.NewDecoder(dBuffer)
	if err = dec.Decode(data); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Parsed %+v\n", data)
}

func TestStoreByteArray(t *testing.T) {

	key := "test:test_byte_arr"

	rdb = rgo.NewClient(testOption)

	b, err := json.Marshal(testData)

	if err != nil {
		t.Fatal(err)
	}

	for i, bb := range b {
		rdb.HSet(ctx, key, i, bb)
	}

	fmt.Printf("saved %v\n", b)

	bMap, err := rdb.HGetAll(ctx, key).Result()

	if err != nil {
		t.Fatal(err)
	}

	bArr := make([]byte, len(bMap))
	for k, v := range bMap {
		i, e1 := strconv.Atoi(k)
		b, e2 := strconv.Atoi(v)
		if e1 != nil || e2 != nil {
			t.Fatal(e1, e2)
		}
		bArr[i] = byte(b)
	}

	fmt.Printf("got %v\n", bArr)

	data := &Class{}
	if err = json.Unmarshal(bArr, data); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Parsed: %+v\n", data)
}
