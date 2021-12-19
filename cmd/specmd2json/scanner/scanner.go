package scanner

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/z5labs/bbcos/cmd/specmd2json/types"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type Interface interface {
	NextItem() Item
}

type Item struct {
	Path []string

	Name        string
	Type        types.Type
	Description string
}

type scanr struct {
	src  []byte
	path []string

	items chan Item
}

func New(r io.Reader) Interface {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	if rc, ok := r.(io.ReadCloser); ok {
		if err = rc.Close(); err != nil {
			panic(err)
		}
	}

	mdParser := goldmark.DefaultParser()
	specMd := mdParser.Parse(text.NewReader(b))
	if !specMd.HasChildren() {
		panic("specmd2json: expecting root not to contain children")
	}

	s := &scanr{
		src:   b,
		path:  make([]string, 0, 8),
		items: make(chan Item),
	}

	go s.scan(specMd, s.scanDoc)

	return s
}

func (s *scanr) NextItem() Item {
	return <-s.items
}

type stateFn func(ast.Node, bool) (ast.WalkStatus, stateFn)

func (s *scanr) scan(node ast.Node, action stateFn) {
	defer s.recover()

	status := ast.WalkContinue
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		status, action = action(n, entering)

		return status, nil
	})
}

func (s *scanr) recover() {
	rerr := recover()
	if rerr == nil {
		s.emit("eof", types.EOF, "eof")
		close(s.items)
		return
	}

	if err, ok := rerr.(error); ok {
		s.items <- Item{
			Path:        s.path[:],
			Name:        "error",
			Type:        types.Err,
			Description: err.Error(),
		}
		close(s.items)
		return
	}
	if item, ok := rerr.(Item); ok {
		s.items <- item
		close(s.items)
		return
	}
	close(s.items)
}

func (s *scanr) emit(name string, typ types.Type, descr string) {
	s.items <- Item{
		Path:        s.path[:],
		Name:        name,
		Type:        typ,
		Description: descr,
	}
}

func (s *scanr) pushPath(path string) {
	s.path = append(s.path, path)
}

func (s *scanr) popPath() {
	s.path = s.path[:len(s.path)-1]
}

func (s *scanr) errorf(format string, args ...interface{}) {
	panic(Item{
		Path:        s.path[:],
		Name:        "",
		Type:        types.Err,
		Description: fmt.Sprintf(format, args...),
	})
}

func (s *scanr) scanDoc(node ast.Node, entering bool) (ast.WalkStatus, stateFn) {
	switch node.(type) {
	case *ast.List:
		if !entering {
			return ast.WalkContinue, s.scanDoc
		}
		return ast.WalkContinue, s.scanListItem
	case *ast.ListItem:
		return s.scanListItem(node, entering)
	default:
		return ast.WalkContinue, s.scanDoc
	}
}

func (s *scanr) scanListItem(node ast.Node, entering bool) (ast.WalkStatus, stateFn) {
	item, ok := node.(*ast.ListItem)
	if !ok {
		s.errorf("parseListItem: given node must be *ast.ListItem")
	}

	text := string(item.FirstChild().Text(s.src))
	name, ty, descr := s.scanItemText(text)
	if !entering {
		if types.InnerType(ty) == types.Object {
			s.popPath()
		}
		return ast.WalkContinue, s.scanDoc
	}

	s.emit(name, ty, descr)
	if types.InnerType(ty) != types.Object {
		return ast.WalkSkipChildren, s.scanDoc
	}

	s.pushPath(name)
	return ast.WalkContinue, s.scanDoc
}

func (s *scanr) scanItemText(text string) (key string, ty types.Type, description string) {
	lParen := strings.Index(text, "(")
	rParen := strings.Index(text, ")")
	if lParen == -1 || rParen == -1 {
		s.errorf("scanItemText: missing complete pair of parens around field type")
	}

	key = strings.TrimSpace(text[:lParen])
	ty = types.Type(strings.TrimSuffix(strings.TrimSpace(text[lParen+1:rParen]), "s"))
	description = strings.TrimPrefix(text[rParen+1:], ": ")

	return
}
