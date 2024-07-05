package pow

import (
	"encoding/json"
	"fmt"
)

const (
	MsgTx    = "TX"
	MsgBlock = "BLOCK"
)

// Msg encapsulates the Data for http communication usages
type Msg[T any] struct {
	Data T      `json:"data"`
	Type string `json:"type"`
}

func WrapData2MsgStr(data interface{}) ([]byte, error) {
	var jsonData []byte
	var err error

	switch v := data.(type) {
	case Transaction:
		msg := Msg[Transaction]{Data: v, Type: MsgTx}
		jsonData, err = json.Marshal(msg)
	case Block:
		msg := Msg[Block]{Data: v, Type: MsgBlock}
		jsonData, err = json.Marshal(msg)
	default:
		return nil, fmt.Errorf("unsupported type")
	}

	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
