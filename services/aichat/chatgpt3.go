package aichat

import (
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"os"
)

// openAI
type Chatgpt3 struct {
}

func (c *Chatgpt3) Reply(msg string) (string, error) {
	apiKey := os.Getenv("CHATGPT_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("chatgpt api key not set")
	}
	cli := gpt3.NewClient(apiKey,
		gpt3.WithDefaultEngine(gpt3.TextDavinci003Engine),
	)

	resp, err := cli.Completion(context.Background(), gpt3.CompletionRequest{
		Prompt:      []string{msg},
		MaxTokens:   gpt3.IntPtr(1000),
		Temperature: gpt3.Float32Ptr(0),
		N:           gpt3.IntPtr(1),
	})

	if err != nil {
		return "", err
	}

	return replaces(resp.Choices[0].Text, map[string]string{
		"答：": "",
	}), nil
}

func (c *Chatgpt3) Name() string {
	return "chatgpt3"
}
