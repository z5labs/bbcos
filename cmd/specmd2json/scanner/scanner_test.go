package scanner

import (
	"bytes"
	"testing"

	"github.com/z5labs/bbcos/cmd/specmd2json/types"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

func TestScanList(t *testing.T) {
	md := `
* **one** (string): one
* **two** (string): two
* **thr** (string): thr`

	items := []Item{
		Item{Path: []string{}, Name: "one", Type: types.String, Description: "one"},
		Item{Path: []string{}, Name: "two", Type: types.String, Description: "two"},
		Item{Path: []string{}, Name: "thr", Type: types.String, Description: "thr"},
		Item{Path: []string{}, Name: "eof", Type: types.EOF, Description: "eof"},
	}

	b := []byte(md)

	mdParser := goldmark.DefaultParser()
	node := mdParser.Parse(text.NewReader(b))
	if !node.HasChildren() {
		panic("scanner: expecting root not to contain children")
	}

	s := New(bytes.NewReader(b))
	for _, exItem := range items {
		compare(t, exItem, s.NextItem())
	}
}

func TestScanObject(t *testing.T) {
	md := `
* **obj** (object): obj
  * **one** (string): one
  * **two** (string): two
  * **thr** (string): thr`

	items := []Item{
		Item{Path: []string{}, Name: "obj", Type: types.Object, Description: "obj"},
		Item{Path: []string{"obj"}, Name: "one", Type: types.String, Description: "one"},
		Item{Path: []string{"obj"}, Name: "two", Type: types.String, Description: "two"},
		Item{Path: []string{"obj"}, Name: "thr", Type: types.String, Description: "thr"},
		Item{Path: []string{}, Name: "eof", Type: types.EOF, Description: "eof"},
	}

	b := []byte(md)

	mdParser := goldmark.DefaultParser()
	node := mdParser.Parse(text.NewReader(b))
	if !node.HasChildren() {
		panic("scanner: expecting root not to contain children")
	}

	s := New(bytes.NewReader(b))
	for _, exItem := range items {
		compare(t, exItem, s.NextItem())
	}
}

func TestScanObjects(t *testing.T) {
	md := `
* **a** (object): a
  * **one** (string): one
  * **two** (string): two
* **b** (object): b
  * **one** (string): one
  * **two** (string): two`

	items := []Item{
		Item{Path: []string{}, Name: "a", Type: types.Object, Description: "a"},
		Item{Path: []string{"a"}, Name: "one", Type: types.String, Description: "one"},
		Item{Path: []string{"a"}, Name: "two", Type: types.String, Description: "two"},
		Item{Path: []string{}, Name: "b", Type: types.Object, Description: "b"},
		Item{Path: []string{"b"}, Name: "one", Type: types.String, Description: "one"},
		Item{Path: []string{"b"}, Name: "two", Type: types.String, Description: "two"},
		Item{Path: []string{}, Name: "eof", Type: types.EOF, Description: "eof"},
	}

	b := []byte(md)

	mdParser := goldmark.DefaultParser()
	node := mdParser.Parse(text.NewReader(b))
	if !node.HasChildren() {
		panic("scanner: expecting root not to contain children")
	}

	s := New(bytes.NewReader(b))
	for _, exItem := range items {
		compare(t, exItem, s.NextItem())
	}
}

func TestNestedObject(t *testing.T) {
	md := `
* **obj** (object): obj
  * **one** (object): one
    * **a** (integer): a
  * **two** (object): two
    * **b** (integer): b`

	items := []Item{
		Item{Path: []string{}, Name: "obj", Type: types.Object, Description: "obj"},
		Item{Path: []string{"obj"}, Name: "one", Type: types.Object, Description: "one"},
		Item{Path: []string{"obj", "one"}, Name: "a", Type: types.Integer, Description: "a"},
		Item{Path: []string{"obj"}, Name: "two", Type: types.Object, Description: "two"},
		Item{Path: []string{"obj", "two"}, Name: "b", Type: types.Integer, Description: "b"},
		Item{Path: []string{}, Name: "eof", Type: types.EOF, Description: "eof"},
	}

	b := []byte(md)

	mdParser := goldmark.DefaultParser()
	node := mdParser.Parse(text.NewReader(b))
	if !node.HasChildren() {
		panic("scanner: expecting root not to contain children")
	}

	s := New(bytes.NewReader(b))
	for _, exItem := range items {
		compare(t, exItem, s.NextItem())
	}
}

func TestScanArray(t *testing.T) {
	md := `
* **arr** (list of objects): arr
  * **one** (string): one
  * **two** (string): two
  * **thr** (string): thr`

	items := []Item{
		Item{Path: []string{}, Name: "arr", Type: types.Objects, Description: "arr"},
		Item{Path: []string{"arr"}, Name: "one", Type: types.String, Description: "one"},
		Item{Path: []string{"arr"}, Name: "two", Type: types.String, Description: "two"},
		Item{Path: []string{"arr"}, Name: "thr", Type: types.String, Description: "thr"},
		Item{Path: []string{}, Name: "eof", Type: types.EOF, Description: "eof"},
	}

	b := []byte(md)

	mdParser := goldmark.DefaultParser()
	node := mdParser.Parse(text.NewReader(b))
	if !node.HasChildren() {
		panic("scanner: expecting root not to contain children")
	}

	s := New(bytes.NewReader(b))
	for _, exItem := range items {
		item := s.NextItem()
		t.Log(item)
		compare(t, exItem, item)
	}
}

func TestScanItemText(t *testing.T) {
	testCases := []struct {
		Name string
		Text string
		Type types.Type
	}{
		{Name: "object", Text: "hello (object): world", Type: types.Object},
		{Name: "objects", Text: "hello (objects): world", Type: types.Object},
		{Name: "string", Text: "hello (string): world", Type: types.String},
		{Name: "list of string", Text: "hello (list of string): world", Type: types.Strings},
	}

	s := new(scanr)
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(subT *testing.T) {
			key, ty, descr := s.scanItemText(testCase.Text)
			if key != "hello" {
				subT.Log("expected key to be:", "hello", "but got", key)
				subT.Fail()
				return
			}

			if ty != testCase.Type {
				subT.Log("expected type to be:", testCase.Type, "but got", ty)
				subT.Fail()
				return
			}

			if descr != "world" {
				subT.Log("expected description to be:", "world", "but got", descr)
				subT.Fail()
				return
			}
		})
	}
}

func compare(t *testing.T, a, b Item) {
	if len(a.Path) != len(b.Path) {
		t.Logf("expected paths to be the same: %v:%v", a.Path, b.Path)
		t.Fail()
		return
	}

	for i, p := range a.Path {
		if p != b.Path[i] {
			t.Logf("found mismatch in item paths: %s:%s", p, b.Path[i])
			t.Fail()
			return
		}
	}

	if a.Name != b.Name {
		t.Logf("expected names to be the same: %s:%s", a.Name, b.Name)
		t.Fail()
		return
	}

	if a.Type != b.Type {
		t.Logf("expected types to be the same: %s:%s", a.Type, b.Type)
		t.Fail()
		return
	}

	if a.Description != b.Description {
		t.Logf("expected descriptions to be the same: %s:%s", a.Description, b.Description)
		t.Fail()
		return
	}
}
