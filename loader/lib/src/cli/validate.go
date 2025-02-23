package loader

import (
	"fmt"
	"time"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/skeleton/export"
)

func ValidateAndMutateConfig(cfg *Config) error {
	if cfg.ObjectPath == "" {
		return fmt.Errorf("object path is required")
	}

	if cfg.StructName == "" {
		return fmt.Errorf("struct name is required")
	}

	if cfg.UserExporterHandler == nil {
		cfg.UserExporterHandler = &export.MyCustomHandler{
			Logger: cfg.Logger,
		}
	}

	if cfg.Logger == nil {
		return fmt.Errorf("logger is required")
	}

	if cfg.PollTimeout == 0 {
		cfg.PollTimeout = 1 * time.Second
	}

	if cfg.IsEnableStats {
		if cfg.StatsInterval == 0 {
			cfg.StatsInterval = 1 * time.Second
		}
	}

	return nil
}
