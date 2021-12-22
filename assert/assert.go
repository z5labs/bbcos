package assert

import (
	"fmt"

	"github.com/spf13/viper"
)

func KeyNotSet(v *viper.Viper, key string) {
	if hasKey := v.Get(key); hasKey != nil {
		panicKeyAlreadySet(key)
	}
}

func panicKeyAlreadySet(key string) {
	panic(fmt.Sprintf("%s has already been set", key))
}
