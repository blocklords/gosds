package reader

import (
	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/remote"
	"github.com/blocklords/gosds/topic"
)

type Reader struct {
	socket  *remote.Socket // SDS Gateway
	address string         // Account address granted for reading
}

func NewReader(gatewaySocket *remote.Socket, address string) *Reader {
	return &Reader{socket: gatewaySocket, address: address}
}

func (r *Reader) Read(t topic.Topic, args map[string]interface{}) message.Reply {
	if t.Level() != topic.FULL_LEVEL {
		return message.Fail(`Topic should contain method name`)
	}

	request := message.Request{
		Command: "smartcontract_read",
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
