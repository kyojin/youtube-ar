package log

import (
	"fmt"
	"strconv"
)

// Field is the field interface.
type Field interface {
	key() string
	value() string
}

type field struct{ k, v string }

func (f field) key() string   { return f.k }
func (f field) value() string { return f.v }

// Int returns a field where value is an int.
func Int(key string, value int) Field {
	return field{k: key, v: strconv.Itoa(value)}
}

// String returns a field where value is a string.
func String(key string, value string) Field {
	return field{k: key, v: value}
}

// Stringer returns a field where value implements fmt.Stringer.
func Stringer(key string, value fmt.Stringer) Field {
	return field{k: key, v: value.String()}
}

// Error returns a field where value is an error.
func Error(key string, value error) Field {
	return field{k: key, v: value.Error()}
}

// Raw returns a field where value is an empty interface{}.
func Raw(key string, value interface{}) Field {
	return field{k: key, v: fmt.Sprintf("%+v", value)}
}
