package primitives_test

import (
	"testing"

	"github.com/GiulianoDecesares/commvent/primitives"
)

const (
	commandType = "TEST_COMMAND"

	emptyCommandString        = "{\"type\":\"TEST_COMMAND\",\"args\":{\"raw\":{}}}"
	parametrizedCommandString = "{\"type\":\"TEST_COMMAND\",\"args\":{\"raw\":{\"TEST\":\"true\",\"TEST_INTEGER\":\"13\",\"TEST_STRING\":\"TESTING_VALUE\"}}}"
)

func TestCommandType(context *testing.T) {
	command := primitives.NewCommand(commandType)

	if command.Type != commandType {
		context.Errorf("Command type should be %s but is %s", commandType, command.Type)
	}
}

func TestEmptyCommandToJSON(context *testing.T) {
	command := primitives.NewCommand(commandType)
	rawJson := string(command.ToJson())

	if rawJson != emptyCommandString {
		context.Errorf("Empty JSON command should be %s but is %s", emptyCommandString, rawJson)
	}
}

func TestParametrizedCommandToJSON(context *testing.T) {
	command := primitives.NewCommand(commandType)

	command.Arguments.SetBool(key, boolValue)
	command.Arguments.SetString(stringKey, stringValue)
	command.Arguments.SetInteger(integerKey, intValue)

	rawJson := string(command.ToJson())

	if rawJson != parametrizedCommandString {
		context.Errorf("Parametrized JSON command should be %s but is %s", parametrizedCommandString, rawJson)
	}
}

func TestNewCommandFromJson(context *testing.T) {
	command := primitives.NewCommandFromJSON([]byte(parametrizedCommandString))

	if !command.Arguments.HasKey(stringKey) {
		context.Errorf("Command from JSON should have %s key", stringKey)
	}

	if !command.Arguments.HasKey(integerKey) {
		context.Errorf("Command from JSON should have %s key", integerKey)
	}

	if command.Arguments.GetInteger(integerKey, intDefaultValue) != intValue {
		context.Errorf("Key %s should have value %d", integerKey, intValue)
	}

	if command.Arguments.GetString(stringKey, stringDefaultvalue) != stringValue {
		context.Errorf("Key %s should have value %s", stringKey, stringValue)
	}
}
