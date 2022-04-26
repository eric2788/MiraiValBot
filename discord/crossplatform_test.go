package discord

import (
	"github.com/eric2788/MiraiValBot/file"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func aTestStartServe(t *testing.T) {
	file.GenerateConfig()
	file.LoadApplicationYaml()
	file.LoadStorage()
	Start()

	// Wait here until CTRL-C or other term signal is received.
	logger.Info("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	_ = client.Close()
}
