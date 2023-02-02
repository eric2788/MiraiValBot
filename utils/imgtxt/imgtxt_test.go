package imgtxt

import (
	"fmt"
	"os"
	"testing"

	"github.com/eric2788/common-utils/request"
	"github.com/google/uuid"
)

func TestGeneratePlayerImage(t *testing.T) {
	msg, err := NewPrependMessage()

	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		msg.Append("=========這是繁體=========")
		msg.Append(fmt.Sprintf("第 %d 名: %s", i+1, "ABC"))

		// 基本资料
		msg.Append("\t基本资料:")
		msg.Append("\t\tKDA: %d | %d | %d")
		msg.Append("\t\t分数: %d")
		msg.Append("\t\t使用角色: %s")
		msg.Append("\t\t所在队伍: %s")

		// 击中分布
		msg.Append("\t击中次数分布")
		msg.Append("\t\t头部: %.1f%% (%d次)")
		msg.Append("\t\t身体: %.1f%% (%d次)")
		msg.Append("\t\t腿部: %.1f%% (%d次)")

		// 行为
		msg.Append("\t行为:")
		msg.Append("\t\tAFK回合次数: %.0f")
		msg.Append("\t\t误击队友伤害: %d")
		msg.Append("\t\t误杀队友次数: %d")
		msg.Append("\t\t被误击队友伤害: %d")
		msg.Append("\t\t被误杀队友次数: %d")
		msg.Append("\t\t拆包次数: %d")
		msg.Append("\t\t装包次数: %d")

		//技能使用

		msg.Append("\t技能使用次数分布:")
		msg.Append("\t\t技能 Q: %d次 (%.1f%%)")
		msg.Append("\t\t技能 E: %d次 (%.1f%%)")
		msg.Append("\t\t技能 C: %d次 (%.1f%%)")
		msg.Append("\t\t技能 X: %d次 (%.1f%%)")

		// 经济
		msg.Append("\t经济:")
		msg.Append("\t\t总支出 $%d")
		msg.Append("\t\t平均支出 $%d")

		// 伤害
		msg.Append("\t伤害分布:")
		msg.Append("\t\t总承受 %d (%.1f%%)")
		msg.Append("\t\t总伤害 %d (%.1f%%)")

	}

	b, err := msg.GenerateImage()
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Create("test.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if _, err := f.Write(b); err != nil {
		t.Fatal(err)
	}
}

func TestGenerateLeaderboardImage(t *testing.T) {

	msg, err := NewPrependMessage()
	if err != nil {
		t.Fatal(err)
	}

	msg.Append("名次\t\t玩家\t\t均分\tK\tD\tA\t爆头率\t友伤\t装包\t拆包")
	for i := 0; i < 10; i++ {
		msg.Append(fmt.Sprintf("%d\t\t%s\t\t%d\t%d\t%d\t%d\t%.1f%%\t%d\t%d\t%d",
			i+1, uuid.New().String()[:10], 50, 1, 2, 3, float64(20), 4, 5, 16,
		))
	}

	b, err := msg.GenerateImage()
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Create("test.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if _, err := f.Write(b); err != nil {
		t.Fatal(err)
	}
}

func TestGenerateImageInImage(t *testing.T) {
	msg, err := NewPrependMessage()
	if err != nil {
		t.Fatal(err)
	}

	msg.Append("hello world!\n")
	msg.Append("1234567\n")
	msg.Append("789456123\n")
	msg.Append("Image:\n")

	img, err := request.GetBytesByUrl("https://img.freepik.com/free-vector/abstract-blue-modern-elegant-design-background_1017-32105.jpg")
	if err != nil {
		t.Skip(err)
	}

	msg.AppendImage(img)

	b, err := msg.GenerateImage()
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Create("test.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if _, err := f.Write(b); err != nil {
		t.Fatal(err)
	}
}
