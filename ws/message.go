package ws

type FrameType uint8

const (
	FrameData FrameType = 0x0
	FramePing FrameType = 0x1
)

type Message struct {
	FrameType `json:"frameType"`
	Method    string      `json:"method"`
	FormId    int         `json:"form_id"`
	Data      interface{} `json:"data"`
}

func NewMessage(formId int, data interface{}) *Message {
	return &Message{
		FrameType: FrameData,
		FormId:    formId,
		Data:      data,
	}
}
