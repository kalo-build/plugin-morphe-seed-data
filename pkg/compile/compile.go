package compile

import (
	"sort"

	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/morphe-go/pkg/registry"
)

func MorpheToSeedSQL(config MorpheCompileConfig) error {
	r, rErr := registry.LoadMorpheRegistry(config.RegistryHooks, config.MorpheLoadRegistryConfig)
	if rErr != nil {
		return rErr
	}

	order := 1

	enumRowCounts := map[string]int{}

	if r.HasEnums() {
		allEnumSQL, compileErr := AllMorpheEnumsToSeedSQL(config, r)
		if compileErr != nil {
			return compileErr
		}

		enumNames := core.MapKeysSorted(allEnumSQL)
		sort.Strings(enumNames)
		for _, enumName := range enumNames {
			sql := allEnumSQL[enumName]
			tableName := GetTableNameFromModel(enumName)
			_, writeErr := config.SeedWriter.WriteSeedFile(tableName, sql, order)
			if writeErr != nil {
				return writeErr
			}

			allEnums := r.GetAllEnums()
			if enum, ok := allEnums[enumName]; ok {
				enumRowCounts[enumName] = len(enum.Entries)
			}
			order++
		}
	}

	if r.HasModels() {
		sortedModels, sortErr := SortModelsByDependency(r)
		if sortErr != nil {
			return sortErr
		}

		for _, model := range sortedModels {
			sql, compileErr := MorpheModelToSeedSQL(config.SeedConfig, config.ModelHooks, r, model, enumRowCounts)
			if compileErr != nil {
				return compileErr
			}

			tableName := GetTableNameFromModel(model.Name)
			_, writeErr := config.SeedWriter.WriteSeedFile(tableName, sql, order)
			if writeErr != nil {
				return writeErr
			}
			order++
		}
	}

	return nil
}
