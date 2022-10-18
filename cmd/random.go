package cmd

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/qq"
	qq2 "github.com/eric2788/MiraiValBot/qq"
	"github.com/eric2788/MiraiValBot/valorant"
	"github.com/eric2788/common-utils/datetime"
)

func randomMember(args []string, source *command.MessageSource) error {
	rand.Seed(time.Now().UnixMicro())
	members := qq2.ValGroupInfo.Members

	if len(members) == 0 {
		reply := qq2.CreateReply(source.Message).Append(message.NewText("群成员列表为空。"))
		_ = qq2.SendGroupMessage(reply)
		return nil
	}

	chosen := members[rand.Intn(len(members))]
	reply := message.NewSendingMessage().Append(message.NewAt(chosen.Uin))
	return qq2.SendGroupMessage(reply)
}

func randomMessage(args []string, source *command.MessageSource) error {

	msg, err := qq2.GetRandomGroupMessage(source.Message.GroupCode)
	if err != nil {
		return err
	} else if msg == nil {
		return fmt.Errorf("随机消息为空")
	}

	reply := message.NewSendingMessage()
	var nick string
	if msg.Sender.CardName == "" {
		nick = msg.Sender.Nickname
	} else {
		nick = msg.Sender.CardName
	}
	reply.Append(qq2.NewTextfLn("%s 在 %s 说过: ", nick, datetime.FormatSeconds(int64(msg.Time))))
	for _, element := range msg.Elements {
		switch element.(type) {
		case *message.ReplyElement:
			continue
		case *message.ForwardElement:
			continue
		default:
			break
		}
		reply.Append(element)
	}

	return qq2.SendGroupMessage(reply)
}

func randomEssence(args []string, source *command.MessageSource) error {

	rand.Seed(time.Now().UnixMicro())

	gpDist, err := source.Client.GetGroupEssenceMsgList(source.Message.GroupCode)

	// why empty ? not sure but let's try other method
	if len(gpDist) == 0 {
		logger.Warnf("群消息為空，正在使用第 2 種方式獲取")
		gpDist, err = source.Client.GetGroupEssenceMsgList(qq2.ValGroupInfo.Uin)
	}

	// why empty ? not sure but let's try other method
	if len(gpDist) == 0 {
		logger.Warnf("群消息為空，正在使用第 3 種方式獲取")
		gpDist, err = bot.Instance.GetGroupEssenceMsgList(source.Message.GroupCode)
	}

	// why empty ? not sure but let's try other method
	if len(gpDist) == 0 {
		logger.Warnf("群消息為空，正在使用第 4 種方式獲取")
		gpDist, err = bot.Instance.GetGroupEssenceMsgList(qq2.ValGroupInfo.Uin)
	}

	if err != nil {
		logger.Warnf("获取群精华消息列表失败: %v", source.Message.GroupCode)
		return err
	}

	if len(gpDist) == 0 {
		reply := qq2.CreateReply(source.Message).Append(message.NewText("群精华消息列表为空。"))
		_ = qq2.SendGroupMessage(reply)
		return nil
	}

	chosen := gpDist[rand.Intn(len(gpDist))]

	seq := int64(chosen.MessageID)
	essence, err := qq2.GetGroupMessage(source.Message.GroupCode, seq)

	if err != nil {
		logger.Warnf("获取群精华消息失败: %+v", chosen)
		return err
	}
	msg := message.NewSendingMessage()

	if essence != nil {
		for _, element := range essence.Elements {
			msg.Append(element)
		}
	} else {
		msg.Append(message.NewText("没有群精华消息"))
	}

	return qq2.SendGroupMessage(msg)
}

func randomAgent(args []string, source *command.MessageSource) error {

	var agentType valorant.AgentType
	if len(args) == 0 {
		agentType = valorant.AllAgents
	} else {
		switch t := strings.ToLower(args[0]); t {
		case "决斗者", "決鬥者", "duelist":
			agentType = valorant.Duelist
		case "守衛", "守卫", "guard":
			agentType = valorant.Guard
		case "控场", "控場", "controller":
			agentType = valorant.Controller
		case "先鋒", "先锋", "inititator":
			agentType = valorant.Initiator
		default:
			return errors.New("无效的 Agent 类型")
		}
	}

	agents, err := valorant.GetAgents(agentType, valorant.SC)
	if err != nil {
		return err
	}

	// impossible
	if len(agents) == 0 {
		return errors.New("角色列表为空")
	}

	rand.Seed(time.Now().UnixMilli())

	random := agents[rand.Intn(len(agents))]

	msg := qq.CreateReply(source.Message)
	msg.Append(qq.NewTextfLn("选中角色: %s", random.DisplayName))
	msg.Append(qq.NewTextfLn("类型: %s", random.Role.DisplayName))
	msg.Append(qq.NewTextfLn("简介: %s", random.Description))

	if random.CharacterTags != nil {
		msg.Append(qq.NewTextfLn("标签: %s", strings.Join(*random.CharacterTags, ", ")))
	}

	skills := make([]string, 0)

	for _, skill := range random.Abilities {
		skills = append(skills, skill.DisplayName)
	}

	msg.Append(qq.NewTextfLn("技能: %s", strings.Join(skills, ", ")))

	img, err := qq.NewImageByUrl(random.FullPortrait)
	if err != nil {
		logger.Errorf("尝试获取角色 %s 图片时出现错误: %v", random.DisplayName, err)
	} else {
		msg.Append(img)
	}

	return qq.SendWithRandomRiskyStrategy(msg)
}

func randomWeapon(args []string, source *command.MessageSource) error {
	var weaponType valorant.WeaponType

	if len(args) == 0 {
		weaponType = valorant.AllWeapons
	} else {
		switch t := strings.ToLower(args[0]); t {
		case "机枪", "重机枪", "heavy":
			weaponType = valorant.Heavy
		case "长枪", "步枪", "rifle":
			weaponType = valorant.Rifle
		case "狙击枪", "狙击", "sniper":
			weaponType = valorant.Sniper
		case "冲锋枪", "轻机枪", "lmg", "smg":
			weaponType = valorant.SMG
		case "手枪", "pistol", "sidearm":
			weaponType = valorant.Pistol
		default:
			return errors.New("未知的武器类型")
		}
	}

	weapons, err := valorant.GetWeapons(weaponType, valorant.SC)
	if err != nil {
		return err
	}

	// impossible
	if len(weapons) == 0 {
		return errors.New("武器列表为空")
	}

	rand.Seed(time.Now().UnixMilli())

	random := weapons[rand.Intn(len(weapons))]

	msg := qq.CreateReply(source.Message)
	msg.Append(qq.NewTextfLn("选中枪械: %s", random.DisplayName))
	msg.Append(qq.NewTextfLn("类型: %s", random.ShopData.CategoryText))
	msg.Append(qq.NewTextfLn("价格: $%d", random.ShopData.Cost))

	img, err := qq.NewImageByUrl(random.DisplayIcon)
	if err != nil {
		logger.Errorf("获取武器 %s 图片时出现错误: %v", random.DisplayName, err)
	} else {
		msg.Append(img)
	}

	return qq.SendWithRandomRiskyStrategy(msg)
}

var (
	randomEssenceCommand = command.NewNode([]string{"essence", "群精华"}, "获取随机一条群精华消息", false, randomEssence)
	randomMemberCommand  = command.NewNode([]string{"member", "成员"}, "随机群成员指令", false, randomMember)
	randomMessageCommand = command.NewNode([]string{"message", "msg", "群消息"}, "随机群消息指令", false, randomMessage)
	randomAgentCommand   = command.NewNode([]string{"agent", "特务", "角色"}, "随机抽选一个瓦角色", false, randomAgent, "[角色类型]")
	randomWeaponCommand  = command.NewNode([]string{"weapon", "武器"}, "随机抽选一个瓦武器", false, randomWeapon, "[武器类型]")
)

var randomCommand = command.NewParent([]string{"random", "随机"}, "随机指令",
	randomMemberCommand,
	randomEssenceCommand,
	randomMessageCommand,
	randomAgentCommand,
	randomWeaponCommand,
)

func init() {
	command.AddCommand(randomCommand)
}
