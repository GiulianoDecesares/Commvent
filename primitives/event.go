package primitives

import "encoding/json"

type Event struct {
	Type        string      `json:"type"`
	Information Information `json:"info"`
}

func NewEvent(eventType string) *Event {
	return &Event{
		Type:        eventType,
		Information: *NewInformation(),
	}
}

func NewEventFromJSON(bytes []byte) *Event {
	event := Event{}
	json.Unmarshal(bytes, &event)

	return &event
}

func (event *Event) ToJson() []byte {
	jsonData, _ := json.Marshal(event)
	return jsonData
}
