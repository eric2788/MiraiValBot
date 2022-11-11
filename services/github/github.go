package github

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/eric2788/MiraiValBot/internal/file"
	gh "github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

var (
	logger = utils.GetModuleLogger("valbot.github")
	ctx    = context.Background()
	client *gh.Client
	config *file.GithubConfig
	lock   sync.Mutex
)

func Init() {
	githubYaml := file.ApplicationYaml.Github
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubYaml.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client = gh.NewClient(tc)
	config = &githubYaml
}

// VerifySuccess verify that the api is working
func VerifySuccess() error {
	_, _, err := client.Users.Get(ctx, "")
	return err
}

func DownloadFile(path string) ([]byte, error) {
	content, resp, err := client.Repositories.DownloadContents(ctx, config.Name, config.Repository, path, nil)
	logger.Debugf("response status of DownloadFile: %s", resp.Status)
	if err != nil || resp.StatusCode != 200 {
		// not found file
		if resp.StatusCode == 404 || strings.Contains(err.Error(), "no file named") {
			return nil, os.ErrNotExist
		} else if err != nil {
			return nil, err
		} else { // no err but response status is not 200
			if b, err := io.ReadAll(resp.Body); err == nil {
				return nil, errors.New(string(b))
			} else {
				logger.Errorf("无法解析错误信息: %v, 将返回 http status", err)
				return nil, errors.New(resp.Status)
			}
		}
	}
	return io.ReadAll(content)
}

func GetFileInfo(path string) (*gh.RepositoryContent, error) {
	f, _, resp, err := client.Repositories.GetContents(ctx, config.Name, config.Repository, path, nil)
	logger.Debugf("response status of GetFileInfo: %s", resp.Status)
	if err != nil {
		if resp.StatusCode == 404 {
			return nil, os.ErrNotExist
		}
		return nil, err
	}
	logger.Debugf("github api rate limit: %+v", resp.Rate)
	return f, nil
}

func ListDir(path string) ([]*gh.RepositoryContent, error) {
	_, dir, resp, err := client.Repositories.GetContents(ctx, config.Name, config.Repository, path, nil)
	logger.Debugf("response status of ListDir: %s", resp.Status)
	if err != nil {
		if resp.StatusCode == 404 {
			return nil, os.ErrNotExist
		}
		return nil, err
	}
	logger.Debugf("github api rate limit: %+v", resp.Rate)
	return dir, nil
}

func RemoveFile(path string) error {
	lock.Lock()
	defer lock.Unlock()
	content, err := GetFileInfo(path)
	if err != nil {
		return err
	} else if content == nil {
		return fmt.Errorf("file is nil")
	}
	var sha *string
	if content.GetType() == "dir" {
		return fmt.Errorf(path + " is dir")
	}
	logger.Debugf("SHA of %s is %q", path, content.GetSHA())
	sha = content.SHA
	author := generateAuthor()
	msg := fmt.Sprintf("remove %s from cloud at %s", path, time.Now().Format(time.RFC3339))
	_, resp, err := client.Repositories.DeleteFile(ctx, config.Name, config.Repository, path, &gh.RepositoryContentFileOptions{
		Author:    author,
		Committer: author,
		Message:   &msg,
		Branch:    &config.Branch,
		SHA:       sha,
	})
	logger.Debugf("response status of UpdateFile: %s", resp.Status)
	if err != nil {
		return err
	}
	logger.Debugf("github api rate limit: %+v", resp.Rate)
	logger.Infof("档案 %s 已成功从 github.com/%s/%s 移除 => 移除时间: %s, SHA: %s",
		path, config.Name, config.Repository,
		author.Date.Format(time.RFC3339),
		content.GetSHA(),
	)
	return nil
}

func UpdateFile(path string, file []byte) error {
	lock.Lock()
	defer lock.Unlock()
	content, err := GetFileInfo(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	var sha *string
	if content != nil {
		if content.GetType() == "dir" {
			return fmt.Errorf(path + " is dir")
		}
		sha = content.SHA
	}
	logger.Debugf("SHA of %s is %q", path, content.GetSHA())
	author := generateAuthor()
	msg := fmt.Sprintf("update %s to cloud with %d bytes at %s", path, len(file), time.Now().Format(time.RFC3339))
	result, resp, err := client.Repositories.UpdateFile(ctx, config.Name, config.Repository, path, &gh.RepositoryContentFileOptions{
		Author:    author,
		Committer: author,
		Message:   &msg,
		Branch:    &config.Branch,
		SHA:       sha,
		Content:   file,
	})
	logger.Debugf("response status of UpdateFile: %s", resp.Status)
	if err != nil {
		return err
	}
	logger.Debugf("github api rate limit: %+v", resp.Rate)
	logger.Infof("档案 %s 已成功上传到 github.com/%s/%s => 档案大小: %d Bytes, 上传时间: %s, SHA: %s => %s",
		path, config.Name, config.Repository,
		len(file),
		author.Date.Format(time.RFC3339),
		content.GetSHA(),
		result.GetSHA(),
	)
	return nil
}

func UploadFile(path string, file []byte) error {
	lock.Lock()
	defer lock.Unlock()
	author := generateAuthor()
	msg := fmt.Sprintf("upload %s to cloud with %d bytes at %s", path, len(file), time.Now().Format(time.RFC3339))
	content, resp, err := client.Repositories.CreateFile(ctx, config.Name, config.Repository, path, &gh.RepositoryContentFileOptions{
		Author:    author,
		Committer: author,
		Message:   &msg,
		Branch:    &config.Branch,
		Content:   file,
	})
	logger.Debugf("response status of UploadFile: %s", resp.Status)
	if err != nil {
		if resp.StatusCode == 422 && strings.Contains(err.Error(), "\"sha\" wasn't supplied") {
			return os.ErrExist
		}
		return err
	}
	logger.Debugf("github api rate limit: %+v", resp.Rate)
	logger.Infof("档案 %s 已成功上传到 github.com/%s/%s => 档案大小: %d Bytes, 上传时间: %s, SHA: %s",
		path, config.Name, config.Repository,
		len(file),
		author.Date.Format(time.RFC3339),
		content.GetSHA(),
	)
	return nil
}

func generateAuthor() *gh.CommitAuthor {
	name, email, now := config.Name, config.Email, time.Now()
	return &gh.CommitAuthor{
		Name:  &name,
		Email: &email,
		Date:  &now,
	}
}
