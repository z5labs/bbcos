package ignition

import (
	"github.com/spf13/viper"
)

// Ignition metadata about the configuration itself.
type Ignition struct {
	v *viper.Viper
}

// New creates a new Ingition.
func New(opts ...Option) Ignition {
	ign := Ignition{
		v: viper.New(),
	}

	for _, opt := range opts {
		opt(ign)
	}

	return ign
}

func (x Ignition) AllSettings() map[string]interface{} {
	return x.v.AllSettings()
}

type Option func(Ignition)
