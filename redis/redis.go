package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"github.com/eric2788/MiraiValBot/file"
	rgo "github.com/go-redis/redis/v8"
	"time"
)

var rdb *rgo.Client
var ctx = context.Background()

func Init() {
	redisConfig := file.ApplicationYaml.Redis
	host := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)
	rdb = rgo.NewClient(&rgo.Options{
		Addr:     host,
		Password: redisConfig.Password,
		DB:       redisConfig.Database,
	})
}

func Close() error {
	return rdb.Close()
}

func Subscribe(ctx context.Context, topic string) *rgo.PubSub {
	return rdb.Subscribe(ctx, topic)
}

func Store(key string, arg interface{}) error {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(arg)
	if err != nil {
		return err
	}
	return StoreBytes(key, buffer.Bytes())
}

func StoreBytes(key string, data []byte) error {
	return rdb.Set(ctx, key, data, time.Hour*86400).Err()
}

func GetBytes(key string) ([]byte, bool, error) {
	b, err := rdb.Get(ctx, key).Bytes()
	return b, err == rgo.Nil, err
}

func Get(key string, arg interface{}) (bool, error) {
	b, notExist, err := GetBytes(key)
	if notExist {
		return false, nil
	} else if err != nil {
		return false, err
	}
	buffer := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buffer)
	return true, dec.Decode(arg)
}
