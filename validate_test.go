package bbcos

import (
	"path"
	"testing"
)

func TestValidateBuiltinSchema(t *testing.T) {
	builtinSchemas := []struct {
		Variant string
		Version string
	}{
		{"fcos", "1.0.0"},
		// {"fcos", "1.1.0"},
		// {"fcos", "1.2.0"},
		// {"fcos", "1.3.0"},
		// {"fcos", "1.4.0"},
		// {"openshift", "4.8.0"},
		// {"openshift", "4.9.0"},
		// {"rhcos", "0.1.0"},
	}

	for _, builtinSchema := range builtinSchemas {
		filename := builtinSchema.Variant + "_v" + builtinSchema.Version + ".schema.json"
		t.Run(filename, func(subT *testing.T) {
			defer func() {
				err := recover()
				if err != nil {
					subT.Fail()
					return
				}
			}()

			spec, err := specs.Open(path.Join("spec", filename))
			if err != nil {
				subT.Error(err)
				return
			}

			compileSchema(builtinSchema.Variant, spec)
		})
	}
}

func TestCustomSchemaValidation(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fail()
			return
		}
	}()

	spec, err := specs.Open(path.Join("spec", "fcos_v1.0.0.schema"))
	if err != nil {
		panic(err)
	}

	NewButaneConfig(FCOS, FCOS_1_0_0).
		ValidateWithSpec(spec)
}

func TestEmptyConfigValidation(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fail()
			return
		}
	}()

	b := NewButaneConfig("", "")
	b.Validate()
}

func TestMissingRequiredIgnitionFieldValidation(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fail()
			return
		}
	}()

	NewButaneConfig(FCOS, FCOS_1_0_0).
		Validate()
}
