package primitives

import "encoding/json"

type Command struct {
	Type      string      `json:"type"`
	Arguments Information `json:"args"`
}

func NewCommand(commandType string) *Command {
	return &Command{
		Type:      commandType,
		Arguments: *NewInformation(),
	}
}

func NewCommandFromJSON(bytes []byte) *Command {
	command := Command{}
	json.Unmarshal(bytes, &command)

	return &command
}

func (command *Command) ToJson() []byte {
	jsonData, _ := json.Marshal(command)
	return jsonData
}
