package passwd

import (
	"github.com/z5labs/bbcos/assert"
	"github.com/z5labs/bbcos/passwd/group"
	"github.com/z5labs/bbcos/passwd/user"

	"github.com/spf13/viper"
)

// Passwd
type Passwd struct {
	v *viper.Viper
}

// New
func New(opts ...Option) Passwd {
	pwd := Passwd{
		v: viper.New(),
	}

	for _, opt := range opts {
		opt(pwd)
	}

	return pwd
}

func (x Passwd) AllSettings() map[string]interface{} {
	return x.v.AllSettings()
}

type Option func(Passwd)

// WithUsers allows you to add users. All users must have a unique name.
func WithUsers(users ...user.User) Option {
	return func(pwd Passwd) {
		assert.KeyNotSet(pwd.v, "users")

		userSettings := make([]map[string]interface{}, 0, len(users))
		for _, u := range users {
			userSettings = append(userSettings, u.AllSettings())
		}

		pwd.v.Set("users", userSettings)
	}
}

// WithGroups allows you to add groups. All groups must have a unique name.
func WithGroups(groups ...group.Group) Option {
	return func(pwd Passwd) {
		assert.KeyNotSet(pwd.v, "groups")

		groupSettings := make([]map[string]interface{}, 0, len(groups))
		for _, g := range groups {
			groupSettings = append(groupSettings, g.AllSettings())
		}

		pwd.v.Set("groups", groupSettings)
	}
}
