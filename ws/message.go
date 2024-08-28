package ws

type FrameType uint8

const (
	FrameData FrameType = 0x0
	FramePing FrameType = 0x1
)

type Message struct {
	FrameType  `json:"frameType"`
	Method     string      `json:"method"`
	FormUserId int         `json:"formUserId"`
	ToUserId   int         `json:"toUserId"`
	Data       interface{} `json:"data"`
}

func NewMessage(formUserId, toUserId int, data interface{}) *Message {
	return &Message{
		FrameType:  FrameData,
		FormUserId: formUserId,
		ToUserId:   toUserId,
		Data:       data,
	}
}
