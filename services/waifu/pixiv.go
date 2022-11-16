package waifu

import (
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/everpcpc/pixiv"
)

var (
	app    = pixiv.NewApp()
	logger = utils.GetModuleLogger("services.waifu")
)

func Init() {
	account, err := pixiv.Login(file.ApplicationYaml.Pixiv.Username, file.ApplicationYaml.Pixiv.Password)
	if err != nil {
		logger.Errorf("登入 pixiv 失敗: %v", err)
	} else {
		logger.Infof("成功登入: %s", account.Name)
	}
}

func getIllust(id uint64) (*pixiv.Illust, error) {
	return app.IllustDetail(id)
}
