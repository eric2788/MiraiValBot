package cache

import (
	"fmt"
	"os"

	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/subosito/gotenv"
)

const StrategyEnvVar = "CACHE_STRATEGY"

var logger = utils.GetModuleLogger("valbot.cache")

type (
	CachedData struct {
		Name string
		Path string
	}

	Service struct {
		path string
		db   CacheDatabase
	}

	CacheDatabase interface {
		Init(path string) error
		Save(path, name string, data []byte) error
		Get(path, name string) ([]byte, error)
		Remove(path, name string) error
		List(path string) []CachedData
	}

	Options struct {
		Type string
		Path string
	}
)

var cacheDbMap = map[string]CacheDatabase{
	"local":  &LocalCache{},
	"github": &GitCache{},
}

func (s *Service) Set(name string, data []byte) error {
	return s.db.Save(s.path, name, data)
}

func (s *Service) Get(name string) ([]byte, error) {
	return s.db.Get(s.path, name)
}

func (s *Service) Remove(name string) error {
	return s.db.Remove(s.path, name)
}

func (s *Service) List() []CachedData {
	return s.db.List(s.path)
}

func New(options ...func(*Options)) (*Service, error) {
	opt := &Options{
		Path: "",
		Type: "local",
	}
	for _, setting := range options {
		setting(opt)
	}
	if db, ok := cacheDbMap[opt.Type]; ok {
		if err := db.Init(opt.Path); err != nil {
			return nil, err
		}
		return &Service{
			db:   db,
			path: opt.Path,
		}, nil
	} else {
		return nil, fmt.Errorf("未知的缓存方式: %s", opt.Type)
	}
}

func NewCache(path string) *Service {
	t := os.Getenv(StrategyEnvVar)
	ser, err := New(
		WithPath(path),
		WithType(t),
	)
	if err != nil {
		panic(err)
	}
	return ser
}

func WithPath(path string) func(*Options) {
	return func(opt *Options) {
		opt.Path = path
	}
}

func WithType(t string) func(*Options) {
	return func(opt *Options) {
		opt.Type = t
	}
}

func init() {
	if err := gotenv.Load(".env.local"); err == nil {
		logger.Debugf("successfully loaded local environment variables.")
	}
	if os.Getenv(StrategyEnvVar) == "" {
		logger.Warn("緩存方式沒有設置, 已默認改用local。")
		os.Setenv(StrategyEnvVar, "local")
	}
}
