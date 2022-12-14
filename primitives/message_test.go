package primitives_test

import (
	"testing"

	"github.com/GiulianoDecesares/commvent/primitives"
)

const (
	messageType = "TESTING_TYPE"

	key = "TEST"

	stringValue        = "TESTING_VALUE"
	stringDefaultvalue = "DEFAULT"

	intValue        = 13
	intDefaultValue = 0

	boolValue        = true
	boolDefaultValue = !boolValue

	floatValue        float64 = 3.14
	defaultFloatValue float64 = 0
)

func TestEmptyMessage(context *testing.T) {
	message := primitives.NewMessage(messageType, "")

	if message.HasKey(key) {
		context.Errorf("%s key has not been added to message but is detected", key)
	}
}

func TestStringMessage(context *testing.T) {
	message := primitives.NewMessage(messageType, "")
	message.SetString(key, stringValue)

	if value, err := message.GetString(key); err == nil {
		if value != stringValue {
			context.Errorf("Value for key %s is not %s", key, stringValue)
		}
	} else {
		context.Errorf("Error while trying to get key %s: %s", key, err.Error())
	}
}

func TestIntegerMessage(context *testing.T) {
	message := primitives.NewMessage(messageType, "")
	message.SetInteger(key, intValue)

	if value, err := message.GetInteger(key); err == nil {
		if value != intValue {
			context.Errorf("Value for key %s is not %d", key, intValue)
		}
	} else {
		context.Errorf("Error while trying to get key %s: %s", key, err.Error())
	}
}

func TestBoolMessage(context *testing.T) {
	message := primitives.NewMessage(messageType, "")
	message.SetBool(key, boolValue)

	if value, err := message.GetBool(key); err == nil {
		if value != boolValue {
			context.Errorf("Value for key %s is not %t", key, boolValue)
		}
	} else {
		context.Errorf("Error while trying to get key %s: %s", key, err.Error())
	}
}

func TestFloatMessage(context *testing.T) {
	message := primitives.NewMessage(messageType, "")
	message.SetFloat(key, floatValue)

	if value, err := message.GetFloat(key); err == nil {
		if value != floatValue {
			context.Errorf("Value for key %s is not %f", key, floatValue)
		}
	} else {
		context.Errorf("Error while trying to get key %s: %s", key, err.Error())
	}
}
