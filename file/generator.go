package file

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"gopkg.in/yaml.v3"
	"os"
)

type Configuration struct {
	Bot     BotConfig     `yaml:"bot"`
	Val     ValConfig     `yaml:"val"`
	Redis   RedisConfig   `yaml:"redis"`
	Discord DiscordConfig `yaml:"discord"`
}

type BotConfig struct {
	Account  int64  `yaml:"account"`
	Password string `yaml:"password"`
}

type ValConfig struct {
	GroupId          int64  `yaml:"groupId"`
	TwitterLookUpUrl string `yaml:"twitterLookUpUrl"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database int    `yaml:"database"`
	Password string `yaml:"password"`
}

type DiscordConfig struct {
	Token       string `yaml:"token"`
	Guild       int64  `yaml:"guild"`
	LogChannel  int64  `yaml:"logChannel"`
	NewsChannel int64  `yaml:"newsChannel"`
}

var defaultConfig = Configuration{
	Bot: BotConfig{
		Account:  123456789,
		Password: "password",
	},
	Val: ValConfig{
		GroupId:          123456789,
		TwitterLookUpUrl: "http://172.17.0.1:8989",
	},
	Redis: RedisConfig{
		Host:     "172.17.0.1",
		Port:     6379,
		Database: 1,
		Password: "",
	},
}

func generate(filename string, generateFunc func() error) {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {

			fmt.Printf("檢測不到 %s，正在生成文件\n", filename)

			if err = generateFunc(); err != nil {
				fmt.Printf("%s 生成失敗: %v\n", filename, err)
				os.Exit(1)
			} else {
				fmt.Printf("已成功生成默認的 %s\n", filename)
			}
		} else {
			panic(fmt.Sprintf("檢測 %s 失敗: %v", filename, err))
		}
	}
}

func GenerateConfig() {
	generate("application.yaml", func() error {
		content, err := yaml.Marshal(defaultConfig)

		if err != nil {
			panic(fmt.Sprintf("解析 Yaml 失敗: %v", err))
		}

		return os.WriteFile("application.yaml", content, 0755)

	})
}

func GenerateDevice() {
	generate("device.json", func() error {
		bot.GenRandomDevice()
		return nil
	})
}
