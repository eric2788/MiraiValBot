package qq

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/eric2788/common-utils/set"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/internal/redis"
)

var ValGroupInfo = &client.GroupInfo{
	Uin:     file.ApplicationYaml.Val.GroupId,
	Members: []*client.GroupMemberInfo{},
}

type MsgContent struct {
	Texts     []string
	At        []int64
	AtDisplay []string
	Faces     []string
	Images    []string
	Replies   []int32
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

// GetGroupEssenceMsgIds with cache
func GetGroupEssenceMsgIds() ([]int64, error) {
	gpDist, err := bot.Instance.GetGroupEssenceMsgList(ValGroupInfo.Code)
	essencesCache := GetEssenceList()

	if err != nil {
		return essencesCache, err
	}

	var messages = set.NewInt64()

	for _, dist := range gpDist {
		messages.Add(int64(dist.MessageID))
	}

	for _, id := range essencesCache {
		messages.Add(id)
	}

	return messages.ToArr(), nil

}

func ParseMsgContent(elements []message.IMessageElement) *MsgContent {

	var content = &MsgContent{
		Texts:     []string{},
		At:        []int64{},
		Replies:   []int32{},
		Faces:     []string{},
		Images:    []string{},
		AtDisplay: []string{},
	}

	// find all texts and at targets
	for _, element := range elements {

		switch e := element.(type) {
		case *message.TextElement:
			content.Texts = append(content.Texts, e.Content)
		case *message.AtElement:
			content.At = append(content.At, e.Target)
			content.AtDisplay = append(content.AtDisplay, e.Display)
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
	return getRandomGroupMessageWithInfo(info)
}

func GetRandomGroupMessageMember(gp, uid int64) (*message.GroupMessage, error) {
	info, err := bot.Instance.GetGroupInfo(gp)
	if err != nil {
		return nil, err
	}
	return getRandomGroupMessageWithMember(info, uid, 10, 0)
}

func getRandomGroupMessageWithMember(info *client.GroupInfo, uid, plus int64, times int) (*message.GroupMessage, error) {

	if times >= 10 {
		return nil, fmt.Errorf("已搜索十次依然找不到该群成员的群信息，请稍后再尝试。")
	}

	gp := info.Code
	rand.Seed(time.Now().UnixNano())
	// MsgSeqAfter ~ LastMsgSeq 範圍內的隨機訊息ID
	id := rand.Int63n(info.LastMsgSeq-file.DataStorage.Setting.MsgSeqAfter) + file.DataStorage.Setting.MsgSeqAfter - plus
	if botSaid.Contains(id) {
		// 略過機器人訊息
		return getRandomGroupMessageWithMember(info, uid, plus, times)
	}
	msgs, err := GetGroupMessages(gp, id, plus, false)
	if err != nil {
		// 不知是什麽，總之重新獲取
		if strings.Contains(err.Error(), "108") {
			logger.Errorf("嘗試獲取隨機消息時出現錯誤: %v, 將重新獲取...", err)
			<-time.After(time.Second) // 緩衝
			return getRandomGroupMessageWithMember(info, uid, plus, times)
		}
		return nil, err
	}

	for _, msg := range msgs {
		if msg.Sender.Uin == bot.Instance.Uin {
			// 不要機器人自己發過的訊息
			logger.Infof("獲取的隨機群訊息為機器人訊息，已略过")
			botSaid.Add(msg.Id)
		} else if msg.Sender.Uin == uid {
			FixGroupImages(msg.GroupCode, msg)
			return msg, nil
		}
	}

	logger.Warnf("找不到 %d 所发送的消息，正在重新获取... (%d次)", uid, times+1)
	<-time.After(time.Second) // 緩衝
	return getRandomGroupMessageWithMember(info, uid, plus, times+1)
}

func getRandomGroupMessageWithInfo(info *client.GroupInfo) (*message.GroupMessage, error) {
	gp := info.Code
	rand.Seed(time.Now().UnixNano())
	// MsgSeqAfter ~ LastMsgSeq 範圍內的隨機訊息ID
	id := rand.Int63n(info.LastMsgSeq-file.DataStorage.Setting.MsgSeqAfter) + file.DataStorage.Setting.MsgSeqAfter
	if botSaid.Contains(id) {
		// 略過機器人訊息
		return getRandomGroupMessageWithInfo(info)
	}
	msg, err := GetGroupMessage(gp, id)
	if err != nil {
		// 不知是什麽，總之重新獲取
		if strings.Contains(err.Error(), "108") {
			logger.Errorf("嘗試獲取隨機消息時出現錯誤: %v, 將重新獲取...", err)
			<-time.After(time.Second) // 緩衝
			return GetRandomGroupMessage(gp)
		}
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

func GetGroupMessages(groupCode int64, seq, plus int64, fixImg bool) (map[int64]*message.GroupMessage, error) {

	results := make(map[int64]*message.GroupMessage)

	for i := seq; i < seq+plus; i++ {
		key := GroupKey(groupCode, fmt.Sprintf("msg:%d", i))
		persistGroupMsg := &PersistentGroupMessage{}
		exist, err := redis.Get(key, persistGroupMsg)
		if err != nil {
			logger.Errorf("嘗試從 redis 獲取群組消息 %d 時出現錯誤: %v, 將使用 API 獲取", i, err)
		} else if exist {
			if msg, err := persistGroupMsg.ToGroupMessage(); err == nil {
				if fixImg {
					FixGroupImages(groupCode, msg)
				}
				results[i] = msg
			} else {
				logger.Errorf("嘗試從 redis 解析 群組消息 %d 時出現錯誤: %v, 將使用 API 獲取", i, err)
			}
		}
	}

	if len(results) >= int(plus) {
		return results, nil
	}

	msgList, err := bot.Instance.GetGroupMessages(groupCode, seq, seq+plus)
	if err != nil {
		logger.Warnf("尝试获取群 %d 的群消息 (%d ~ %d) 时出现错误: %v", groupCode, seq, seq+plus, err)

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

	for _, msg := range msgList {
		key := GroupKey(groupCode, fmt.Sprintf("msg:%d", msg.Id))
		persistGroupMsg := &PersistentGroupMessage{}
		err = persistGroupMsg.Parse(msg)
		if err != nil {
			logger.Warnf("嘗試序列化群組消息時出現錯誤: %v", err)
		} else {
			err = redis.Store(key, persistGroupMsg)
			if err != nil {
				logger.Warnf("Redis 儲存群組消息時出現錯誤: %v", err)
			}
		}
		if fixImg {
			FixGroupImages(groupCode, msg)
		}
		results[int64(msg.Id)] = msg
	}

	return results, nil
}

func GetGroupMessage(groupCode int64, seq int64) (*message.GroupMessage, error) {

	key := GroupKey(groupCode, fmt.Sprintf("msg:%d", seq))

	persistGroupMsg := &PersistentGroupMessage{}
	exist, err := redis.Get(key, persistGroupMsg)
	if err != nil {
		logger.Errorf("嘗試從 redis 獲取群組消息時出現錯誤: %v, 將使用 API 獲取", err)
	} else if exist {
		if msg, err := persistGroupMsg.ToGroupMessage(); err == nil {
			//修復圖片
			FixGroupImages(groupCode, msg)
			return msg, nil
		} else {
			logger.Errorf("嘗試從 redis 解析 群組消息 時出現錯誤: %v, 將使用 API 獲取", err)
		}
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
			err = persistGroupMsg.Parse(msg)
			if err != nil {
				logger.Warnf("嘗試序列化群組消息時出現錯誤: %v", err)
			} else {
				err = redis.Store(key, persistGroupMsg)
				if err != nil {
					logger.Warnf("Redis 儲存群組消息時出現錯誤: %v", err)
				}
			}

		}
		//修復圖片
		FixGroupImages(groupCode, msg)
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
