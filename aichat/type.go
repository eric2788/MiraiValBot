package aichat

import (
	"fmt"
	"strings"

	"github.com/eric2788/common-utils/request"
)




type AIReply interface {
	Reply(msg string) (string, error)
	Name() string
}

func getAiReply(url, msg string) ([]byte, error) {
	return request.GetBytesByUrl(fmt.Sprintf(url, msg))
}

func replaces(msg string, replacer map[string]string) (string){
	for src, dst := range replacer {
		msg = strings.ReplaceAll(msg, src, dst)
	}
	return msg
}