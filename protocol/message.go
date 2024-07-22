package protocol

const (
	MessageTypeOpen = iota
	MessageTypeClose
	MessageTypePing
	MessageTypePong
	MessageTypeEmpty
	MessageTypeEmit
	MessageTypeAckRequest
	MessageTypeAckResponse
)

type Message struct {
	Type   int
	AckId  int
	Method string
	Args   string
	Source string
}
