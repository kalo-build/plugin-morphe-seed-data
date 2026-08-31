package compile

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-seed-data/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-seed-data/pkg/compile/hook"
)

func AllMorpheEnumsToSeedSQL(config MorpheCompileConfig, r *registry.Registry) (map[string]string, error) {
	allEnumSQL := map[string]string{}
	for enumName, enum := range r.GetAllEnums() {
		sql, err := MorpheEnumToSeedSQL(config.SeedConfig, config.EnumHooks, enum)
		if err != nil {
			return nil, err
		}
		allEnumSQL[enumName] = sql
	}
	return allEnumSQL, nil
}

func MorpheEnumToSeedSQL(seedConfig cfg.SeedConfig, hooks hook.CompileMorpheEnum, enum yaml.Enum) (string, error) {
	if hooks.OnCompileMorpheEnumStart != nil {
		var startErr error
		seedConfig, enum, startErr = hooks.OnCompileMorpheEnumStart(seedConfig, enum)
		if startErr != nil {
			return "", triggerEnumFailure(hooks, seedConfig, enum, startErr)
		}
	}

	sql, err := morpheEnumToSeedSQL(seedConfig, enum)
	if err != nil {
		return "", triggerEnumFailure(hooks, seedConfig, enum, err)
	}

	if hooks.OnCompileMorpheEnumSuccess != nil {
		sql, err = hooks.OnCompileMorpheEnumSuccess(sql)
		if err != nil {
			return "", triggerEnumFailure(hooks, seedConfig, enum, err)
		}
	}
	return sql, nil
}

func morpheEnumToSeedSQL(seedConfig cfg.SeedConfig, enum yaml.Enum) (string, error) {
	tableName := GetTableNameFromModel(enum.Name)
	if seedConfig.Schema != "" {
		tableName = seedConfig.Schema + "." + tableName
	}

	entryKeys := core.MapKeysSorted(enum.Entries)
	sort.Strings(entryKeys)

	var lines []string
	lines = append(lines, fmt.Sprintf("-- Seed data for %s", tableName))

	for idx, key := range entryKeys {
		value := enum.Entries[key]
		id := idx + 1
		formattedValue := FormatSQLValue(value)
		line := fmt.Sprintf("INSERT INTO %s (id, \"key\", value) VALUES (%d, '%s', %s);",
			tableName, id, key, formattedValue)
		lines = append(lines, line)
	}

	lines = append(lines, "")
	return strings.Join(lines, "\n"), nil
}

func triggerEnumFailure(hooks hook.CompileMorpheEnum, config cfg.SeedConfig, enum yaml.Enum, err error) error {
	if hooks.OnCompileMorpheEnumFailure != nil {
		return hooks.OnCompileMorpheEnumFailure(config, enum, err)
	}
	return err
}
