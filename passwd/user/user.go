package user

import (
	"github.com/z5labs/bbcos/assert"

	"github.com/spf13/viper"
)

type User struct {
	v *viper.Viper
}

func New(name string, opts ...Option) User {
	user := User{
		v: viper.New(),
	}
	user.v.Set("name", name)

	for _, opt := range opts {
		opt(user)
	}

	return user
}

func (x User) AllSettings() map[string]interface{} {
	return x.v.AllSettings()
}

type Option func(User)

func WithPasswordHash(hash string) Option {
	return func(user User) {
		assert.KeyNotSet(user.v, "password_hash")
		user.v.Set("password_hash", hash)
	}
}

func WithSSHAuthorizedKeys(keys []string) Option {
	return func(user User) {
		assert.KeyNotSet(user.v, "ssh_authorized_keys")
		user.v.Set("ssh_authorized_keys", keys)
	}
}

func WithUserID(id int) Option {
	return func(user User) {
		assert.KeyNotSet(user.v, "uid")
		user.v.Set("uid", id)
	}
}

func WithGECOS(gecos string) Option {
	return func(user User) {
		assert.KeyNotSet(user.v, "gecos")
		user.v.Set("gecos", gecos)
	}
}

func WithHomeDirectory(dir string) Option {
	return func(user User) {
		assert.KeyNotSet(user.v, "home_dir")
		user.v.Set("home_dir", dir)
	}
}

func ShouldCreateUserHomeDir(b bool) Option {
	return func(user User) {
		assert.KeyNotSet(user.v, "no_create_home")
		user.v.Set("no_create_home", b)
	}
}

func WithPrimaryGroup(name string) Option {
	return func(user User) {
		assert.KeyNotSet(user.v, "primary_group")
		user.v.Set("primary_group", name)
	}
}

func WithSupplementaryGroups(names ...string) Option {
	return func(user User) {
		assert.KeyNotSet(user.v, "groups")
		user.v.Set("groups", names)
	}
}

func ShouldCreateUserGroup(b bool) Option {
	return func(user User) {
		assert.KeyNotSet(user.v, "no_user_group")
		user.v.Set("no_user_group", b)
	}
}

func ShouldAddToLastAndFailLogDBs(b bool) Option {
	return func(user User) {
		assert.KeyNotSet(user.v, "no_log_init")
		user.v.Set("no_log_init", b)
	}
}

func WithShell(name string) Option {
	return func(user User) {
		assert.KeyNotSet(user.v, "shell")
		user.v.Set("shell", name)
	}
}

func ShouldBeSystemAccount(b bool) Option {
	return func(user User) {
		assert.KeyNotSet(user.v, "system")
		user.v.Set("system", b)
	}
}
