package primitives_test

import (
	"testing"

	"github.com/GiulianoDecesares/commvent/primitives"
)

const (
	eventType = "TEST_EVENT"

	emptyEventString        = "{\"type\":\"TEST_EVENT\",\"info\":{\"raw\":{}}}"
	parametrizedEventString = "{\"type\":\"TEST_EVENT\",\"info\":{\"raw\":{\"TEST\":\"true\",\"TEST_INTEGER\":\"13\",\"TEST_STRING\":\"TESTING_VALUE\"}}}"
)

const (
	integerKey = key + "_INTEGER"
	stringKey  = key + "_STRING"
)

func TestEventType(context *testing.T) {
	event := primitives.NewEvent(eventType)

	if event.Type != eventType {
		context.Errorf("Event type should be %s but is %s", eventType, event.Type)
	}
}

func TestEmptyEventToJSON(context *testing.T) {
	event := primitives.NewEvent(eventType)
	rawJson := string(event.ToJson())

	if rawJson != emptyEventString {
		context.Errorf("Empty JSON event should be %s but is %s", emptyEventString, rawJson)
	}
}

func TestParametrizedEventToJSON(context *testing.T) {
	event := primitives.NewEvent(eventType)

	event.Information.SetBool(key, boolValue)
	event.Information.SetString(stringKey, stringValue)
	event.Information.SetInteger(integerKey, intValue)

	rawJson := string(event.ToJson())

	if rawJson != parametrizedEventString {
		context.Errorf("Parametrized JSON event should be %s but is %s", parametrizedEventString, rawJson)
	}
}

func TestNewEventFromJson(context *testing.T) {
	event := primitives.NewEventFromJSON([]byte(parametrizedEventString))

	if !event.Information.HasKey(stringKey) {
		context.Errorf("Event from JSON should have %s key", stringKey)
	}

	if !event.Information.HasKey(integerKey) {
		context.Errorf("Event from JSON should have %s key", integerKey)
	}

	if event.Information.GetInteger(integerKey, intDefaultValue) != intValue {
		context.Errorf("Key %s should have value %d", integerKey, intValue)
	}

	if event.Information.GetString(stringKey, stringDefaultvalue) != stringValue {
		context.Errorf("Key %s should have value %s", stringKey, stringValue)
	}
}
