package redis

import (
	"fmt"
	rgo "github.com/go-redis/redis/v8"
	"testing"
)

func ainit() {
	rdb = rgo.NewClient(&rgo.Options{
		Addr:     "192.168.0.127:6379",
		Password: "",
		DB:       1,
	})
}

func testList(v string) (int64, error) {
	key := "test:test_list"
	if err := ListAdd(key, v); err != nil && err != ListExists {
		return -1, err
	}
	return ListPos(key, v)
}

func testLists(vv []string) (map[string]int64, error) {
	key := "test:test_list"
	result := make(map[string]int64)
	for _, v := range vv {
		if err := ListAdd(key, v); err != nil {
			fmt.Println(err)
			result[v] = -1
		}
	}
	for _, v := range vv {
		if pos, err := ListPos(key, v); err != nil {
			fmt.Println(err)
			result[v] = -1
		} else {
			result[v] = pos
		}
	}
	return result, nil
}

func aTestSaveList(t *testing.T) {

	list := []string{"z", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	for _, v := range list {
		if index, err := testList(v); err != nil {
			t.Fatal(err)
		} else {
			fmt.Println(index)
		}
	}
}

func SaveArrToByte(t *testing.T) {
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
