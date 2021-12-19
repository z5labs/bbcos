package parser

import (
	"io"
	"strings"
	"sync"

	"github.com/z5labs/bbcos/cmd/specmd2json/keys"
	"github.com/z5labs/bbcos/cmd/specmd2json/scanner"
	"github.com/z5labs/bbcos/cmd/specmd2json/types"

	"github.com/spf13/viper"
)

const (
	delim = "."

	schemaVersion = "https://json-schema.org/draft/2020-12/schema"
)

type stateFn func(*parser, *viper.Viper) stateFn

type parser struct {
	scanner.Interface

	arrs []bool
	vals []string
}

func markArrayStatus(arrs []bool, path []string, isArr bool) []bool {
	if len(path) == len(arrs)-1 {
		arrs[len(path)] = isArr
		return arrs
	}
	return append(arrs[:len(arrs)], isArr)
}

func Parse(r io.Reader) *viper.Viper {
	cfg := viper.New()
	p := &parser{
		Interface: scanner.New(r),
		arrs:      make([]bool, 0, 5),
		vals:      make([]string, 0, 4),
	}

	for action := parseDoc(p, cfg); action != nil; {
		action = action(p, cfg)
	}

	return cfg
}

func parseDoc(p *parser, cfg *viper.Viper) stateFn {
	item := p.NextItem()
	if item.Type == types.EOF {
		return nil
	}
	if item.Type == types.Err {
		panic(buildPathWith(bufPool, item.Path, p.arrs, item.Name, "error") + ": " + item.Description)
	}
	if len(item.Path) < len(p.arrs) {
		p.arrs = p.arrs[:len(item.Path)]
	}

	prevIsArr := false
	if len(p.arrs) > 0 {
		prevIsArr = p.arrs[len(p.arrs)-1]
	}

	p.vals = p.vals[:0]
	if prevIsArr {
		p.vals = append(p.vals, "items")
	}
	p.vals = append(p.vals, "properties")
	p.vals = append(p.vals, item.Name)

	p.set(cfg, item.Path, types.Normalize(item.Type), keys.Type.String())
	p.set(cfg, item.Path, item.Description, keys.Description.String())

	isArr := types.IsArray(item.Type)
	if isArr {
		p.set(cfg, item.Path, types.InnerType(item.Type).String(), keys.Items.String(), keys.Type.String())
	}
	p.arrs = markArrayStatus(p.arrs, item.Path, isArr)
	return parseDoc
}

func (p *parser) set(cfg *viper.Viper, path []string, val string, keyPath ...string) {
	p.vals = append(p.vals, keyPath...)
	key := buildPathWith(bufPool, path, p.arrs, p.vals...)
	p.vals = p.vals[:len(p.vals)-len(keyPath)]

	cfg.Set(key, val)
}

var bufPool = &sync.Pool{
	New: func() interface{} {
		return new(strings.Builder)
	},
}

type pool interface {
	Get() interface{}
	Put(interface{})
}

func buildPathWith(bufPool pool, path []string, arrs []bool, vals ...string) string {
	if len(path) != len(arrs) {
		panic("parser: expected path and arrs to have same length")
	}

	buf := bufPool.Get().(*strings.Builder)
	defer buf.Reset()
	defer bufPool.Put(buf)

	if len(path) > 0 {
		buf.WriteString("properties")
		buf.WriteString(delim)
	}

	for i, p := range path {
		buf.WriteString(p)
		buf.WriteString(delim)
		if i == len(path)-1 {
			break
		}
		if arrs[i] {
			buf.WriteString("items")
			buf.WriteString(delim)
		}
		buf.WriteString("properties")
		buf.WriteString(delim)
	}
	for i, val := range vals {
		buf.WriteString(val)
		if i != len(vals)-1 {
			buf.WriteString(delim)
		}
	}
	return buf.String()
}
