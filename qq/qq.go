package qq

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/redis"
)

var ValGroupInfo = &client.GroupInfo{
	Uin:     file.ApplicationYaml.Val.GroupId,
	Members: []*client.GroupMemberInfo{},
}

type MsgContent struct {
	Texts   []string
	At      []int64
	Faces   []string
	Images  []string
	Replies []int32
}

func (msg *MsgContent) String() string {
	return strings.Join(msg.Texts, "")
}

var logger = utils.GetModuleLogger("valbot.qq")

func InitValGroupInfo(bot *bot.Bot) {

	ginfo, err := bot.GetGroupInfo(file.ApplicationYaml.Val.GroupId)

	if err != nil {
		logger.Fatalf("群Id %d 無效! ", file.ApplicationYaml.Val.GroupId)
		return
	}

	members, err := bot.GetGroupMembers(ginfo)

	if err != nil {
		logger.Fatalf("獲取群 %d 的成員列表失敗", file.ApplicationYaml.Val.GroupId)
		return
	}

	ginfo.Update(func(info *client.GroupInfo) {
		info.Members = members
		info.Uin = file.ApplicationYaml.Val.GroupId
	})

	ValGroupInfo = ginfo

	logger.Infof("以指定 %s (%d) 为 瓦群。(共 %d 個成員)", ValGroupInfo.Name, ValGroupInfo.Uin, len(ValGroupInfo.Members))
}

func RefreshGroupInfo() {
	ValGroupInfo.Update(func(info *client.GroupInfo) {
		ginfo, err := bot.Instance.GetGroupInfo(file.ApplicationYaml.Val.GroupId)
		if err != nil {
			logger.Warnf("刷新群资料时出现错误: %v", err)
		} else {
			ValGroupInfo.Name = ginfo.Name
			ValGroupInfo.GroupCreateTime = ginfo.GroupCreateTime
			ValGroupInfo.LastMsgSeq = ginfo.LastMsgSeq
			ValGroupInfo.Code = ginfo.Code
			ValGroupInfo.GroupLevel = ginfo.GroupLevel
			ValGroupInfo.OwnerUin = ginfo.OwnerUin
			ValGroupInfo.MaxMemberCount = ginfo.MaxMemberCount
			ValGroupInfo.OwnerUin = ginfo.OwnerUin
			ValGroupInfo.MemberCount = ginfo.MemberCount
		}
	})
}

func ParseMsgContent(elements []message.IMessageElement) *MsgContent {

	var content = &MsgContent{
		Texts:   []string{},
		At:      []int64{},
		Replies: []int32{},
		Faces:   []string{},
		Images:  []string{},
	}

	// find all texts and at targets
	for _, element := range elements {

		switch e := element.(type) {
		case *message.TextElement:
			content.Texts = append(content.Texts, e.Content)
		case *message.AtElement:
			content.At = append(content.At, e.Target)
		case *message.FaceElement:
			content.Faces = append(content.Faces, e.Name)
		case *message.GroupImageElement:
			content.Images = append(content.Images, e.Url)
		case *message.ReplyElement:
			content.Replies = append(content.Replies, e.ReplySeq)
		}
	}

	return content
}

func RefreshGroupMember() {
	ValGroupInfo.Update(func(info *client.GroupInfo) {
		if members, err := bot.Instance.GetGroupMembers(info); err != nil {
			logger.Warnf("更新群成员列表时出现错误: %v", err)
		} else {
			info.Members = members
		}
	})
}

func FindGroupMember(uid int64) *client.GroupMemberInfo {
	RefreshGroupMember()
	return ValGroupInfo.FindMember(uid)
}

func FindOtherGroupMember(members []*client.GroupMemberInfo, uid int64) *client.GroupMemberInfo {
	for _, member := range members {
		m := *member
		if m.Uin == uid {
			return member
		}
	}
	return nil
}

var GroupKey = func(groupCode int64, key string) string { return fmt.Sprintf("qq:group_%d:%s", groupCode, key) }
var toPrivateKey = func(uid int64, key string) string { return fmt.Sprintf("qq:private_%d:%s", uid, key) }

func GetRandomGroupMessage(gp int64) (*message.GroupMessage, error) {
	info, err := bot.Instance.GetGroupInfo(gp)
	if err != nil {
		return nil, err
	}
	return getRandomGroupMessageWithInfo(gp, info)
}

func getRandomGroupMessageWithInfo(gp int64, info *client.GroupInfo) (*message.GroupMessage, error) {
	rand.Seed(time.Now().UnixMicro())
	// MsgSeqAfter ~ LastMsgSeq 範圍內的隨機訊息ID
	id := rand.Int63n(info.LastMsgSeq-file.DataStorage.Setting.MsgSeqAfter) + file.DataStorage.Setting.MsgSeqAfter
	if botSaid.Contains(id) {
		// 略過機器人訊息
		return getRandomGroupMessageWithInfo(gp, info)
	}
	msg, err := GetGroupMessage(gp, id)
	if err != nil {
		return nil, err
	} else if msg.Sender.Uin == bot.Instance.Uin {
		// 不要機器人自己發過的訊息
		logger.Infof("獲取的隨機群訊息為機器人訊息，正在重新獲取...")
		botSaid.Add(msg.Id)
		<-time.After(time.Second) // 緩衝
		return GetRandomGroupMessage(gp)
	}
	return msg, nil
}

func GetGroupMessage(groupCode int64, seq int64) (*message.GroupMessage, error) {

	key := GroupKey(groupCode, fmt.Sprintf("msg:%d", seq))

	persistGroupMsg := &PersistentGroupMessage{}
	exist, err := redis.GetProto(key, persistGroupMsg)
	if err != nil {
		logger.Errorf("嘗試從 redis 獲取群組消息時出現錯誤: %v, 將使用 API 獲取", err)
	} else if exist {
		return persistGroupMsg.ToGroupMessage(), nil
	}

	msgList, err := bot.Instance.GetGroupMessages(groupCode, seq, seq+1)

	if err != nil {
		logger.Warnf("尝试获取群 %d 的群消息 (%d) 时出现错误: %v", groupCode, seq, err)

		// get msg error: 104 <= 消息不存在
		// 即 機器人加群前的消息，需要略過
		if strings.Contains(err.Error(), "104") {
			file.UpdateStorage(func() {
				if file.DataStorage.Setting.MsgSeqAfter < seq {
					file.DataStorage.Setting.MsgSeqAfter = seq
					logger.Warnf("已調整機器人消息獲取最低範圍為 %v", seq)
				}
			})
		}
		return nil, err
	}
	if len(msgList) > 0 {
		msg := msgList[0]
		// 非 bot 訊息才儲存
		if msg.Sender.Uin != bot.Instance.Uin {
			persistGroupMsg.Parse(msg)
			err = redis.StoreProto(key, persistGroupMsg)
			if err != nil {
				logger.Warnf("Redis 儲存群組消息時出現錯誤: %v", err)
			}
		}
		//修復圖片
		fixGroupImages(groupCode, msg)
		return msg, nil
	} else {
		return nil, nil
	}
}

func IsMuted(uid int64) bool {
	member := FindGroupMember(uid)
	if member == nil {
		return false
	}
	return member.ShutUpTimestamp > time.Now().Unix()
}

// reLogin - 參考了 Sora233/Mirai-Template 中的重連方式
func reLogin(qBot *bot.Bot) error {
	if qBot.Online.Load() {
		return nil
	}
	logger.Info("嘗試使用緩存會話登錄...")
	token, err := os.ReadFile("./session.token")
	if err == nil {
		err = qBot.TokenLogin(token)
		if err == nil {
			return nil
		}
	}
	logger.Warnf("緩存會話登錄失敗: %v", err)
	logger.Info("將嘗試使用普通登錄。")
	err = bot.CommonLogin()
	if err == nil {
		bot.SaveToken()
	}
	return err
}
