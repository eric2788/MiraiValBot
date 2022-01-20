package redis

import (
	"fmt"
	rgo "github.com/go-redis/redis/v8"
	"testing"
)

func SaveArrToByte(t *testing.T) {
	rdb = rgo.NewClient(&rgo.Options{
		Addr:     "192.168.0.127:6379",
		Password: "",
		DB:       1,
	})

	data := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	if err := Store("test:test_arr", data); err != nil {
		t.Fatal(err)
	}

	getter := make([]int64, 0)

	if b, _, err := GetBytes("test:test_arr"); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(string(b))
	}

	if exist, err := Get("test:test_arr", &getter); err != nil {
		t.Fatal(err)
	} else if !exist {
		t.Fatal("not exist")
	} else {
		fmt.Println(getter)
	}
}
