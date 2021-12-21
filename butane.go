// Package bbcos provides a builder API for constructing Butane configs.
package bbcos

import (
	"fmt"
	"io"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

const (
	STORAGE     = "storage"
	IGNITION    = "ignition"
	SYSTEMD     = "systemd"
	PASSWD      = "passwd"
	KERNEL_ARGS = "kernel_arguments"
	BOOT_DEVICE = "boot_device"
)

// Builder provides a classic builder pattern for constructing Butane configs.
type Builder struct {
	v *viper.Viper
}

// Create a new Butane config builder
func New() *Builder {
	return &Builder{
		v: viper.New(),
	}
}

const cfgFile = "butane.yaml"

func (b *Builder) WriteTo(w io.Writer) (int64, error) {
	// Write config to an in-mem file since Viper doesn't currently
	// support writing to an io.Writer directly
	fs := afero.NewMemMapFs()
	b.v.SetFs(fs)
	b.v.SetConfigFile(cfgFile)
	err := b.v.WriteConfig()
	if err != nil {
		panic(err)
	}

	f, err := fs.Open(cfgFile)
	if err != nil {
		panic(err)
	}

	n, err := io.Copy(w, f)
	if err != nil {
		panic(err)
	}
	return n, nil
}

// Used to differentiate configs for different operating systems.
func (b *Builder) Variant(v string) *Builder {
	b.v.Set("variant", v)
	return b
}

// The semantic version of the spec used for validating the constructed config.
func (b *Builder) Version(v string) *Builder {
	b.v.Set("version", v)
	return b
}

// Igniter provides a classic builder pattern for providing metadata
// about the config itself.
//
type Igniter struct {
	v *viper.Viper
}

// WithIgnition
func (b *Builder) WithIgnition(f func(*Igniter)) *Builder {
	b.assertKeyNotSet(IGNITION)
	b.v.Set(IGNITION, map[string]interface{}{})
	f(NewIgniter(b.v.Sub(IGNITION)))
	return b
}

// NewIgniter creates a new Ingiter.
func NewIgniter(v *viper.Viper) *Igniter {
	return &Igniter{
		v: v,
	}
}

// Storager provides a classic builder pattern for describing the
// desired state of the system’s storage devices.
//
type Storager struct {
	v *viper.Viper
}

// WithStorage
func (b *Builder) WithStorage(f func(*Storager)) *Builder {
	b.assertKeyNotSet(STORAGE)
	b.v.Set(STORAGE, map[string]interface{}{})
	f(NewStorager(b.v.Sub(STORAGE)))
	return b
}

// NewStorager creates a new Storager.
func NewStorager(v *viper.Viper) *Storager {
	return &Storager{
		v: v,
	}
}

type Systemder struct {
	v *viper.Viper
}

func (b *Builder) WithSystemd(f func(*Systemder)) *Builder {
	b.assertKeyNotSet(SYSTEMD)
	b.v.Set(SYSTEMD, map[string]interface{}{})
	f(NewSystemder(b.v.Sub(SYSTEMD)))
	return b
}

func NewSystemder(v *viper.Viper) *Systemder {
	return &Systemder{
		v: v,
	}
}

type Passwder struct {
	v *viper.Viper
}

func (b *Builder) WithPasswder(f func(*Passwder)) *Builder {
	b.assertKeyNotSet(PASSWD)
	b.v.Set(PASSWD, map[string]interface{}{})
	f(NewPasswder(b.v.Sub(PASSWD)))
	return b
}

func NewPasswder(v *viper.Viper) *Passwder {
	return &Passwder{
		v: v,
	}
}

type User struct {
	v *viper.Viper
}

func (p *Passwder) WithUser(f func(*User)) *Passwder {
	users, ok := p.v.Get("users").([]map[string]interface{})
	if !ok {
		users = []map[string]interface{}{}
		p.v.Set("users", users)
	}

	v := viper.New()
	user := NewUser(v)
	f(user)

	users = append(users, v.AllSettings())
	return p
}

func NewUser(v *viper.Viper) *User {
	return &User{
		v: v,
	}
}

func (u *User) WithName(name string) *User {
	u.v.Set("name", name)
	return u
}

func (u *User) WithPasswordHash(hash string) *User {
	u.v.Set("password_hash", hash)
	return u
}

func (u *User) WithSSHAuthorizedKeys(keys []string) *User {
	u.v.Set("ssh_authorized_keys", keys)
	return u
}

type Kerneler struct {
	v *viper.Viper
}

func (b *Builder) WithKernelArgs(f func(*Kerneler)) *Builder {
	b.assertKeyNotSet(KERNEL_ARGS)
	b.v.Set(KERNEL_ARGS, map[string]interface{}{})
	f(NewKerneler(b.v.Sub(KERNEL_ARGS)))
	return b
}

func NewKerneler(v *viper.Viper) *Kerneler {
	return &Kerneler{
		v: v,
	}
}

type Booter struct {
	v *viper.Viper
}

func (b *Builder) WithBootDevice(f func(*Booter)) *Builder {
	b.assertKeyNotSet(BOOT_DEVICE)
	b.v.Set(BOOT_DEVICE, map[string]interface{}{})
	f(NewBooter(b.v.Sub(BOOT_DEVICE)))
	return b
}

func NewBooter(v *viper.Viper) *Booter {
	return &Booter{
		v: v,
	}
}

func (b *Builder) assertKeyNotSet(key string) {
	if hasKey := b.v.Get(key); hasKey != nil {
		panicKeyAlreadySet(key)
	}
}

func panicKeyAlreadySet(key string) {
	panic(fmt.Sprintf("bbcos: %s has already been configured", key))
}
