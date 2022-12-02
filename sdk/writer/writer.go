package writer

import (
	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/remote"
	"github.com/blocklords/gosds/topic"
)

type Writer struct {
	socket  *remote.Socket // SDS Gateway host
	address string         // Account address granted for reading
}

func NewWriter(gatewaySocket *remote.Socket, address string) *Writer {
	return &Writer{socket: gatewaySocket, address: address}
}

func (r *Writer) Write(t topic.Topic, args map[string]interface{}) message.Reply {
	if t.Level() != topic.FULL_LEVEL {
		return message.Fail(`Topic should contain method name`)
	}

	request := message.Request{
		Command: "smartcontract_write",
		Param: map[string]interface{}{
			"topic_string": t.ToString(topic.FULL_LEVEL),
			"arguments":    args,
			"address":      r.address,
		},
	}

	params, err := r.socket.RequestRemoteService(&request)
	if err != nil {
		return message.Fail(err.Error())
	}

	return message.Reply{Status: "OK", Message: "", Params: params}
}

func (r *Writer) AddToPool(t topic.Topic, args map[string]interface{}) message.Reply {
	if t.Level() != topic.FULL_LEVEL {
		return message.Fail(`Topic should contain method name`)
	}

	request := message.Request{
		Command: "pool_add",
		Param: map[string]interface{}{
			"topic_string": t.ToString(topic.FULL_LEVEL),
			"arguments":    args,
			"address":      r.address,
		},
	}

	params, err := r.socket.RequestRemoteService(&request)
	if err != nil {
		return message.Fail(err.Error())
	}

	return message.Reply{Status: "OK", Message: "", Params: params}
}
