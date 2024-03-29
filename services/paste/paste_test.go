package paste

import (
	"os"
	"strings"
	"testing"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func generateTestMessage() *message.SendingMessage {
	msg := message.NewSendingMessage()

	for i := 0; i < 10; i++ {
		msg.Append(qq.NewTextLn("=========這是繁體========="))
		msg.Append(qq.NewTextfLn("第 %d 名: %s", i+1))

		// 基本资料
		msg.Append(qq.NewTextLn("\t基本资料:"))
		msg.Append(qq.NewTextfLn("\t\tKDA: %d | %d | %d"))
		msg.Append(qq.NewTextfLn("\t\t分数: %d"))
		msg.Append(qq.NewTextfLn("\t\t使用角色: %s"))
		msg.Append(qq.NewTextfLn("\t\t所在队伍: %s"))

		// 击中分布
		msg.Append(qq.NewTextLn("\t击中次数分布"))
		msg.Append(qq.NewTextfLn("\t\t头部: %.1f%% (%d次)"))
		msg.Append(qq.NewTextfLn("\t\t身体: %.1f%% (%d次)"))
		msg.Append(qq.NewTextfLn("\t\t腿部: %.1f%% (%d次)"))

		// 行为
		msg.Append(qq.NewTextLn("\t行为:"))
		msg.Append(qq.NewTextfLn("\t\tAFK回合次数: %.0f"))
		msg.Append(qq.NewTextfLn("\t\t误击队友伤害: %d"))
		msg.Append(qq.NewTextfLn("\t\t误杀队友次数: %d"))
		msg.Append(qq.NewTextfLn("\t\t被误击队友伤害: %d"))
		msg.Append(qq.NewTextfLn("\t\t被误杀队友次数: %d"))
		msg.Append(qq.NewTextfLn("\t\t拆包次数: %d"))
		msg.Append(qq.NewTextfLn("\t\t装包次数: %d"))

		//技能使用

		msg.Append(qq.NewTextLn("\t技能使用次数分布:"))
		msg.Append(qq.NewTextfLn("\t\t技能 Q: %d次 (%.1f%%)"))
		msg.Append(qq.NewTextfLn("\t\t技能 E: %d次 (%.1f%%)"))
		msg.Append(qq.NewTextfLn("\t\t技能 C: %d次 (%.1f%%)"))
		msg.Append(qq.NewTextfLn("\t\t技能 X: %d次 (%.1f%%)"))

		// 经济
		msg.Append(qq.NewTextLn("\t经济:"))
		msg.Append(qq.NewTextfLn("\t\t总支出 $%d"))
		msg.Append(qq.NewTextfLn("\t\t平均支出 $%d"))

		// 伤害
		msg.Append(qq.NewTextLn("\t伤害分布:"))
		msg.Append(qq.NewTextfLn("\t\t总承受 %d (%.1f%%)"))
		msg.Append(qq.NewTextfLn("\t\t总伤害 %d (%.1f%%)"))

	}
	return msg
}

func TestCreatePasteMe(t *testing.T) {

	msg := generateTestMessage()
	content := strings.Join(qq.ParseMsgContent(msg.Elements).Texts, "")

	url, err := CreatePasteMe("plain", content)
	if err != nil {
		t.Skip(err)
	}
	t.Logf(url)
}

func TestCreatePasteBin(t *testing.T) {
	msg := generateTestMessage()
	content := strings.Join(qq.ParseMsgContent(msg.Elements).Texts, "")

	if os.Getenv("PASTEBIN_API_KEY") == "" {
		t.Logf("API key is not set, skipped the test.")
		return
	}

	url, err := CreatePasteBin("test-message", content, "yaml")

	if err != nil {
		t.Fatal(err)
	}

	t.Log(url)
}
