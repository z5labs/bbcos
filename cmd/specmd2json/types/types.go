// Package types contains the supported types for JSON schemas.
package types

// Type represents a type which can be validated by a JSON schema.
type Type string

// String maps a type to its text based representation.
func (t Type) String() string {
	return string(t)
}

const (
	Err Type = "error"
	EOF Type = "eof"

	Object   Type = "object"
	Objects  Type = "list of object"
	Array    Type = "array"
	Arrays   Type = "list of array"
	String   Type = "string"
	Strings  Type = "list of string"
	Boolean  Type = "boolean"
	Booleans Type = "list of boolean"
	Integer  Type = "integer"
	Integers Type = "list of integer"
)

// IsArray checks whether or not a given type
// represents an array of some other type.
// e.g. strings -> array of strings
//
func IsArray(t Type) bool {
	return t == Array || t == Arrays || t == Objects || t == Strings || t == Booleans || t == Integers
}

func Normalize(t Type) string {
	if IsArray(t) {
		return "array"
	}
	return string(t)
}

func InnerType(t Type) Type {
	switch t {
	case Arrays:
		return Array
	case Objects:
		return Object
	case Strings:
		return String
	case Booleans:
		return Boolean
	case Integers:
		return Integer
	default:
		return t
	}
}
