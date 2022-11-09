package qq

import (
	"github.com/Mrs4s/MiraiGo/client/pb/msg"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/RomiChan/protobuf/proto"
)

type PersistentGroupMessage struct {
	Id         int32
	InternalId int32
	GroupCode  int64
	GroupName  string
	Sender     *message.Sender
	Time       int32
	Elements   []byte
}

func (m *PersistentGroupMessage) Parse(gp *message.GroupMessage) error {
	m.Id = gp.Id
	m.InternalId = gp.InternalId
	m.GroupCode = gp.GroupCode
	m.GroupName = gp.GroupName
	m.Sender = gp.Sender
	m.Time = gp.Time
	protoElements := message.ToProtoElems(gp.Elements, true)
	b, err := proto.Marshal(&protoElements)
	if err != nil {
		return err
	}
	m.Elements = b
	return nil
}

func (m *PersistentGroupMessage) ToGroupMessage() (*message.GroupMessage, error) {
	gp := &message.GroupMessage{}
	gp.Id = m.Id
	gp.InternalId = m.InternalId
	gp.GroupCode = m.GroupCode
	gp.GroupName = m.GroupName
	gp.Sender = m.Sender
	gp.Time = m.Time
	var elements []*msg.Elem
	err := proto.Unmarshal(m.Elements, &elements)
	if err != nil {
		return nil, err
	}
	gp.Elements = message.ParseMessageElems(elements)
	return gp, nil
}
