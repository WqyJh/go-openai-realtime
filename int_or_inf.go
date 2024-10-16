package openairt

import (
	"encoding/json"
)

// IntOrInf is a type that can be either an int or "inf".
type IntOrInf struct {
	value int
	inf   bool
}

// Inf returns an IntOrInf with the value "inf".
func Inf() IntOrInf {
	return IntOrInf{
		inf: true,
	}
}

// Int returns an IntOrInf with the given value.
func Int(value int) IntOrInf {
	return IntOrInf{
		value: value,
	}
}

// IsInf returns true if the value is "inf".
func (m IntOrInf) IsInf() bool {
	return m.inf
}

// Value returns the value as an int.
func (m IntOrInf) Value() int {
	return m.value
}

// MarshalJSON marshals the IntOrInf to JSON.
func (m IntOrInf) MarshalJSON() ([]byte, error) {
	if m.inf {
		return []byte("\"inf\""), nil
	}
	return json.Marshal(m.value)
}

// UnmarshalJSON unmarshals the IntOrInf from JSON.
func (m *IntOrInf) UnmarshalJSON(data []byte) error {
	if string(data) == "\"inf\"" {
		m.inf = true
		return nil
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, &m.value)
}
