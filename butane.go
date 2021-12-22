// Package bbcos provides a builder API for constructing Butane configs.
package bbcos

import (
	"io"

	"github.com/z5labs/bbcos/assert"
	"github.com/z5labs/bbcos/ignition"
	"github.com/z5labs/bbcos/passwd"
	"github.com/z5labs/bbcos/storage"
	"github.com/z5labs/bbcos/systemd"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

const cfgFile = "butane.yaml"

// Variant
type Variant string

const (
	FCOS      Variant = "fcos"
	OPENSHIFT Variant = "openshift"
	RHCOS     Variant = "rhcos"
)

// Version
type Version string

const (
	FCOS_1_0_0      Version = "1.0.0"
	FCOS_1_1_0      Version = "1.1.0"
	FCOS_1_2_0      Version = "1.2.0"
	FCOS_1_3_0      Version = "1.3.0"
	FCOS_1_4_0      Version = "1.4.0"
	OPENSHIFT_4_8_0 Version = "4.8.0"
	OPENSHIFT_4_9_0 Version = "4.9.0"
	RHCOS_0_1_0     Version = "0.1.0"
)

type ButaneConfig struct {
	v *viper.Viper
}

func NewButaneConfig(variant Variant, version Version, opts ...ButaneOption) ButaneConfig {
	cfg := ButaneConfig{
		v: viper.New(),
	}
	cfg.v.Set("version", string(version))
	cfg.v.Set("variant", string(variant))

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

func (cfg ButaneConfig) WriteTo(w io.Writer) (int64, error) {
	// Write config to an in-mem file since Viper doesn't currently
	// support writing to an io.Writer directly
	fs := afero.NewMemMapFs()
	cfg.v.SetFs(fs)
	cfg.v.SetConfigFile(cfgFile)
	err := cfg.v.WriteConfig()
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

type ButaneOption func(ButaneConfig)

func WithIgnition(i ignition.Ignition) ButaneOption {
	return func(cfg ButaneConfig) {
		assert.KeyNotSet(cfg.v, "ignition")
		cfg.v.Set("ignition", i.AllSettings())
	}
}

func WithStorage(s storage.Storage) ButaneOption {
	return func(cfg ButaneConfig) {
		assert.KeyNotSet(cfg.v, "storage")
		cfg.v.Set("storage", s.AllSettings())
	}
}

func WithSystemd(s systemd.Systemd) ButaneOption {
	return func(cfg ButaneConfig) {
		assert.KeyNotSet(cfg.v, "systemd")
		cfg.v.Set("systemd", s.AllSettings())
	}
}

func WithPasswd(p passwd.Passwd) ButaneOption {
	return func(cfg ButaneConfig) {
		assert.KeyNotSet(cfg.v, "passwd")
		cfg.v.Set("passwd", p.AllSettings())
	}
}
