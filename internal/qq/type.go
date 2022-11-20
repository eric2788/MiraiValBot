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
	Elements   [][]byte
}

func (m *PersistentGroupMessage) Parse(gp *message.GroupMessage) error {
	m.Id = gp.Id
	m.InternalId = gp.InternalId
	m.GroupCode = gp.GroupCode
	m.GroupName = gp.GroupName
	m.Sender = gp.Sender
	m.Time = gp.Time
	protoElements := message.ToProtoElems(gp.Elements, true)
	m.Elements = make([][]byte, len(protoElements))
	for i, ele := range protoElements {
		b, err := proto.Marshal(ele)
		if err != nil {
			return err
		}
		m.Elements[i] = b
	}
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
	elements := make([]*msg.Elem, len(m.Elements))
	for i, b := range m.Elements {
		var elem msg.Elem
		err := proto.Unmarshal(b, &elem)
		if err != nil {
			return nil, err
		}
		elements[i] = &elem
	}
	gp.Elements = message.ParseMessageElems(elements)
	return gp, nil
}
