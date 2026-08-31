package cfg

import "fmt"

type SeedConfig struct {
	RowCount int
	Seed     int64
	Schema   string
}

func (c SeedConfig) Validate() error {
	if c.RowCount < 1 {
		return fmt.Errorf("rowCount must be at least 1")
	}
	return nil
}

func DefaultSeedConfig() SeedConfig {
	return SeedConfig{
		RowCount: 5,
		Seed:     42,
	}
}
