package hook

import (
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-seed-data/pkg/compile/cfg"
)

type CompileMorpheModel struct {
	OnCompileMorpheModelStart   func(cfg.SeedConfig, yaml.Model) (cfg.SeedConfig, yaml.Model, error)
	OnCompileMorpheModelSuccess func(string) (string, error)
	OnCompileMorpheModelFailure func(cfg.SeedConfig, yaml.Model, error) error
}
