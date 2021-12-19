package parser

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/z5labs/bbcos/cmd/specmd2json/keys"
	"github.com/z5labs/bbcos/cmd/specmd2json/types"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

const specMdUrl = "https://raw.githubusercontent.com/coreos/butane/main/docs/config-fcos-v1_0.md"

// func TestParse(t *testing.T) {
// 	specSchemaFilePath, err := filepath.Abs(specSchemaFilePath)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
//
// 	expectedSpecSchema, err := ioutil.ReadFile(specSchemaFilePath)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
//
// 	resp, err := http.Get(specMdUrl)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	defer resp.Body.Close()
//
// 	spec := Parse(resp.Body)
//
// 	var b bytes.Buffer
// 	err = writeTo(spec, &b)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	actualSpecSchema := b.Bytes()
//
// 	opts := jsondiff.DefaultJSONOptions()
// 	diff, s := jsondiff.Compare(expectedSpecSchema, actualSpecSchema, &opts)
// 	if diff != jsondiff.FullMatch {
// 		t.Fail()
// 		t.Log(s)
// 		return
// 	}
// }

func TestParseObject(t *testing.T) {
	md := []byte(`
* **obj** (object): obj
	* **one** (string): one
	* **two** (string): two
	* **thr** (string): thr`)

	spec := Parse(bytes.NewReader(md))
	spec = spec.Sub(keys.Properties.String())

	if spec.Get("obj") == nil {
		t.Log("expected obj to be set")
		t.Fail()
		return
	}

	obj := spec.Sub("obj")
	ty := obj.GetString(keys.Type.String())
	descr := obj.GetString(keys.Description.String())
	if ty == "" || ty != "object" {
		t.Log("expected type be set to:", "object")
		t.Fail()
		return
	}
	if descr == "" || descr != "obj" {
		t.Log("expected description be set to:", "obj")
		t.Fail()
		return
	}

	fields := obj.Sub(keys.Properties.String())
	if fields == nil {
		t.Log("expected obj to be an object with properties")
		t.Fail()
		return
	}

	for _, id := range []string{"one", "two", "thr"} {
		if fields.Get(id) == nil {
			t.Log("expected", id, "to be set")
			t.Fail()
			return
		}

		field := fields.Sub(id)
		descr := field.GetString(keys.Description.String())
		if descr == "" {
			t.Log("expected description to be set")
			t.Fail()
			return
		}
		if descr != id {
			t.Log("expected description to be:", id)
			t.Fail()
			return
		}
	}
}

func TestParseObjects(t *testing.T) {
	md := []byte(`
* **obj1** (object): obj1
	* **one** (string): one
	* **two** (string): two
	* **thr** (string): thr
* **obj2** (object): obj2
	* **one** (string): one
	* **two** (string): two
	* **thr** (string): thr`)

	spec := Parse(bytes.NewReader(md))
	spec = spec.Sub(keys.Properties.String())

	if spec.Get("obj1") == nil {
		t.Log("expected obj1 to be set")
		t.Fail()
		return
	}

	obj := spec.Sub("obj1")
	ty := obj.GetString(keys.Type.String())
	descr := obj.GetString(keys.Description.String())
	if ty == "" || ty != "object" {
		t.Log("expected type be set to:", "object")
		t.Fail()
		return
	}
	if descr == "" || descr != "obj1" {
		t.Log("expected description be set to:", "obj1")
		t.Fail()
		return
	}

	fields := obj.Sub(keys.Properties.String())
	if fields == nil {
		t.Log("expected obj2 to be an object with properties")
		t.Fail()
		return
	}

	for _, id := range []string{"one", "two", "thr"} {
		if fields.Get(id) == nil {
			t.Log("expected", id, "to be set")
			t.Fail()
			return
		}

		field := fields.Sub(id)
		descr := field.GetString(keys.Description.String())
		if descr == "" {
			t.Log("expected description to be set")
			t.Fail()
			return
		}
		if descr != id {
			t.Log("expected description to be:", id)
			t.Fail()
			return
		}
	}

	if spec.Get("obj2") == nil {
		t.Log("expected obj2 to be set")
		t.Fail()
		return
	}

	obj = spec.Sub("obj2")
	ty = obj.GetString(keys.Type.String())
	descr = obj.GetString(keys.Description.String())
	if ty == "" || ty != "object" {
		t.Log("expected type be set to:", "object")
		t.Fail()
		return
	}
	if descr == "" || descr != "obj2" {
		t.Log("expected description be set to:", "obj2")
		t.Fail()
		return
	}

	fields = obj.Sub(keys.Properties.String())
	if fields == nil {
		t.Log("expected obj2 to be an object with properties")
		t.Fail()
		return
	}

	for _, id := range []string{"one", "two", "thr"} {
		if fields.Get(id) == nil {
			t.Log("expected", id, "to be set")
			t.Fail()
			return
		}

		field := fields.Sub(id)
		descr := field.GetString(keys.Description.String())
		if descr == "" {
			t.Log("expected description to be set")
			t.Fail()
			return
		}
		if descr != id {
			t.Log("expected description to be:", id)
			t.Fail()
			return
		}
	}
}

func TestNestedObject(t *testing.T) {
	md := []byte(`
* **obj** (object): obj
	* **a** (object): a
		* **a** (string): a
	* **b** (object): b
		* **b** (string): b
	* **c** (object): c
		* **c** (string): c`)

	spec := Parse(bytes.NewReader(md))
	spec = spec.Sub(keys.Properties.String())

	if spec.Get("obj") == nil {
		t.Log("expected obj to be set")
		t.Fail()
		return
	}

	obj := spec.Sub("obj")
	ty := obj.GetString(keys.Type.String())
	descr := obj.GetString(keys.Description.String())
	if ty == "" || ty != "object" {
		t.Log("expected type be set to:", types.Object)
		t.Fail()
		return
	}
	if descr == "" || descr != "obj" {
		t.Log("expected description be set to:", "obj")
		t.Fail()
		return
	}

	fields := obj.Sub(keys.Properties.String())
	if fields == nil {
		t.Log("expected obj to be an object with properties")
		t.Fail()
		return
	}

	for _, id := range []string{"a", "b", "c"} {
		if fields.Get(id) == nil {
			t.Log("expected", id, "to be set")
			t.Fail()
			return
		}

		field := fields.Sub(id)
		ty := field.GetString(keys.Type.String())
		descr := field.GetString(keys.Description.String())
		if ty == "" || ty != "object" {
			t.Log("expected type be set to:", types.Object)
			t.Fail()
			return
		}
		if descr == "" || descr != id {
			t.Log("expected description be set to:", descr)
			t.Fail()
			return
		}

		subFields := field.Sub(keys.Properties.String())
		if subFields == nil {
			t.Log("expected", id, "to be an object with properties:", field.AllSettings())
			t.Fail()
			return
		}

		subField := subFields.Sub(id)
		ty = subField.GetString(keys.Type.String())
		descr = subField.GetString(keys.Description.String())
		if ty == "" || ty != "string" {
			t.Log("expected type be set to:", types.String)
			t.Fail()
			return
		}
		if descr == "" || descr != id {
			t.Log("expected description be set to:", descr)
			t.Fail()
			return
		}
	}
}

func TestArrayOfObjects(t *testing.T) {
	md := []byte(`
* **arr** (list of objects): arr
	* **one** (string): one
	* **two** (string): two
	* **thr** (string): thr`)

	spec := Parse(bytes.NewReader(md))
	spec = spec.Sub(keys.Properties.String())

	if spec.Get("arr") == nil {
		t.Log("expected arr to be set")
		t.Fail()
		return
	}

	obj := spec.Sub("arr")
	ty := obj.GetString(keys.Type.String())
	descr := obj.GetString(keys.Description.String())
	if ty == "" || ty != types.Normalize(types.Objects) {
		t.Log("expected type be set to:", types.Objects)
		t.Fail()
		return
	}
	if descr == "" || descr != "arr" {
		t.Log("expected description be set to:", "arr")
		t.Fail()
		return
	}

	var items struct {
		Type       string
		Properties map[string]interface{}
	}

	err := obj.UnmarshalKey("items", &items)
	if err != nil {
		t.Log("unexpected error when unmarshalling array items:", err)
		t.Fail()
		return
	}
	if items.Type != types.Object.String() {
		t.Logf("expected array items type to be: %s", types.Object)
		t.Fail()
		return
	}

	for _, id := range []string{"one", "two", "thr"} {
		item, ok := items.Properties[id]
		if !ok || item == nil {
			t.Log("expected", id, "to be set")
			t.Fail()
			return
		}
		itemMap := item.(map[string]interface{})

		typ, ok := itemMap[keys.Type.String()]
		if !ok || typ == nil {
			t.Log("expected type to be set")
			t.Fail()
			return
		}
		if typ.(string) != types.String.String() {
			t.Log("expected type to be:", types.String)
			t.Fail()
			return
		}

		descr, ok := itemMap[keys.Description.String()]
		if !ok || descr == nil {
			t.Log("expected description to be set")
			t.Fail()
			return
		}
		if descr.(string) != id {
			t.Log("expected description to be:", id)
			t.Fail()
			return
		}
	}
}

func TestArrayOfStrings(t *testing.T) {
	md := []byte(`* **arr** (list of strings): arr`)

	spec := Parse(bytes.NewReader(md))
	spec = spec.Sub(keys.Properties.String())

	if spec.Get("arr") == nil {
		t.Log("expected arr to be set")
		t.Fail()
		return
	}

	obj := spec.Sub("arr")
	ty := obj.GetString(keys.Type.String())
	descr := obj.GetString(keys.Description.String())
	if ty == "" || ty != types.Normalize(types.Strings) {
		t.Log("expected type be set to:", types.Strings)
		t.Fail()
		return
	}
	if descr == "" || descr != "arr" {
		t.Log("expected description be set to:", "arr")
		t.Fail()
		return
	}

	items := obj.Sub(keys.Items.String())
	if items == nil {
		t.Logf("expected array items to be set: %v", obj.AllSettings())
		t.Fail()
		return
	}

	ty = items.GetString(keys.Type.String())
	if ty != types.String.String() {
		t.Logf("expected array items type to be: %s", types.String)
		t.Fail()
		return
	}
}

func TestNestedArrays(t *testing.T) {
	md := []byte(`
* **arr** (list of objects): arr
	* **a** (list of objects): a
		* **a** (list of string): a
	* **b** (list of objects): b
		* **b** (list of string): b
	* **c** (list of objects): c
		* **c** (list of string): c`)

	spec := Parse(bytes.NewReader(md))
	spec = spec.Sub(keys.Properties.String())

	if spec.Get("arr") == nil {
		t.Log("expected arr to be set")
		t.Fail()
		return
	}

	arr := spec.Sub("arr")
	ty := arr.GetString(keys.Type.String())
	descr := arr.GetString(keys.Description.String())
	if ty == "" || ty != "array" {
		t.Log("expected type be set to:", types.Objects)
		t.Fail()
		return
	}
	if descr == "" || descr != "arr" {
		t.Log("expected description be set to:", "arr")
		t.Fail()
		return
	}

	var items struct {
		Type       string
		Properties map[string]interface{}
	}

	err := arr.UnmarshalKey("items", &items)
	if err != nil {
		t.Log("unexpected error when unmarshalling array items:", err)
		t.Fail()
		return
	}
	if items.Type != types.Object.String() {
		t.Logf("expected array items type to be: %s", types.Object)
		t.Fail()
		return
	}

	for _, id := range []string{"a", "b", "c"} {
		item, ok := items.Properties[id]
		if !ok || item == nil {
			t.Log("expected", id, "to be set")
			t.Fail()
			return
		}
		itemMap := item.(map[string]interface{})

		typ, ok := itemMap[keys.Type.String()]
		if !ok || typ == nil {
			t.Log("expected type to be set")
			t.Fail()
			return
		}
		if typ.(string) != types.Array.String() {
			t.Logf("expected type to be: %s (actual: %s)", types.String, typ)
			t.Fail()
			return
		}

		descr, ok := itemMap[keys.Description.String()]
		if !ok || descr == nil {
			t.Log("expected description to be set")
			t.Fail()
			return
		}
		if descr.(string) != id {
			t.Logf("expected description to be: %s (actual: %s)", id, descr)
			t.Fail()
			return
		}

		subItems, ok := itemMap[keys.Items.String()]
		if !ok || subItems == nil {
			t.Log("expected", keys.Items, "to be set")
			t.Fail()
			return
		}
		subItemsMap := subItems.(map[string]interface{})

		typ, ok = subItemsMap[keys.Type.String()]
		if !ok || typ == nil {
			t.Log("expected type to be set", subItemsMap)
			t.Fail()
			return
		}
		if typ.(string) != types.Object.String() {
			t.Logf("expected type to be: %s (actual: %s)", types.String, typ)
			t.Fail()
			return
		}

		subProps, ok := subItemsMap[keys.Properties.String()]
		if !ok || subProps == nil {
			t.Log("expected properties to be set")
			t.Fail()
			return
		}
		subPropsMap := subProps.(map[string]interface{})

		subProp, ok := subPropsMap[id].(map[string]interface{})
		if !ok || subProp == nil {
			t.Log("expected sub property to be present")
			t.Fail()
			return
		}

		typ, ok = subProp[keys.Type.String()]
		if !ok || typ == nil {
			t.Log("expected type to be set")
			t.Fail()
			return
		}
		if typ.(string) != types.Array.String() {
			t.Logf("expected type to be: %s (actual: %s)", types.Array, typ)
			t.Fail()
			return
		}

		descr, ok = subProp[keys.Description.String()]
		if !ok || descr == nil {
			t.Log("expected description to be set")
			t.Fail()
			return
		}
		if descr.(string) != id {
			t.Logf("expected description to be: %s (actual: %s)", id, descr)
			t.Fail()
			return
		}

		subSubItems, ok := subProp[keys.Items.String()]
		if !ok || subSubItems == nil {
			t.Log("expected items to be set")
			t.Fail()
			return
		}
		subSubItemsMap := subSubItems.(map[string]interface{})

		typ, ok = subSubItemsMap[keys.Type.String()]
		if !ok || typ == nil {
			t.Log("expected type to be set")
			t.Fail()
			return
		}
		if typ.(string) != types.String.String() {
			t.Logf("expected type to be: %s (actual: %s)", types.String, typ)
			t.Fail()
			return
		}
	}
}

func TestBuildPathWith(t *testing.T) {
	testCases := []struct {
		Name string
		Path []string
		Arrs []bool
		Vals []string
		Ex   string
	}{
		{
			Name: "rootPropField",
			Path: []string{},
			Arrs: []bool{},
			Vals: []string{"a"},
			Ex:   "a",
		},
		{
			Name: "objProp",
			Path: []string{},
			Arrs: []bool{},
			Vals: []string{"obj", "type"},
			Ex:   "obj.type",
		},
		{
			Name: "nestedObjProp",
			Path: []string{"obj"},
			Arrs: []bool{false},
			Vals: []string{"properties", "a", "type"},
			Ex:   "properties.obj.properties.a.type",
		},
		{
			Name: "arrayProp",
			Path: []string{},
			Arrs: []bool{},
			Vals: []string{"arr", "type"},
			Ex:   "arr.type",
		},
		{
			Name: "arrayObjectPropField",
			Path: []string{"arr", "obj"},
			Arrs: []bool{true, false},
			Vals: []string{"properties", "a", "type"},
			Ex:   "properties.arr.items.properties.obj.properties.a.type",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(subT *testing.T) {
			actual := buildPathWith(bufPool, testCase.Path, testCase.Arrs, testCase.Vals...)

			if testCase.Ex != actual {
				subT.Log(actual)
				subT.Fail()
				return
			}
		})
	}
}

func TestMarkArrayStatus(t *testing.T) {
	testCases := []struct {
		Name  string
		Arrs  []bool
		Path  []string
		IsArr bool
		Ex    []bool
	}{
		{
			Name:  "root",
			Arrs:  []bool{},
			Path:  []string{},
			IsArr: true,
			Ex:    []bool{true},
		},
		{
			Name:  "withParent",
			Arrs:  []bool{true},
			Path:  []string{""},
			IsArr: true,
			Ex:    []bool{true, true},
		},
		{
			Name:  "sibling",
			Arrs:  []bool{true, true},
			Path:  []string{""},
			IsArr: false,
			Ex:    []bool{true, false},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(subT *testing.T) {
			out := markArrayStatus(testCase.Arrs, testCase.Path, testCase.IsArr)
			if len(testCase.Ex) != len(out) {
				subT.Log("expected slice lengths to match:", testCase.Ex, out)
				subT.Fail()
				return
			}

			for i, ex := range testCase.Ex {
				if ex != out[i] {
					subT.Log("expected values to match:", ex, out[i])
					subT.Fail()
					return
				}
			}
		})
	}
}

func BenchmarkParse(b *testing.B) {
	resp, err := http.Get(specMdUrl)
	if err != nil {
		b.Error(err)
		return
	}
	defer resp.Body.Close()

	mdBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		b.Error(err)
		return
	}
	if err = resp.Body.Close(); err != nil {
		b.Error(err)
		return
	}
	var buf bytes.Buffer

	for i := 0; i < b.N; i++ {
		buf.Reset()
		buf.Write(mdBytes)

		spec := Parse(&buf)
		if spec == nil {
			b.Fail()
		}
	}
}

func writeTo(v *viper.Viper, w io.Writer) error {
	fs := afero.NewMemMapFs()
	v.SetFs(fs)
	v.SetConfigFile("test.json")
	err := v.WriteConfig()
	if err != nil {
		return err
	}

	f, err := fs.Open("test.json")
	if err != nil {
		return err
	}

	_, err = io.Copy(w, f)
	return err
}
