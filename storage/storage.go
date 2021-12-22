package storage

import (
	"github.com/spf13/viper"
)

// Storage provides a classic builder pattern for describing the
// desired state of the system’s storage devices.
//
type Storage struct {
	v *viper.Viper
}

// New creates a new Storager.
func New(opts ...Option) Storage {
	storage := Storage{
		v: viper.New(),
	}

	for _, opt := range opts {
		opt(storage)
	}

	return storage
}

func (x Storage) AllSettings() map[string]interface{} {
	return x.v.AllSettings()
}

type Option func(Storage)
