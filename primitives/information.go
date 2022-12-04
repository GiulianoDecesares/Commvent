package primitives

import (
	"fmt"
	"strconv"
)

type Information struct {
	Raw map[string]string `json:"raw"`
}

func NewInformation() *Information {
	return &Information{
		Raw: map[string]string{},
	}
}

func (info *Information) HasKey(key string) bool {
	_, exists := info.Raw[key]
	return exists
}

func (info *Information) SetBool(key string, value bool) {
	info.Raw[key] = strconv.FormatBool(value)
}

func (info *Information) GetBool(key string, defaultValue bool) bool {
	var result bool = defaultValue

	if value, exists := info.Raw[key]; exists {
		if booleanResult, err := strconv.ParseBool(value); err == nil {
			result = booleanResult
		}
	}

	return result
}

func (info *Information) SetString(key string, value string) {
	info.Raw[key] = value
}

func (info *Information) GetString(key string, defaultValue string) string {
	var result string = defaultValue

	if value, exists := info.Raw[key]; exists {
		result = value
	}

	return result
}

func (info *Information) SetInteger(key string, value int) {
	info.Raw[key] = strconv.Itoa(value)
}

func (info *Information) GetInteger(key string, defaultValue int) int {
	var result int = defaultValue

	if value, exists := info.Raw[key]; exists {
		if integer, err := strconv.Atoi(value); err == nil {
			result = integer
		}
	}

	return result
}

func (info *Information) SetFloat(key string, value float64) {
	info.Raw[key] = fmt.Sprintf("%v", value)
}

func (info *Information) GetFloat(key string, defaultValue float64) float64 {
	var result float64 = defaultValue

	if rawValue, exists := info.Raw[key]; exists {
		if value, err := strconv.ParseFloat(rawValue, 64); err == nil {
			result = value
		}
	}

	return result
}
