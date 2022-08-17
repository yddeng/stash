package protocol

import (
	"initialthree/codec/ss"
)

func NewEncoder() *ss.Encoder {
	return ss.NewEncoder("center_msg", "center_req", "center_resp")
}

func NewReceiver() *ss.Receiver {
	return ss.NewReceiver("center_msg", "center_req", "center_resp")
}
