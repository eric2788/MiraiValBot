package main

import (
	"github.com/eric2788/MiraiValBot/utils/test"
	"os"
	"os/signal"
	"syscall"

	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/services/discord"
	"github.com/sirupsen/logrus"
)

var logger = logrus.WithField("simulate", "discord")

func main() {

	test.InitTesting()

	file.GenerateConfig()
	file.LoadApplicationYaml()
	file.LoadStorage()
	discord.Start()

	// Wait here until CTRL-C or other term signal is received.
	logger.Info("Discord Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	_ = discord.Close()
}
