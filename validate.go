// validate.go provides functionality for validating a config against
// a particular spec.

package bbcos

import (
	"embed"
	"encoding/json"
	"io"
	"io/ioutil"
	"path"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/spf13/afero"
)

//go:embed spec
var specs embed.FS

// ValidateWithSpec is used to validate your configuration against a provided spec.
//
// The spec should be a JSON schema. For help formatting your spec,
// please reference the out-of-box supported specs in the spec folder.
//
func (b *Builder) ValidateWithSpec(spec io.Reader) *Builder {
	b.mustBeValidAgainstSpec("custom.schema", spec)
	return b
}

// Validate is used to validate your configuration against a particaular spec
// matching the set variation and version fields.
//
func (b *Builder) Validate() *Builder {
	if !b.v.IsSet("variant") || !b.v.IsSet("version") {
		panic("bbcos: both variation and version must be set before validating")
	}

	variant := strings.ToLower(b.v.GetString("variant"))
	version := strings.TrimPrefix(strings.ToLower(b.v.GetString("version")), "v")
	filename := variant + "_v" + version + ".schema"
	spec, err := specs.Open(path.Join("spec", filename))
	if err != nil {
		panic(err)
	}

	b.mustBeValidAgainstSpec(filename, spec)
	return b
}

func (b *Builder) mustBeValidAgainstSpec(specName string, spec io.Reader) {
	schema := compileSchema(specName, spec)
	cfg := b.jsonifyConfig()
	if err := schema.Validate(cfg); err != nil {
		panic(err)
	}
}

func compileSchema(filename string, spec io.Reader) *jsonschema.Schema {
	schemaB, err := ioutil.ReadAll(spec)
	if err != nil {
		panic(err)
	}

	return jsonschema.MustCompileString(filename, string(schemaB))
}

func (b *Builder) jsonifyConfig() interface{} {
	fs := afero.NewMemMapFs()
	b.v.SetFs(fs)
	b.v.SetConfigFile("butane.json")
	err := b.v.WriteConfig()
	if err != nil {
		panic(err)
	}

	f, err := fs.Open("butane.json")
	if err != nil {
		panic(err)
	}

	var cfg interface{}
	d := json.NewDecoder(f)
	if err = d.Decode(&cfg); err != nil {
		panic(err)
	}
	return cfg
}
