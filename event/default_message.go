package event

type DefaultMessage struct {
	key     string
	source  string
	content string
}

func NewMessage(key, source, content string) *DefaultMessage {
	return &DefaultMessage{
		key:     key,
		source:  source,
		content: content,
	}
}

func (m *DefaultMessage) GetKey() string {
	return m.key
}

func (m *DefaultMessage) GetSource() string {
	return m.source
}

func (m *DefaultMessage) GetContent() string {
	return m.content
}
