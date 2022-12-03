package server

import "encoding/json"

type Event struct {
	Type        string            `json:"type"`
	Information map[string]string `json:"info"`
}

func (event *Event) ToJson() []byte {
	jsonData, _ := json.Marshal(event)
	return jsonData
}
