package aichat

import (
	"context"
	"fmt"
	"os"

	"github.com/PullRequestInc/go-gpt3"
)

// openAI
type Chatgpt3 struct {
}

func (c *Chatgpt3) Reply(msg string) (string, error) {
	apiKey := os.Getenv("CHATGPT_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("chatgpt api key not set")
	}
	cli := gpt3.NewClient(apiKey)

	resp, err := cli.Completion(context.Background(), gpt3.CompletionRequest{
		Prompt:      []string{msg},
		MaxTokens:   gpt3.IntPtr(512),
		Temperature: gpt3.Float32Ptr(0),
	})

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Text, nil
}

func (c *Chatgpt3) Name() string {
	return "chatgpt3"
}
