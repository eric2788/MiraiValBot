package waifu

import (
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/everpcpc/pixiv"
	"os"
	"time"
)

var (
	app    *pixiv.AppPixivAPI
	logger = utils.GetModuleLogger("services.waifu")
)

func Init() {
	account, err := pixiv.Login(file.ApplicationYaml.Pixiv.Username, file.ApplicationYaml.Pixiv.Password)
	if err != nil {
		logger.Errorf("Password 登入 pixiv 失敗: %v", err)
		if os.Getenv("PIXIV_TOKEN") != "" && os.Getenv("PIXIV_REFRESH_TOKEN") != "" {
			logger.Infof("將嘗試使用環境變數中的 pixiv token 進行登入")
			account, err = pixiv.LoadAuth(os.Getenv("PIXIV_TOKEN"), os.Getenv("PIXIV_REFRESH_TOKEN"), time.Now().Add(time.Second*3600))
			if err != nil {
				logger.Errorf("Token 登入 pixiv 失敗: %v", err)
			}
		}
	}

	if account != nil {
		logger.Infof("成功登入: %s", account.Name)
	}
	app = pixiv.NewApp()
}

func getIllust(id uint64) (*pixiv.Illust, error) {
	return app.IllustDetail(id)
}
