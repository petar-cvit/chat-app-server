package models

type Messages struct {
	Messages []*Message `json:"messages"`
}

type Message struct {
	Text   string `json:"text"`
	Time   string `json:"time"`
	Issuer string `json:"issuer"`
}

func (m *Messages) AppendMessage(message *Message) {
	m.Messages = append(m.Messages, message)
}
