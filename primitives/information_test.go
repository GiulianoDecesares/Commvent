package primitives_test

import (
	"strconv"
	"testing"

	"github.com/GiulianoDecesares/commvent/primitives"
)

const (
	key = "TEST"
)

const (
	stringValue        = "TESTING_VALUE"
	stringDefaultvalue = "DEFAULT"

	intValue        = 13
	intDefaultValue = 0

	boolValue        = true
	boolDefaultValue = !boolValue

	floatValue        float64 = 3.14
	defaultFloatValue float64 = 0
)

func TestInformationNoneParameter(context *testing.T) {
	info := primitives.NewInformation()

	if info.HasKey(key) {
		context.Errorf("%s key has not been added to information but is detected", key)
	}
}

func TestInformationStringParameter(context *testing.T) {
	info := primitives.NewInformation()

	info.SetString(key, stringValue)

	if !info.HasKey(key) {
		context.Errorf("%s key has been added to information but is not detected", key)
	} else if info.GetString(key, stringDefaultvalue) != stringValue {
		context.Errorf("Value for key %s is not %s", key, stringValue)
	}
}

func TestInformationIntParameter(context *testing.T) {
	info := primitives.NewInformation()

	info.SetInteger(key, intValue)

	if !info.HasKey(key) {
		context.Errorf("%s key has been added to information but is not detected", key)
	} else if info.GetInteger(key, intDefaultValue) != intValue {
		context.Errorf("Value for key %s is not %d", key, intValue)
	}
}

func TestInformationBoolParameter(context *testing.T) {
	info := primitives.NewInformation()

	info.SetBool(key, boolValue)

	if !info.HasKey(key) {
		context.Errorf("%s key has been added to information but is not detected", key)
	} else if info.GetBool(key, boolDefaultValue) != boolValue {
		context.Errorf("Value for key %s is not %s", key, strconv.FormatBool(boolValue))
	}
}

func TestInformationFloatParameter(context *testing.T) {
	info := primitives.NewInformation()

	info.SetFloat(key, floatValue)

	if !info.HasKey(key) {
		context.Errorf("%s key has been added to information but is not detected", key)
	} else if info.GetFloat(key, defaultFloatValue) != floatValue {
		context.Errorf("Value for key %s is not %f", key, floatValue)
	}
}
