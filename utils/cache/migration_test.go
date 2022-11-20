package cache

import (
	"fmt"
	"github.com/eric2788/MiraiValBot/utils/test"
	"testing"
	"time"

	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/services/github"
)

func init() {
	test.InitTesting()
	github.Init()
}

func TestMigration(t *testing.T) {

	if file.ApplicationYaml.Github.AccessToken == "" {
		logger.Debugf("skipping migration test because no github access token is provided")
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
