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
func (cfg ButaneConfig) ValidateWithSpec(spec io.Reader) ButaneConfig {
	cfg.mustBeValidAgainstSpec("custom.schema", spec)
	return cfg
}

// Validate is used to validate your configuration against a particaular spec
// matching the set variation and version fields.
//
func (cfg ButaneConfig) Validate() ButaneConfig {
	if !cfg.v.IsSet("variant") || !cfg.v.IsSet("version") {
		panic("bbcos: both variation and version must be set before validating")
	}

	variant := strings.ToLower(cfg.v.GetString("variant"))
	version := strings.TrimPrefix(strings.ToLower(cfg.v.GetString("version")), "v")
	filename := variant + "_v" + version + ".schema.json"
	spec, err := specs.Open(path.Join("spec", filename))
	if err != nil {
		panic(err)
	}

	cfg.mustBeValidAgainstSpec(filename, spec)
	return cfg
}

func (cfg ButaneConfig) mustBeValidAgainstSpec(specName string, spec io.Reader) {
	schema := compileSchema(specName, spec)
	jsonCfg := cfg.jsonifyConfig()
	if err := schema.Validate(jsonCfg); err != nil {
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

func (cfg ButaneConfig) jsonifyConfig() interface{} {
	fs := afero.NewMemMapFs()
	cfg.v.SetFs(fs)
	cfg.v.SetConfigFile("butane.json")
	err := cfg.v.WriteConfig()
	if err != nil {
		panic(err)
	}

	f, err := fs.Open("butane.json")
	if err != nil {
		panic(err)
	}

	var jsonCfg interface{}
	d := json.NewDecoder(f)
	if err = d.Decode(&jsonCfg); err != nil {
		panic(err)
	}
	return jsonCfg
}
