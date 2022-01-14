package qq

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

var testSerialized = `
{"Id":3742,"InternalId":-1423422452,"GroupCode":486518527,"GroupName":"","Sender":{"Uin":2899929243,"Nickname":"我什么也不知道","CardName":"","AnonymousInfo":null,"IsFriend":true},"Time":1641659134,"Elements":[{"Text":{"Str":"aw","Link":null,"Attr6Buf":null,"Attr7Buf":null,"Buf":null,"PbReserve":null},"Face":null,"OnlineImage":null,"NotOnlineImage":null,"TransElemInfo":null,"MarketFace":null,"CustomFace":null,"ElemFlags2":null,"RichMsg":null,"GroupFile":null,"ExtraInfo":null,"VideoFile":null,"AnonGroupMsg":null,"QQWalletMsg":null,"CustomElem":null,"GeneralFlags":null,"SrcMsg":null,"LightApp":null,"CommonElem":null}],"OriginalObject":{"Head":{"FromUin":2899929243,"ToUin":3585680664,"MsgType":82,"C2CCmd":1,"MsgSeq":3742,"MsgTime":1641659134,"MsgUid":null,"C2CTmpMsgHead":null,"GroupInfo":{"GroupCode":486518527,"GroupType":1,"GroupInfoSeq":484,"GroupCard":"我什么也不知道","GroupRank":null,"GroupLevel":1,"GroupCardType":2,"GroupName":null},"FromAppid":null,"FromInstid":null,"UserActive":null,"DiscussInfo":null,"FromNick":null,"AuthUin":null,"AuthNick":null,"MsgFlag":16,"AuthRemark":null,"GroupName":null,"MutiltransHead":null,"MsgInstCtrl":null,"PublicAccountGroupSendFlag":null,"WseqInC2CMsghead":null,"Cpid":null,"ExtGroupKeyInfo":null,"MultiCompatibleText":null,"AuthSex":null,"IsSrcMsg":null},"Content":{"PkgNum":1,"PkgIndex":0,"DivSeq":0,"AutoReply":null},"Body":{"RichText":{"Attr":{"CodePage":0,"Time":1641659134,"Random":-1423422452,"Color":0,"Size":9,"Effect":0,"CharSet":136,"PitchAndFamily":0,"FontName":"Microsoft YaHei","ReserveData":null},"Elems":[{"Text":{"Str":"aw","Link":null,"Attr6Buf":null,"Attr7Buf":null,"Buf":null,"PbReserve":null},"Face":null,"OnlineImage":null,"NotOnlineImage":null,"TransElemInfo":null,"MarketFace":null,"CustomFace":null,"ElemFlags2":null,"RichMsg":null,"GroupFile":null,"ExtraInfo":null,"VideoFile":null,"AnonGroupMsg":null,"QQWalletMsg":null,"CustomElem":null,"GeneralFlags":null,"SrcMsg":null,"LightApp":null,"CommonElem":null},{"Text":null,"Face":null,"OnlineImage":null,"NotOnlineImage":null,"TransElemInfo":null,"MarketFace":null,"CustomFace":null,"ElemFlags2":null,"RichMsg":null,"GroupFile":null,"ExtraInfo":null,"VideoFile":null,"AnonGroupMsg":null,"QQWalletMsg":null,"CustomElem":null,"GeneralFlags":{"BubbleDiyTextId":null,"GroupFlagNew":null,"Uin":null,"RpId":null,"PrpFold":null,"LongTextFlag":null,"LongTextResid":null,"GroupType":null,"ToUinFlag":null,"GlamourLevel":null,"MemberLevel":null,"GroupRankSeq":2,"OlympicTorch":null,"BabyqGuideMsgCookie":null,"Uin32ExpertFlag":null,"BubbleSubId":null,"PendantId":null,"RpIndex":null,"PbReserve":"CAZ4gIAEyAEA8AEA+AEAkAIAmAMAoAMAsAMAwAMA0AMAigQECAEQEbgEAMAEAMoEAPgEgIAIiAUA"},"SrcMsg":null,"LightApp":null,"CommonElem":null},{"Text":null,"Face":null,"OnlineImage":null,"NotOnlineImage":null,"TransElemInfo":null,"MarketFace":null,"CustomFace":null,"ElemFlags2":{"ColorTextId":0,"MsgId":null,"WhisperSessionId":null,"PttChangeBit":null,"VipStatus":null,"CompatibleId":null,"Insts":null,"MsgRptCnt":1,"SrcInst":null,"Longtitude":null,"Latitude":null,"CustomFont":null,"PcSupportDef":null,"CrmFlags":null},"RichMsg":null,"GroupFile":null,"ExtraInfo":null,"VideoFile":null,"AnonGroupMsg":null,"QQWalletMsg":null,"CustomElem":null,"GeneralFlags":null,"SrcMsg":null,"LightApp":null,"CommonElem":null},{"Text":null,"Face":null,"OnlineImage":null,"NotOnlineImage":null,"TransElemInfo":null,"MarketFace":null,"CustomFace":null,"ElemFlags2":null,"RichMsg":null,"GroupFile":null,"ExtraInfo":{"Nick":"5oiR5LuA5LmI5Lmf5LiN55+l6YGT","GroupCard":null,"Level":1,"Flags":16,"GroupMask":null,"MsgTailId":null,"SenderTitle":null,"ApnsTips":null,"Uin":null,"MsgStateFlag":null,"ApnsSoundType":null,"NewGroupFlag":null},"VideoFile":null,"AnonGroupMsg":null,"QQWalletMsg":null,"CustomElem":null,"GeneralFlags":null,"SrcMsg":null,"LightApp":null,"CommonElem":null}],"NotOnlineFile":null,"Ptt":null},"MsgContent":null,"MsgEncryptContent":null}}}
`

func TestParseFromPersistence(t *testing.T) {
	var persist = &PersistentGroupMessage{}
	if err := json.Unmarshal([]byte(testSerialized), persist); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", *persist)
	gpMsg := persist.ToGroupMessage()
	fmt.Printf("%+v", *gpMsg)
}

func TestDuration(t *testing.T) {

	var sec time.Duration = 10

	fmt.Printf("%d 秒", sec/time.Second)
}
