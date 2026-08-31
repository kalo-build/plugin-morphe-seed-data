package compile

import (
	"path"

	r "github.com/kalo-build/morphe-go/pkg/registry"
	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-seed-data/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-seed-data/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-seed-data/pkg/compile/write"
)

type MorpheCompileConfig struct {
	rcfg.MorpheLoadRegistryConfig
	cfg.SeedConfig

	RegistryHooks r.LoadMorpheRegistryHooks

	SeedWriter write.SeedSQLWriter

	ModelHooks hook.CompileMorpheModel
	EnumHooks  hook.CompileMorpheEnum
}

func DefaultMorpheCompileConfig(
	yamlRegistryPath string,
	baseOutputDirPath string,
) MorpheCompileConfig {
	return MorpheCompileConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryEnumsDirPath:      path.Join(yamlRegistryPath, "enums"),
			RegistryModelsDirPath:     path.Join(yamlRegistryPath, "models"),
			RegistryStructuresDirPath: path.Join(yamlRegistryPath, "structures"),
			RegistryEntitiesDirPath:   path.Join(yamlRegistryPath, "entities"),
		},
		SeedConfig: cfg.DefaultSeedConfig(),

		RegistryHooks: r.LoadMorpheRegistryHooks{},

		SeedWriter: &SeedFileWriter{
			TargetDirPath: baseOutputDirPath,
		},

		ModelHooks: hook.CompileMorpheModel{},
		EnumHooks:  hook.CompileMorpheEnum{},
	}
}
