package compile_test

import (
	"strings"
	"testing"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/plugin-morphe-seed-data/pkg/compile"
	"github.com/kalo-build/plugin-morphe-seed-data/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-seed-data/pkg/compile/hook"
)

type CompileModelsTestSuite struct {
	suite.Suite
}

func TestCompileModelsTestSuite(t *testing.T) {
	suite.Run(t, new(CompileModelsTestSuite))
}

func (suite *CompileModelsTestSuite) newMinimalRegistry() *registry.Registry {
	r := &registry.Registry{}
	r.SetModel("Item", yaml.Model{
		Name: "Item",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: yaml.ModelFieldTypeAutoIncrement},
			"Name": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	})
	return r
}

func (suite *CompileModelsTestSuite) TestMorpheModelToSeedSQL_BasicFields() {
	r := suite.newMinimalRegistry()
	seedConfig := cfg.SeedConfig{RowCount: 2, Seed: 42}
	model, _ := r.GetModel("Item")

	sql, err := compile.MorpheModelToSeedSQL(seedConfig, hook.CompileMorpheModel{}, r, model, nil)

	suite.NoError(err)
	suite.Contains(sql, "-- Seed data for items")
	suite.Contains(sql, "INSERT INTO items")
	insertCount := 0
	for _, line := range strings.Split(sql, "\n") {
		if strings.HasPrefix(line, "INSERT") {
			insertCount++
		}
	}
	suite.Equal(2, insertCount)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToSeedSQL_Deterministic() {
	r := suite.newMinimalRegistry()
	seedConfig := cfg.SeedConfig{RowCount: 3, Seed: 42}
	model, _ := r.GetModel("Item")

	sql1, err1 := compile.MorpheModelToSeedSQL(seedConfig, hook.CompileMorpheModel{}, r, model, nil)
	sql2, err2 := compile.MorpheModelToSeedSQL(seedConfig, hook.CompileMorpheModel{}, r, model, nil)

	suite.NoError(err1)
	suite.NoError(err2)
	suite.Equal(sql1, sql2)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToSeedSQL_WithSchema() {
	r := suite.newMinimalRegistry()
	seedConfig := cfg.SeedConfig{RowCount: 1, Seed: 42, Schema: "app"}
	model, _ := r.GetModel("Item")

	sql, err := compile.MorpheModelToSeedSQL(seedConfig, hook.CompileMorpheModel{}, r, model, nil)

	suite.NoError(err)
	suite.Contains(sql, "INSERT INTO app.items")
}

func (suite *CompileModelsTestSuite) TestMorpheModelToSeedSQL_AutoIncrementSequential() {
	r := suite.newMinimalRegistry()
	seedConfig := cfg.SeedConfig{RowCount: 3, Seed: 42}
	model, _ := r.GetModel("Item")

	sql, err := compile.MorpheModelToSeedSQL(seedConfig, hook.CompileMorpheModel{}, r, model, nil)

	suite.NoError(err)
	var insertLines []string
	for _, line := range strings.Split(sql, "\n") {
		if strings.HasPrefix(line, "INSERT") {
			insertLines = append(insertLines, line)
		}
	}
	suite.Len(insertLines, 3)
	suite.Contains(insertLines[0], "1,")
	suite.Contains(insertLines[1], "2,")
	suite.Contains(insertLines[2], "3,")
}

func (suite *CompileModelsTestSuite) TestMorpheModelToSeedSQL_AllFieldTypes() {
	r := &registry.Registry{}
	model := yaml.Model{
		Name: "FullModel",
		Fields: map[string]yaml.ModelField{
			"ID":        {Type: yaml.ModelFieldTypeUUID},
			"Label":     {Type: yaml.ModelFieldTypeString},
			"Count":     {Type: yaml.ModelFieldTypeInteger},
			"Price":     {Type: yaml.ModelFieldTypeFloat},
			"Active":    {Type: yaml.ModelFieldTypeBoolean},
			"CreatedAt": {Type: yaml.ModelFieldTypeTime},
			"BirthDate": {Type: yaml.ModelFieldTypeDate},
			"Secret":    {Type: yaml.ModelFieldTypeProtected},
			"Hash":      {Type: yaml.ModelFieldTypeSealed},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	}
	r.SetModel("FullModel", model)

	seedConfig := cfg.SeedConfig{RowCount: 1, Seed: 42}
	sql, err := compile.MorpheModelToSeedSQL(seedConfig, hook.CompileMorpheModel{}, r, model, nil)

	suite.NoError(err)
	suite.Contains(sql, "INSERT INTO full_models")
	suite.Contains(sql, "active")
	suite.Contains(sql, "birth_date")
	suite.Contains(sql, "created_at")
}

func (suite *CompileModelsTestSuite) TestMorpheModelToSeedSQL_OptionalField() {
	r := &registry.Registry{}
	model := yaml.Model{
		Name: "Profile",
		Fields: map[string]yaml.ModelField{
			"ID":  {Type: yaml.ModelFieldTypeAutoIncrement},
			"Bio": {Type: yaml.ModelFieldTypeString, Attributes: []string{"optional"}},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	}
	r.SetModel("Profile", model)

	seedConfig := cfg.SeedConfig{RowCount: 20, Seed: 42}
	sql, err := compile.MorpheModelToSeedSQL(seedConfig, hook.CompileMorpheModel{}, r, model, nil)

	suite.NoError(err)
	suite.Contains(sql, "NULL")
}

func (suite *CompileModelsTestSuite) TestMorpheModelToSeedSQL_WithFormatAttributes() {
	r := &registry.Registry{}
	model := yaml.Model{
		Name: "Contact",
		Fields: map[string]yaml.ModelField{
			"ID":    {Type: yaml.ModelFieldTypeAutoIncrement},
			"Email": {Type: yaml.ModelFieldTypeString, Attributes: []string{"format=email"}},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	}
	r.SetModel("Contact", model)

	seedConfig := cfg.SeedConfig{RowCount: 3, Seed: 42}
	sql, err := compile.MorpheModelToSeedSQL(seedConfig, hook.CompileMorpheModel{}, r, model, nil)

	suite.NoError(err)
	suite.Contains(sql, "@")
}

func (suite *CompileModelsTestSuite) TestMorpheModelToSeedSQL_EnumFK() {
	r := &registry.Registry{}
	r.SetEnum("Status", yaml.Enum{
		Name: "Status",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"ACTIVE":   "Active",
			"INACTIVE": "Inactive",
		},
	})
	model := yaml.Model{
		Name: "Task",
		Fields: map[string]yaml.ModelField{
			"ID":     {Type: yaml.ModelFieldTypeAutoIncrement},
			"Status": {Type: "Status"},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	}
	r.SetModel("Task", model)

	enumRowCounts := map[string]int{"Status": 2}
	seedConfig := cfg.SeedConfig{RowCount: 5, Seed: 42}
	sql, err := compile.MorpheModelToSeedSQL(seedConfig, hook.CompileMorpheModel{}, r, model, enumRowCounts)

	suite.NoError(err)
	suite.Contains(sql, "status_id")
}

func (suite *CompileModelsTestSuite) TestMorpheModelToSeedSQL_ForOneRelation() {
	r := &registry.Registry{}
	r.SetModel("Company", yaml.Model{
		Name: "Company",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeAutoIncrement},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	})
	model := yaml.Model{
		Name: "Employee",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: yaml.ModelFieldTypeAutoIncrement},
			"Name": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Company": {Type: "ForOne"},
		},
	}
	r.SetModel("Employee", model)

	seedConfig := cfg.SeedConfig{RowCount: 2, Seed: 42}
	sql, err := compile.MorpheModelToSeedSQL(seedConfig, hook.CompileMorpheModel{}, r, model, nil)

	suite.NoError(err)
	suite.Contains(sql, "company_id")
}

func (suite *CompileModelsTestSuite) TestMorpheModelToSeedSQL_AliasedRelation() {
	r := &registry.Registry{}
	r.SetModel("Contact", yaml.Model{
		Name: "Contact",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeAutoIncrement},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	})
	model := yaml.Model{
		Name: "Person",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeAutoIncrement},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"WorkContact":     {Type: "ForOne", Aliased: "Contact"},
			"PersonalContact": {Type: "ForOne", Aliased: "Contact"},
		},
	}
	r.SetModel("Person", model)

	seedConfig := cfg.SeedConfig{RowCount: 2, Seed: 42}
	sql, err := compile.MorpheModelToSeedSQL(seedConfig, hook.CompileMorpheModel{}, r, model, nil)

	suite.NoError(err)
	suite.Contains(sql, "personal_contact_id")
	suite.Contains(sql, "work_contact_id")
}

func (suite *CompileModelsTestSuite) TestMorpheModelToSeedSQL_InvalidRowCount() {
	r := suite.newMinimalRegistry()
	seedConfig := cfg.SeedConfig{RowCount: 0, Seed: 42}
	model, _ := r.GetModel("Item")

	_, err := compile.MorpheModelToSeedSQL(seedConfig, hook.CompileMorpheModel{}, r, model, nil)

	suite.Error(err)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToSeedSQL_WithHooks() {
	r := suite.newMinimalRegistry()
	seedConfig := cfg.SeedConfig{RowCount: 1, Seed: 42}
	model, _ := r.GetModel("Item")

	startCalled := false
	successCalled := false

	hooks := hook.CompileMorpheModel{
		OnCompileMorpheModelStart: func(c cfg.SeedConfig, m yaml.Model) (cfg.SeedConfig, yaml.Model, error) {
			startCalled = true
			return c, m, nil
		},
		OnCompileMorpheModelSuccess: func(sql string) (string, error) {
			successCalled = true
			return sql, nil
		},
	}

	_, err := compile.MorpheModelToSeedSQL(seedConfig, hooks, r, model, nil)

	suite.NoError(err)
	suite.True(startCalled)
	suite.True(successCalled)
}

func (suite *CompileModelsTestSuite) TestGetModelDependencies() {
	r := &registry.Registry{}
	r.SetModel("Company", yaml.Model{
		Name: "Company",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeAutoIncrement},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	})

	model := yaml.Model{
		Name: "Person",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeAutoIncrement},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Company": {Type: "ForOne"},
			"Friend":  {Type: "HasMany"},
		},
	}

	deps := compile.GetModelDependencies(r, model)

	suite.Contains(deps, "Company")
	suite.NotContains(deps, "Friend")
}

func (suite *CompileModelsTestSuite) TestSortModelsByDependency() {
	r := &registry.Registry{}
	r.SetModel("Company", yaml.Model{
		Name: "Company",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeAutoIncrement},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	})
	r.SetModel("Person", yaml.Model{
		Name: "Person",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeAutoIncrement},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Company": {Type: "ForOne"},
		},
	})

	sorted, err := compile.SortModelsByDependency(r)

	suite.NoError(err)
	suite.Len(sorted, 2)
	suite.Equal("Company", sorted[0].Name)
	suite.Equal("Person", sorted[1].Name)
}
