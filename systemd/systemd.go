package systemd

import (
	"github.com/spf13/viper"
)

//  Systemd
type Systemd struct {
	v *viper.Viper
}

// New
func New(opts ...Option) Systemd {
	systemd := Systemd{
		v: viper.New(),
	}

	for _, opt := range opts {
		opt(systemd)
	}

	return systemd
}

func (x Systemd) AllSettings() map[string]interface{} {
	return x.v.AllSettings()
}

// Option
type Option func(Systemd)
