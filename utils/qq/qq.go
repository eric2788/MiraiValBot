package qq

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/redis"
	"strings"
)

var ValGroupInfo = &client.GroupInfo{
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
	ginfo.Members = members

	ValGroupInfo = ginfo

	logger.Infof("以指定 %s (%d) 为 瓦群。(共 %d 個成員)", ValGroupInfo.Name, ValGroupInfo.Uin, len(ValGroupInfo.Members))
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
			content.Texts = append(content.Texts, e.Display)
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

func FindGroupMember(uid int64) *client.GroupMemberInfo {
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

var key = func(g int64, m int64) string { return fmt.Sprintf("qq:group_%d:msg:%d", g, m) }

func GetGroupMessage(groupCode int64, seq int64) (*message.GroupMessage, error) {
	persistGroupMsg := &PersistentGroupMessage{}
	exist, err := redis.Get(key(groupCode, seq), persistGroupMsg)
	if err != nil {
		return nil, err
	} else if exist {
		return persistGroupMsg.ToGroupMessage(), nil
	}

	msgList, err := bot.Instance.GetGroupMessages(groupCode, seq, seq+1)

	if err != nil {
		return nil, err
	}
	if len(msgList) > 0 {
		msg := msgList[0]
		persistGroupMsg.Parse(msg)
		err = redis.Store(key(groupCode, seq), persistGroupMsg)
		if err != nil {
			logger.Warnf("Redis 儲存群組消息時出現錯誤: %v", err)
		}
		return msg, nil
	} else {
		return nil, nil
	}
}
