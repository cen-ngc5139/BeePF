package loader

import (
	"fmt"
	"time"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/skeleton/export"
)

func ValidateAndMutateConfig(cfg *Config) error {
	if cfg.ObjectBytes == nil && cfg.ObjectPath == "" {
		return fmt.Errorf("object file is required")
	}

	if cfg.ObjectBytes != nil && cfg.ObjectPath != "" {
		return fmt.Errorf("object file and object bytes cannot both be set")
	}

	if cfg.Logger == nil {
		return fmt.Errorf("logger is required")
	}

	if cfg.PollTimeout == 0 {
		cfg.PollTimeout = 1 * time.Second
	}

	if cfg.Properties.Maps == nil || cfg.Properties.EventHandler == nil {
		cfg.Properties.EventHandler = &export.MyCustomHandler{Logger: cfg.Logger}
	}

	if cfg.Properties.Stats != nil {
		if cfg.Properties.Stats.Interval == 0 {
			cfg.Properties.Stats.Interval = 1 * time.Second
		}
	}

	return nil
}
