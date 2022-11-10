package cmd

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/qq"
	"github.com/eric2788/MiraiValBot/valorant"
	"github.com/eric2788/common-utils/datetime"
)

func randomMember(args []string, source *command.MessageSource) error {
	rand.Seed(time.Now().UnixMicro())
	members := qq.ValGroupInfo.Members

	if len(members) == 0 {
		reply := qq.CreateReply(source.Message).Append(message.NewText("群成员列表为空。"))
		_ = qq.SendGroupMessage(reply)
		return nil
	}

	chosen := members[rand.Intn(len(members))]
	at := message.NewAt(chosen.Uin)
	at.Display = "@" + chosen.Nickname
	reply := message.NewSendingMessage().Append(at)
	return qq.SendGroupMessage(reply)
}

func randomLong(args []string, source *command.MessageSource) error {
	msg := qq.CreateReply(source.Message)
	backup := "https://phqghume.github.io/img/"
	rand.Seed(time.Now().UnixMicro())
	random := rand.Intn(58) + 1
	ext := ".jpg"
	if random > 48 {
		ext = ".gif"
	}
	imgLink := fmt.Sprintf("%slong%%20(%d)%s", backup, random, ext)
	img, err := qq.NewImageByUrl(imgLink)
	if err != nil {
		return err
	}
	msg.Append(img)
	return qq.SendGroupMessage(msg)
}

func randomChoice(args []string, source *command.MessageSource) error {

	msg := qq.CreateReply(source.Message)

	if len(args) == 0 {
		msg.Append(message.NewText("选项不得为空。"))
		return qq.SendGroupMessage(msg)
	}

	rand.Seed(time.Now().UnixMicro())
	chosen := args[rand.Intn(len(args))]

	msg.Append(qq.NewTextf(chosen))

	return qq.SendGroupMessage(msg)
}

func randomMessage(args []string, source *command.MessageSource) error {

	msg, err := qq.GetRandomGroupMessage(source.Message.GroupCode)
	if err != nil {
		return err
	} else if msg == nil || len(msg.Elements) == 0 {
		return fmt.Errorf("随机消息为空")
	}

	reply := message.NewSendingMessage()
	var nick string
	if msg.Sender.CardName == "" {
		nick = msg.Sender.Nickname
	} else {
		nick = msg.Sender.CardName
	}
	reply.Append(qq.NewTextfLn("%s 在 %s 说过: ", nick, datetime.FormatSeconds(int64(msg.Time))))
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

	return qq.SendGroupMessage(reply)
}

func randomPhoto(args []string, source *command.MessageSource) error {
	rand.Seed(time.Now().UnixMicro())
	imgs := qq.GetImageList()

	if len(imgs) == 0 {
		return qq.SendGroupMessage(qq.CreateReply(source.Message).Append(qq.NewTextf("群图片缓存列表为空。")))
	}

	logger.Debugf("成功索取 %d 张群图片缓存。", len(imgs))

	chosen := imgs[rand.Intn(len(imgs))]
	img, err := qq.NewImageByByte(chosen)
	if err != nil {
		return err
	}
	return qq.SendGroupMessage(message.NewSendingMessage().Append(img))
}

func randomEssence(args []string, source *command.MessageSource) error {

	rand.Seed(time.Now().UnixMicro())

	msgids, err := qq.GetGroupEssenceMsgIds()
	if err != nil {
		logger.Warnf("获取群精华消息列表失败: %v, 將使用純緩存列表", source.Message.GroupCode)

		// 快取是 0
		if len(msgids) == 0 {
			return err
		}

	}

	if len(msgids) == 0 {
		reply := qq.CreateReply(source.Message).Append(message.NewText("群精华消息列表为空。"))
		_ = qq.SendGroupMessage(reply)
		return nil
	}

	chosen := msgids[rand.Intn(len(msgids))]
	essence, err := qq.GetGroupEssenceMessage(chosen)

	if err != nil {
		logger.Warnf("获取群精华消息失败: %d", chosen)
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

	return qq.SendGroupMessage(msg)
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

func randomBundle(args []string, source *command.MessageSource) error {

	rand.Seed(time.Now().UnixMicro())

	bundles, err := valorant.GetBundles(valorant.SC)
	if err != nil {
		return err
	}

	chosen := bundles[rand.Intn(len(bundles))]

	msg := qq.CreateReply(source.Message)

	msg.Append(qq.NewTextfLn("选中套装: %s", chosen.DisplayName))
	if chosen.ExtraDescription != "" {
		msg.Append(qq.NewTextfLn("简介: %s", chosen.ExtraDescription))
	}
	if chosen.PromoDescription != "" {
		msg.Append(qq.NewTextfLn("推广: %s", chosen.PromoDescription))
	}

	img, err := qq.NewImageByUrl(chosen.DisplayIcon)
	if err != nil {
		img, err = qq.NewImageByUrl(chosen.DisplayIcon2)
	}

	if err != nil {
		logger.Errorf("索取套装图片时出现错误: %v", err)
		msg.Append(qq.NewTextf("[图片]"))
	} else {
		msg.Append(img)
	}

	return qq.SendWithRandomRiskyStrategy(msg)
}

func randomSkin(args []string, source *command.MessageSource) error {

	name := strings.ToLower(args[0])

	weapons, err := valorant.GetWeapons(valorant.AllWeapons, valorant.SC)
	if err != nil {
		return err
	}

	weapon := valorant.GetWeapon(weapons, name)
	msg := qq.CreateReply(source.Message)

	if weapon == nil {
		msg.Append(qq.NewTextf("没有此武器: %s, 请使用简中武器名称。", name))
		return qq.SendGroupMessage(msg)
	}

	rand.Seed(time.Now().UnixMicro())

	skin := weapon.Skins[rand.Intn(len(weapon.Skins))]

	if len(skin.Chromas) == 0 {
		msg.Append(qq.NewTextfLn("选中皮肤: %s", skin.DisplayName))
		img, err := qq.NewImageByUrl(skin.DisplayIcon)
		if err != nil {
			logger.Errorf("索取皮肤图片时错误: %v", err)
			msg.Append(qq.NewTextf("[图片]"))
		} else {
			msg.Append(img)
		}
	} else {

		rand.Seed(time.Now().UnixMicro())

		chroma := skin.Chromas[rand.Intn(len(skin.Chromas))]

		msg.Append(qq.NewTextfLn("选中皮肤: %s", chroma.DisplayName))

		icon := chroma.FullRender
		if icon == "" {
			icon = chroma.DisplayIcon
		}

		if icon != "" {
			img, err := qq.NewImageByUrl(icon)
			if err != nil {
				logger.Errorf("索取皮肤图片时错误: %v", err)
				msg.Append(qq.NewTextf("[图片]"))
			} else {
				msg.Append(img)
			}
		}

		if chroma.StreamedVideo != "" {

			if err := qq.SendGroupMessage(msg); err != nil {
				return err
			}

			msg := message.NewSendingMessage()
			if video, err := qq.NewVideoByUrl(chroma.StreamedVideo, icon); err != nil {
				logger.Error(err)
			} else {
				msg.Append(video)
				return qq.SendGroupMessage(msg)
			}

		}
	}
	return qq.SendWithRandomRiskyStrategy(msg)
}

var (
	randomEssenceCommand = command.NewNode([]string{"essence", "群精华"}, "获取随机一条群精华消息", false, randomEssence)
	randomMemberCommand  = command.NewNode([]string{"member", "成员"}, "随机群成员指令", false, randomMember)
	randomMessageCommand = command.NewNode([]string{"message", "msg", "群消息"}, "随机群消息指令", false, randomMessage)
	randomChoiceCommand  = command.NewNode([]string{"choice", "选项"}, "随机选项指令", false, randomChoice)
	randomAgentCommand   = command.NewNode([]string{"agent", "特务", "角色"}, "随机抽选一个瓦角色", false, randomAgent, "[角色类型]")
	randomWeaponCommand  = command.NewNode([]string{"weapon", "武器"}, "随机抽选一个瓦武器", false, randomWeapon, "[武器类型]")
	randomBundleCommand  = command.NewNode([]string{"bundle", "套装"}, "随机抽选一个瓦套装", false, randomBundle)
	randomSkinCommand    = command.NewNode([]string{"skin", "皮肤"}, "随机抽选一个瓦皮肤", false, randomSkin, "<武器名称>")
	randomDragonCommand  = command.NewNode([]string{"long", "dragon", "龙图"}, "随机抽选一张龙图", false, randomLong)
	randomPhotoCommand   = command.NewNode([]string{"photo", "image", "图片"}, "随机抽选一张群图片", false, randomPhoto)
)

var randomCommand = command.NewParent([]string{"random", "随机"}, "随机指令",
	randomMemberCommand,
	randomEssenceCommand,
	randomChoiceCommand,
	randomMessageCommand,
	randomAgentCommand,
	randomWeaponCommand,
	randomBundleCommand,
	randomSkinCommand,
	randomDragonCommand,
	randomPhotoCommand,
)

func init() {
	command.AddCommand(randomCommand)
}
