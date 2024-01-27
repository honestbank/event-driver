package event

type Message interface {
	GetKey() string
	GetSource() string
	GetContent() string
	SetKey(string)
	SetSource(string)
	SetContent(string)
}
