package hook

import (
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-seed-data/pkg/compile/cfg"
)

type CompileMorpheEnum struct {
	OnCompileMorpheEnumStart   func(cfg.SeedConfig, yaml.Enum) (cfg.SeedConfig, yaml.Enum, error)
	OnCompileMorpheEnumSuccess func(string) (string, error)
	OnCompileMorpheEnumFailure func(cfg.SeedConfig, yaml.Enum, error) error
}
