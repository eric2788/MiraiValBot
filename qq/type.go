package qq

import (
	"github.com/Mrs4s/MiraiGo/client/pb/msg"
	"github.com/Mrs4s/MiraiGo/message"
)

type PersistentGroupMessage struct {
	Id             int32
	InternalId     int32
	GroupCode      int64
	GroupName      string
	Sender         *message.Sender
	Time           int32
	Elements       []*msg.Elem
	OriginalObject *msg.Message
	// OriginalElements []*msg.Elem
}

func (m *PersistentGroupMessage) Parse(gp *message.GroupMessage) {
	m.Id = gp.Id
	m.InternalId = gp.InternalId
	m.GroupCode = gp.GroupCode
	m.GroupName = gp.GroupName
	m.Sender = gp.Sender
	m.Time = gp.Time
	m.OriginalObject = gp.OriginalObject
	m.Elements = message.ToProtoElems(gp.Elements, true)
}

func (m *PersistentGroupMessage) ToGroupMessage() *message.GroupMessage {
	gp := &message.GroupMessage{}
	gp.Id = m.Id
	gp.InternalId = m.InternalId
	gp.GroupCode = m.GroupCode
	gp.GroupName = m.GroupName
	gp.Sender = m.Sender
	gp.Time = m.Time
	gp.Elements = message.ParseMessageElems(m.Elements)
	return gp
}
