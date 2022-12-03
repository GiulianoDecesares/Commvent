package server

import "encoding/json"

type Command struct {
	Type      string            `json:"type"`
	Arguments map[string]string `json:"args"`
}

func (command *Command) ToJson() []byte {
	jsonData, _ := json.Marshal(command)
	return jsonData
}
