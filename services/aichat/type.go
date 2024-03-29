package aichat

import (
	"strings"

	"github.com/eric2788/common-utils/request"
)

type AIReply interface {
	Reply(msg string) (string, error)
	Name() string
}

func getAiReply(url, msg string) ([]byte, error) {
	msg = strings.ReplaceAll(msg, " ", "")
	url = strings.ReplaceAll(url, "{msg}", msg)
	return request.GetBytesByUrl(url)
}

func replaces(msg string, replacer map[string]string) string {
	for src, dst := range replacer {
		msg = strings.ReplaceAll(msg, src, dst)
	}
	return strings.TrimSpace(msg)
}
