package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/RomiChan/protobuf/proto"
	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/utils/compress"
	rgo "github.com/go-redis/redis/v8"
)

var (
	rdb    *rgo.Client
	ctx    = context.Background()
	logger = utils.GetModuleLogger("valbot.redis")
)

var posArg rgo.LPosArgs

const (
	Permanent   = time.Duration(0)
	ShortMoment = time.Minute * 10
	OneDay      = time.Hour * 24
)

func Init() {
	redisConfig := file.ApplicationYaml.Redis
	host := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)
	rdb = rgo.NewClient(&rgo.Options{
		Addr:         host,
		Password:     redisConfig.Password,
		DB:           redisConfig.Database,
		PoolSize:     64,
		MinIdleConns: 12,
		PoolTimeout:  5 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  1 * time.Minute,
	})
}

func Close() error {
	return rdb.Close()
}

func Subscribe(ctx context.Context, topic string) *rgo.PubSub {
	return rdb.Subscribe(ctx, topic)
}

func ListKeys(prefix string) (keys []string, err error) {
	keys, cursor, err := rdb.Scan(ctx, 0, prefix+":*", 0).Result()
	if cursor == 0 {
		return nil, nil
	}
	return
}

func HasKey(key string) (bool, error) {
	re, err := rdb.Exists(ctx, key).Result()
	return re == 1, err
}

func StoreTimely(key string, arg interface{}, duration time.Duration) error {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(arg)
	if err != nil {
		return err
	}
	return StoreBytes(key, buffer.Bytes(), duration)
}

func Store(key string, arg interface{}) error {
	return StoreTimely(key, arg, Permanent)
}

func StoreTemp(key string, arg interface{}) error {
	return StoreTimely(key, arg, ShortMoment)
}

func StoreProtoTimely(key string, arg interface{}, duration time.Duration) error {
	b, err := proto.Marshal(arg)
	if err != nil {
		return err
	}
	return StoreBytes(key, b, duration)
}

func StoreProto(key string, arg interface{}) error {
	return StoreProtoTimely(key, arg, Permanent)
}

func StoreProtoTemp(key string, arg interface{}) error {
	return StoreProtoTimely(key, arg, ShortMoment)
}

func StoreBytes(key string, data []byte, duration time.Duration) error {
	logger.Debugf("original data size of %s: %dB", key, len(data))
	data = compress.DoCompress(data)
	logger.Debugf("compressed data size of %s: %dB", key, len(data))
	return rdb.Set(ctx, key, data, duration).Err()
}

func GetBytes(key string) ([]byte, bool, error) {
	b, err := rdb.Get(ctx, key).Bytes()
	if b != nil {
		logger.Debugf("compressed data size of %s: %dB", key, len(b))
		b = compress.DoUnCompress(b)
		logger.Debugf("original data size of %s: %dB", key, len(b))
	}
	return b, err == rgo.Nil, err
}

func SetAdd(key string, value interface{}) error {
	return rdb.SAdd(ctx, key, value).Err()
}

func SetRemove(key string, value interface{}) error {
	return rdb.SRem(ctx, key, value).Err()
}

func ListPos(key, value string) (int64, error) {
	index, err := rdb.LPos(ctx, key, value, posArg).Result()
	if err == rgo.Nil {
		return -1, nil
	} else {
		return index, err
	}
}

func GetMapValue(key, mapKey string) (string, error) {
	value, err := rdb.HGet(ctx, key, mapKey).Result()
	if err == rgo.Nil {
		return "", nil
	} else {
		return value, err
	}
}

var ListExists = errors.New("this key in list exists")

func ListAdd(key, value string) error {
	i, err := rdb.LPos(ctx, key, value, posArg).Result()
	if err != nil {
		if err != rgo.Nil {
			return err
		}
	} else if i >= 0 {
		return nil
	}
	return rdb.RPush(ctx, key, value).Err()
}

func ListIndex(key string, index int64) (string, error) {
	s, err := rdb.LIndex(ctx, key, index).Result()
	if err == rgo.Nil {
		return "", nil
	} else {
		return s, err
	}
}

func ListRem(key, value string) error {
	return rdb.LRem(ctx, key, 1, value).Err()
}

func SetContains(key string, value interface{}) (bool, error) {
	return rdb.SIsMember(ctx, key, value).Result()
}

func Delete(key string) error {
	return rdb.Del(ctx, key).Err()
}

func GetProto(key string, arg interface{}) (bool, error) {
	b, notExist, err := GetBytes(key)
	if notExist {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, proto.Unmarshal(b, arg)
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

func GetString(key string) (string, error) {
	s, err := rdb.Get(ctx, key).Result()
	if err == rgo.Nil {
		return "", nil
	} else {
		return s, err
	}
}
