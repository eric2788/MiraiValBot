package test

import (
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
	"path/filepath"
	"runtime"
)

var logger = utils.GetModuleLogger("utils.test")

func InitTesting() {
	logrus.SetLevel(logrus.DebugLevel)
	logger.Debugf("Logging Level set to debug")
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		logger.Warnf("unable to get the current filename")
		return
	}
	dirname := filepath.Dir(filename)
	if err := gotenv.OverLoad(dirname + "\\.env.local"); err == nil {
		logger.Debugf("successfully loaded local environment variables.")
	}
}
