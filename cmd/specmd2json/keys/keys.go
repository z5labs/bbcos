// Package keys provides valid JSON schema keys.
package keys

// Key represents a valid JSON schema key.
type Key string

// String maps the key to its corresponding text encoding e.g. ID -> $id.
func (k Key) String() string {
	return string(k)
}

const (
	Schema      Key = "$schema"
	ID          Key = "$id"
	Title       Key = "title"
	Description Key = "description"
	Type        Key = "type"
	Properties  Key = "properties"
	Items       Key = "items"
)
