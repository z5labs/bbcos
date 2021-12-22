package group

import (
	"github.com/z5labs/bbcos/assert"

	"github.com/spf13/viper"
)

type Group struct {
	v *viper.Viper
}

func New(name string, opts ...Option) Group {
	group := Group{
		v: viper.New(),
	}
	group.v.Set("name", name)

	for _, opt := range opts {
		opt(group)
	}

	return group
}

func (x Group) AllSettings() map[string]interface{} {
	return x.v.AllSettings()
}

type Option func(Group)

func WithName(name string) Option {
	return func(g Group) {
		assert.KeyNotSet(g.v, "name")
		g.v.Set("name", name)
	}
}

func WithGroupID(gid int) Option {
	return func(g Group) {
		assert.KeyNotSet(g.v, "gid")
		g.v.Set("gid", gid)
	}
}

func WithPasswordHash(hash string) Option {
	return func(g Group) {
		assert.KeyNotSet(g.v, "password_hash")
		g.v.Set("password_hash", hash)
	}
}

func ShouldBeSystemGroup(b bool) Option {
	return func(g Group) {
		assert.KeyNotSet(g.v, "system")
		g.v.Set("system", b)
	}
}
