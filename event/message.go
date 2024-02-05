package event

type Message struct {
	key     string
	source  string
	content string
}

func NewMessage(key, source, content string) *Message {
	return &Message{
		key:     key,
		source:  source,
		content: content,
	}
}

func (m *Message) GetKey() string {
	return m.key
}

func (m *Message) GetSource() string {
	return m.source
}

func (m *Message) GetContent() string {
	return m.content
}

func (m *Message) SetKey(key string) {
	m.key = key
}

func (m *Message) SetSource(source string) {
	m.source = source
}

func (m *Message) SetContent(content string) {
	m.content = content
}