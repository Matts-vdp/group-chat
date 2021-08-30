package mess

import "encoding/json"

type Message struct {
	Sender string
	Data   string
}

func (mess Message) String() string {
	return mess.Sender + ": " + mess.Data
}

func (mess Message) ToJson() (string, error) {
	str, err := json.Marshal(mess)
	return string(str), err
}
func FromJson(js []byte) (Message, error) {
	var mess Message
	err := json.Unmarshal(js, &mess)
	return mess, err
}
