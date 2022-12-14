package primitives

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Message struct {
	Type        string            `json:"type"`
	SenderID    string            `json:"senderId"`
	Information map[string]string `json:"info"`
}

func NewMessage(messageType string, senderId string) *Message {
	return &Message{
		Type:        messageType,
		SenderID:    senderId,
		Information: make(map[string]string),
	}
}

func (message *Message) HasKey(key string) bool {
	_, exists := message.Information[key]
	return exists
}

func (message *Message) SetBool(key string, value bool) {
	message.Information[key] = strconv.FormatBool(value)
}

func (message *Message) GetBool(key string) (bool, error) {
	var err error = message.checkKey(key)
	var result bool = false

	if err == nil {
		result, err = strconv.ParseBool(message.Information[key])
	}

	return result, err
}

func (message *Message) SetString(key string, value string) {
	message.Information[key] = value
}

func (message *Message) GetString(key string) (string, error) {
	var err error = message.checkKey(key)
	var result string = ""

	if err == nil {
		result = message.Information[key]
	}

	return result, err
}

func (message *Message) SetInteger(key string, value int) {
	message.Information[key] = strconv.Itoa(value)
}

func (message *Message) GetInteger(key string) (int, error) {
	var err error = message.checkKey(key)
	var result int = 0

	if err == nil {
		result, err = strconv.Atoi(message.Information[key])
	}

	return result, err
}

func (message *Message) SetFloat(key string, value float64) {
	message.Information[key] = fmt.Sprintf("%v", value)
}

func (message *Message) GetFloat(key string) (float64, error) {
	var err error = message.checkKey(key)
	var result float64 = 0

	if err == nil {
		result, err = strconv.ParseFloat(message.Information[key], 64)
	}

	return result, err
}

func (message *Message) ToJson() []byte {
	value, _ := json.Marshal(message)
	return value
}

func (message *Message) ToString() string {
	return string(message.ToJson())
}

func (message *Message) checkKey(key string) error {
	var result error = nil

	if !message.HasKey(key) {
		result = fmt.Errorf("Key %s doesn't exists", key)
	}

	return result
}
