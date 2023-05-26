package command

import (
	"fmt"
	"github.com/eric2788/common-utils/stream"
	"strings"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/common-utils/array"
)

const Prefix = "!"

type MessageSource struct {
	Client  *client.QQClient
	Message *message.GroupMessage
}

type Response struct {
	Ignore   bool
	Content  string
	ShowHelp bool
}

type CmdHandler func(args []string, source *MessageSource) error

type Node struct {
	Command      string
	Alias        []string
	ChildNodes   []Node
	AdminOnly    bool
	Placeholders []string
	Description  string
	Handler      *CmdHandler
}

var commandTree []Node

func AddCommand(node Node) {
	commandTree = append(commandTree, node)
	logger.Infof("已成功註冊指令 %s 且其 %d 個子指令", node.Command, len(node.ChildNodes))
}

// NewParent 新增含分支指令的指令節點
func NewParent(names []string, description string, nodes ...Node) Node {
	if len(names) == 0 {
		panic("指令名稱最少需要一個參數")
	}

	cmd, alias := names[0], names[1:]
	return Node{
		Command:      cmd,
		Alias:        alias,
		ChildNodes:   nodes,
		AdminOnly:    false,
		Placeholders: []string{},
		Description:  description,
		Handler:      nil,
	}
}

// NewNode 新增指令處理器
func NewNode(names []string, description string, adminOnly bool, handler CmdHandler, placeholders ...string) Node {
	if len(names) == 0 {
		panic("指令名稱最少需要一個參數")
	}
	cmd, alias := names[0], names[1:]

	return Node{
		Command:      cmd,
		Alias:        alias,
		ChildNodes:   []Node{},
		AdminOnly:    adminOnly,
		Placeholders: placeholders,
		Description:  description,
		Handler:      &handler,
	}
}

func InvokeCommand(content string, admin bool, source *MessageSource) (*Response, error) {
	if isNotCommand(content) {
		return &Response{
			Ignore: true,
		}, nil
	}
	logger.Debugf("收到指令輸入: %s", content)
	commands := strings.Split(strings.TrimPrefix(content, Prefix), " ")
	cmd, plainArgs := commands[0], commands[1:]

	var args []string

	// remove all empty or space args
	for _, arg := range plainArgs {
		if strings.TrimSpace(arg) != "" {
			args = append(args, arg)
		}
	}

	logger.Debugf("original args(%d): %s", len(plainArgs), strings.Join(plainArgs, ", "))
	logger.Debugf("filtered args(%d): %s", len(args), strings.Join(args, ", "))

	for _, node := range commandTree {
		labels := append(node.Alias, node.Command)

		logger.Debugf("指令对比: %s vs %v", cmd, labels)

		if array.Contains(labels, cmd) {
			return invokeCommandInternal(node, admin, []string{}, args, source)
		}
	}
	return &Response{
		Content:  showHelpLines([]string{}, commandTree),
		ShowHelp: true,
	}, nil
}

func invokeCommandInternal(node Node, admin bool, parents []string, args []string, source *MessageSource) (*Response, error) {

	// 非管理員且為管理員指令
	if node.AdminOnly && !admin {
		return &Response{
			Content:  "只有管理員才可使用此指令",
			ShowHelp: false,
		}, nil
	}

	// 有子指令
	if len(node.ChildNodes) > 0 {

		parents = append(parents, node.Command)

		// 有參數輸入
		if len(args) > 0 {
			cmd, args := args[0], args[1:]
			for _, subNode := range node.ChildNodes {
				labels := append(subNode.Alias, subNode.Command)
				logger.Debugf("指令对比: %s vs %v, 参数: %v", cmd, labels, args)
				if array.Contains(labels, cmd) {
					return invokeCommandInternal(subNode, admin, parents, args, source)
				}
			}
		}

		// 找不到匹配，返回幫助
		return &Response{
			Content:  showHelpLines(parents, node.ChildNodes),
			ShowHelp: true,
		}, nil
	}

	// 參數過少
	if len(filterNecessary(node.Placeholders)) > len(args) {
		return &Response{
			Content:  showHelpLine(parents, node),
			ShowHelp: true,
		}, nil
	}

	if node.Handler == nil {
		return nil, fmt.Errorf("此指令的處理器為 NULL")
	}

	err := (*node.Handler)(args, source)
	return &Response{}, err
}

func isNotCommand(content string) bool {
	return !strings.HasPrefix(content, Prefix)
}

func showHelpLines(parents []string, nodes []Node) string {
	var lines []string
	for _, node := range nodes {
		subCommandLine := showHelpLine(parents, node)
		lines = append(lines, subCommandLine)
	}
	return strings.Join(lines, "\n")
}

func filterNecessary(placeholders []string) []string {
	return stream.From(placeholders).
		Filter(func(s string) bool {
			return !(strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]"))
		}).
		ToArr()
}

func showHelpLine(parents []string, node Node) string {
	line := Prefix
	for _, parent := range parents {
		line += parent + " "
	}
	line += node.Command
	for _, placeholder := range node.Placeholders {
		line += " " + placeholder
	}
	return fmt.Sprintf("%s - %s", line, node.Description)
}

func ExtractPrefix(line string) string {
	content := make([]string, 0)
	for _, words := range strings.Split(line, " ") {
		if strings.HasPrefix(words, Prefix) {
			continue
		}
		content = append(content, words)
	}
	return strings.Join(content, " ")
}
