package github

import (
	"os"
	"testing"
	"time"

	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/utils/compress"
	gh "github.com/google/go-github/v48/github"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGithubAccess(t *testing.T) {
	if file.ApplicationYaml.Github.AccessToken == "" {
		t.Log("token is empty, skipped test.")
		return
	}
	user, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	// Rate.Limit should most likely be 5000 when authorized.
	t.Logf("Rate: %#v\n", resp.Rate)

	// If a Token Expiration has been set, it will be displayed.
	if !resp.TokenExpiration.IsZero() {
		t.Logf("Token Expiration: %v\n", resp.TokenExpiration)
	}

	t.Logf("\n%v\n", gh.Stringify(user))
}

func TestGithubRepo(t *testing.T) {
	if file.ApplicationYaml.Github.AccessToken == "" {
		t.Log("token is empty, skipped test.")
		return
	}
	_, files, resp, err := client.Repositories.GetContents(ctx, "sysnapse", "cloud", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Rate: %+v", resp.Rate)
	for _, file := range files {
		t.Logf("Name: %s", file.GetName())
		t.Logf("Path: %s", file.GetPath())
		t.Logf("Download: %s", file.GetDownloadURL())
		t.Logf("Type: %s", file.GetType())

		if file.GetType() != "file" {
			continue
		}

		f, _, resp, err := client.Repositories.GetContents(ctx, "sysnapse", "cloud", file.GetPath(), nil)
		if err != nil {
			t.Log(err)
		} else {
			t.Logf("Rate: %+v", resp.Rate)
			if content, err := f.GetContent(); err == nil {
				t.Logf("Content: %s", content)
			} else {
				t.Logf("Cannot get content: %s", err)
			}
		}
	}
}

func TestGithubUploadAndRead(t *testing.T) {
	if file.ApplicationYaml.Github.AccessToken == "" {
		t.Log("token is empty, skipped test.")
		return
	}
	path := "test/upload.txt"
	text := "this is a file that uploaded from github v3 api, date: " + time.Now().Format(time.UnixDate)
	err := UpdateFile(path, []byte(text))
	if err != nil {
		t.Fatal(err)
	}
	file, err := GetFileInfo(path)
	if err != nil {
		t.Fatal(err)
	}
	if c, err := file.GetContent(); err != nil {
		t.Fatal(err)
	} else {
		t.Logf(c)
		assert.Equal(t, c, text)
	}
}

func TestUploadReadCompressed(t *testing.T) {
	if file.ApplicationYaml.Github.AccessToken == "" {
		t.Log("token is empty, skipped test.")
		return
	}
	compress.SwitchType("zlib")
	path := "test/binaryfile"
	text := "this is a binary file that compressed and uploaded from github v3 api, date: " + time.Now().Format(time.UnixDate)
	data := []byte(text)
	compressed := compress.DoCompress(data)

	t.Logf("before size: %d, after size: %d", len(data), len(compressed))

	err := UpdateFile(path, compressed)
	if err != nil {
		t.Fatal(err)
	}

	b, err := DownloadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	result := compress.DoUnCompress(b)
	t.Logf("compressed size: %d, uncompressed size: %d", len(b), len(result))

	assert.Equal(t, text, string(data))
	t.Logf("text: %s", text)
	t.Logf("result: %s", result)
	assert.Equal(t, text, string(result))
}

func TestUploadExistError(t *testing.T) {
	if file.ApplicationYaml.Github.AccessToken == "" {
		t.Log("token is empty, skipped test.")
		return
	}
	path := "test/test_err.txt"
	// upload first
	_ = UploadFile(path, []byte("hawidhaiwhdiahdiawhida"))
	err := UploadFile(path, []byte("hawidhaiwhdiahdiawhida"))
	assert.NotNil(t, err)
	assert.Equal(t, os.ErrExist, err)
	t.Log(err)
}

func TestNotExistError(t *testing.T) {
	if file.ApplicationYaml.Github.AccessToken == "" {
		t.Log("token is empty, skipped test.")
		return
	}
	_, err1 := GetFileInfo("test/awdawdawdawda")
	_, err2 := DownloadFile("test/awdawdawdadaw")

	assert.NotNil(t, err1)
	assert.NotNil(t, err2)

	assert.Equal(t, err1, err2)
	assert.Equal(t, err1, os.ErrNotExist)
}

func TestListDir(t *testing.T) {
	if file.ApplicationYaml.Github.AccessToken == "" {
		t.Log("token is empty, skipped test.")
		return
	}
	dir, err := ListDir("test")
	if err != nil {
		t.Fatal(err)
	}

	for _, d := range dir {
		t.Logf("Name: %s, Type: %s, Path: %s", d.GetName(), d.GetType(), d.GetPath())
	}
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	Init()
}
