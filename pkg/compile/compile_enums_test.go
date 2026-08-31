package compile_test

import (
	"testing"

	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/plugin-morphe-seed-data/pkg/compile"
	"github.com/kalo-build/plugin-morphe-seed-data/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-seed-data/pkg/compile/hook"
)

type CompileEnumsTestSuite struct {
	suite.Suite
}

func TestCompileEnumsTestSuite(t *testing.T) {
	suite.Run(t, new(CompileEnumsTestSuite))
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToSeedSQL_Basic() {
	seedConfig := cfg.DefaultSeedConfig()
	enum := yaml.Enum{
		Name: "Status",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"ACTIVE":   "Active",
			"INACTIVE": "Inactive",
		},
	}

	sql, err := compile.MorpheEnumToSeedSQL(seedConfig, hook.CompileMorpheEnum{}, enum)

	suite.NoError(err)
	suite.Contains(sql, "-- Seed data for statuses")
	suite.Contains(sql, "INSERT INTO statuses (id, \"key\", value) VALUES (1, 'ACTIVE', 'Active');")
	suite.Contains(sql, "INSERT INTO statuses (id, \"key\", value) VALUES (2, 'INACTIVE', 'Inactive');")
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToSeedSQL_WithSchema() {
	seedConfig := cfg.DefaultSeedConfig()
	seedConfig.Schema = "app"
	enum := yaml.Enum{
		Name: "Priority",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"LOW":  "Low",
			"HIGH": "High",
		},
	}

	sql, err := compile.MorpheEnumToSeedSQL(seedConfig, hook.CompileMorpheEnum{}, enum)

	suite.NoError(err)
	suite.Contains(sql, "INSERT INTO app.priorities")
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToSeedSQL_IntegerEnum() {
	seedConfig := cfg.DefaultSeedConfig()
	enum := yaml.Enum{
		Name: "Level",
		Type: yaml.EnumTypeInteger,
		Entries: map[string]any{
			"LOW":    1,
			"MEDIUM": 2,
			"HIGH":   3,
		},
	}

	sql, err := compile.MorpheEnumToSeedSQL(seedConfig, hook.CompileMorpheEnum{}, enum)

	suite.NoError(err)
	suite.Contains(sql, "INSERT INTO levels")
	suite.Contains(sql, "VALUES (1, 'HIGH', 3)")
	suite.Contains(sql, "VALUES (2, 'LOW', 1)")
	suite.Contains(sql, "VALUES (3, 'MEDIUM', 2)")
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToSeedSQL_SortedByKey() {
	seedConfig := cfg.DefaultSeedConfig()
	enum := yaml.Enum{
		Name: "Color",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Z_ZEBRA":  "Zebra",
			"A_APPLE":  "Apple",
			"M_MANGO":  "Mango",
		},
	}

	sql, err := compile.MorpheEnumToSeedSQL(seedConfig, hook.CompileMorpheEnum{}, enum)

	suite.NoError(err)
	suite.Contains(sql, "VALUES (1, 'A_APPLE', 'Apple')")
	suite.Contains(sql, "VALUES (2, 'M_MANGO', 'Mango')")
	suite.Contains(sql, "VALUES (3, 'Z_ZEBRA', 'Zebra')")
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToSeedSQL_WithHooks() {
	seedConfig := cfg.DefaultSeedConfig()
	enum := yaml.Enum{
		Name: "Role",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"ADMIN": "Admin",
			"USER":  "User",
		},
	}

	startCalled := false
	successCalled := false

	hooks := hook.CompileMorpheEnum{
		OnCompileMorpheEnumStart: func(c cfg.SeedConfig, e yaml.Enum) (cfg.SeedConfig, yaml.Enum, error) {
			startCalled = true
			return c, e, nil
		},
		OnCompileMorpheEnumSuccess: func(sql string) (string, error) {
			successCalled = true
			return sql, nil
		},
	}

	_, err := compile.MorpheEnumToSeedSQL(seedConfig, hooks, enum)

	suite.NoError(err)
	suite.True(startCalled)
	suite.True(successCalled)
}
