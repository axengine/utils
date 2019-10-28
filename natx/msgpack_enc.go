package natx

import (
	msgpack "github.com/vmihailenco/msgpack/v4"
)

const (
	MSGP_ENCODER = "msgp"
)

type MsgpEncoder struct {
	// Empty
}

// Encode
func (ge *MsgpEncoder) Encode(subject string, v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

// Decode
func (ge *MsgpEncoder) Decode(subject string, data []byte, vPtr interface{}) (err error) {
	return msgpack.Unmarshal(data, vPtr)
}
