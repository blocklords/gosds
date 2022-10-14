package message

import (
	"encoding/json"
)

type Reply struct {
	Status  string
	Message string
	Params  map[string]interface{}
}

func Fail(err string) Reply {
	return Reply{Status: "fail", Message: err}
}

func (r *Reply) IsOK() bool { return r.Status == "ok" || r.Status == "OK" }

/* Convert to format understood by the protocol */
func (reply *Reply) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"status":   reply.Status,
		"message:": reply.Message,
		"params":   reply.Params,
	}
}

func (reply *Reply) ToString() string {
	return string(reply.ToBytes())
}

func (reply *Reply) ToBytes() []byte {
	interfaces := reply.ToJSON()
	byt, err := json.Marshal(interfaces)
	if err != nil {
		return []byte{}
	}

	return byt
}

func ParseReply(msgs []string) (Reply, error) {
	msg := ""
	for _, v := range msgs {
		msg += v
	}
	var dat map[string]interface{}

	if err := json.Unmarshal([]byte(msg), &dat); err != nil {
		return Reply{}, err
	}

	replyMessage := ""
	if dat["message"] != nil {
		replyMessage = dat["message"].(string)
	}

	return Reply{Status: dat["status"].(string), Params: dat["params"].(map[string]interface{}), Message: replyMessage}, nil
}
