package cache

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/services/github"
	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	if err := gotenv.Load("../../services/github/.env.local"); err == nil {
		logger.Debugf("successfully loaded local environment variables.")
	}
	file.ApplicationYaml.Github.AccessToken = os.Getenv("GITHUB_PAT_TOKEN")
	github.Init()
}

func TestMigration(t *testing.T) {

	if file.ApplicationYaml.Github.AccessToken == "" {
		return
	}

	git, err := New(
		WithType("github"),
		WithPath("test_migration"),
	)
	if err != nil {
		t.Fatal(err)
	}
	local, err := New(
		WithType("local"),
		WithPath("test_migration"),
	)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 15; i++ {
		id := fmt.Sprint(i)
		go func(id string) {
			if err := local.Set(id, []byte(id)); err != nil {
				t.Log(err)
			}
		}(id)
	}

	<-time.After(time.Second * 3)
	Migrate(local, git, true).Wait()
}
