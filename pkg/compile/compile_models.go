package compile

import (
	"fmt"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/go-util/strcase"
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-seed-data/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-seed-data/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-seed-data/pkg/typemap"
)

type modelDependency struct {
	Name         string
	Dependencies []string
}

func AllMorpheModelsToSeedSQL(config MorpheCompileConfig, r *registry.Registry, enumRowCounts map[string]int) (map[string]string, error) {
	allModelSQL := map[string]string{}
	for modelName, model := range r.GetAllModels() {
		sql, err := MorpheModelToSeedSQL(config.SeedConfig, config.ModelHooks, r, model, enumRowCounts)
		if err != nil {
			return nil, err
		}
		allModelSQL[modelName] = sql
	}
	return allModelSQL, nil
}

func MorpheModelToSeedSQL(seedConfig cfg.SeedConfig, hooks hook.CompileMorpheModel, r *registry.Registry, model yaml.Model, enumRowCounts map[string]int) (string, error) {
	if hooks.OnCompileMorpheModelStart != nil {
		var startErr error
		seedConfig, model, startErr = hooks.OnCompileMorpheModelStart(seedConfig, model)
		if startErr != nil {
			return "", triggerModelFailure(hooks, seedConfig, model, startErr)
		}
	}

	sql, err := morpheModelToSeedSQL(seedConfig, r, model, enumRowCounts)
	if err != nil {
		return "", triggerModelFailure(hooks, seedConfig, model, err)
	}

	if hooks.OnCompileMorpheModelSuccess != nil {
		sql, err = hooks.OnCompileMorpheModelSuccess(sql)
		if err != nil {
			return "", triggerModelFailure(hooks, seedConfig, model, err)
		}
	}
	return sql, nil
}

func morpheModelToSeedSQL(seedConfig cfg.SeedConfig, r *registry.Registry, model yaml.Model, enumRowCounts map[string]int) (string, error) {
	validateErr := seedConfig.Validate()
	if validateErr != nil {
		return "", validateErr
	}

	validateModelErr := model.Validate(r.GetAllEnums())
	if validateModelErr != nil {
		return "", validateModelErr
	}

	tableName := GetTableNameFromModel(model.Name)
	if seedConfig.Schema != "" {
		tableName = seedConfig.Schema + "." + tableName
	}

	columns, columnFieldMeta, err := buildColumnMetadata(r, model)
	if err != nil {
		return "", err
	}

	faker := gofakeit.New(uint64(seedConfig.Seed))

	var lines []string
	lines = append(lines, fmt.Sprintf("-- Seed data for %s", tableName))

	for rowIdx := 0; rowIdx < seedConfig.RowCount; rowIdx++ {
		values := make([]string, len(columns))
		for colIdx := range columns {
			meta := columnFieldMeta[colIdx]
			val := generateColumnValue(faker, meta, rowIdx, seedConfig.RowCount, enumRowCounts)
			values[colIdx] = FormatSQLValue(val)
		}

		columnList := strings.Join(columns, ", ")
		valueList := strings.Join(values, ", ")
		line := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", tableName, columnList, valueList)
		lines = append(lines, line)
	}

	lines = append(lines, "")
	return strings.Join(lines, "\n"), nil
}

type columnMeta struct {
	FieldName  string
	FieldType  yaml.ModelFieldType
	IsOptional bool
	IsAutoInc  bool
	IsEnumFK   bool
	EnumName   string
	IsModelFK  bool
	Hints      typemap.FormatHints
}

func buildColumnMetadata(r *registry.Registry, model yaml.Model) ([]string, []columnMeta, error) {
	var columns []string
	var metas []columnMeta

	fieldNames := core.MapKeysSorted(model.Fields)
	for _, fieldName := range fieldNames {
		field := model.Fields[fieldName]
		colName := GetColumnNameFromField(fieldName)
		hints := typemap.ParseFormatHints(field.Attributes)
		isOptional := typemap.HasAttribute(field.Attributes, "optional")

		if yaml.IsModelFieldTypePrimitive(field.Type) {
			columns = append(columns, colName)
			metas = append(metas, columnMeta{
				FieldName:  fieldName,
				FieldType:  field.Type,
				IsOptional: isOptional,
				IsAutoInc:  field.Type == yaml.ModelFieldTypeAutoIncrement,
				Hints:      hints,
			})
			continue
		}

		_, enumErr := r.GetEnum(string(field.Type))
		if enumErr != nil {
			return nil, nil, fmt.Errorf("field '%s' has unsupported type '%s'", fieldName, field.Type)
		}

		colName = colName + "_id"
		columns = append(columns, colName)
		metas = append(metas, columnMeta{
			FieldName:  fieldName,
			FieldType:  field.Type,
			IsOptional: isOptional,
			IsEnumFK:   true,
			EnumName:   string(field.Type),
			Hints:      hints,
		})
	}

	relatedNames := core.MapKeysSorted(model.Related)
	for _, relatedName := range relatedNames {
		relation := model.Related[relatedName]
		relationType := strings.ToLower(relation.Type)

		if strings.HasPrefix(relationType, "foronepoly") {
			typeColName := strcase.ToSnakeCaseLower(relatedName) + "_type"
			idColName := strcase.ToSnakeCaseLower(relatedName) + "_id"
			columns = append(columns, typeColName)
			metas = append(metas, columnMeta{FieldName: relatedName, FieldType: yaml.ModelFieldTypeString})
			columns = append(columns, idColName)
			metas = append(metas, columnMeta{FieldName: relatedName, FieldType: yaml.ModelFieldTypeString})
			continue
		}

		if strings.Contains(relationType, "poly") {
			continue
		}

		if !strings.HasPrefix(relationType, "for") || !isRelationOne(relationType) {
			continue
		}

		targetModelName := relatedName
		if relation.Aliased != "" {
			targetModelName = relation.Aliased
		}

		relatedModel, modelErr := r.GetModel(targetModelName)
		if modelErr != nil {
			return nil, nil, modelErr
		}

		primaryID, hasPrimary := relatedModel.Identifiers["primary"]
		if !hasPrimary || len(primaryID.Fields) != 1 {
			return nil, nil, fmt.Errorf("related model %s must have a single-field primary identifier", targetModelName)
		}

		targetPrimaryIdName := primaryID.Fields[0]
		colName := GetForeignKeyColumnName(relatedName, targetPrimaryIdName)
		columns = append(columns, colName)
		metas = append(metas, columnMeta{
			FieldName: relatedName,
			FieldType: yaml.ModelFieldTypeInteger,
			IsModelFK: true,
		})
	}

	return columns, metas, nil
}

func generateColumnValue(faker *gofakeit.Faker, meta columnMeta, rowIdx int, totalRows int, enumRowCounts map[string]int) any {
	if meta.IsOptional && faker.Number(0, 3) == 0 {
		return nil
	}

	if meta.IsAutoInc {
		return rowIdx + 1
	}

	if meta.IsEnumFK {
		enumCount, ok := enumRowCounts[meta.EnumName]
		if !ok || enumCount == 0 {
			enumCount = 3
		}
		return faker.Number(1, enumCount)
	}

	if meta.IsModelFK {
		return faker.Number(1, max(totalRows, 1))
	}

	return typemap.GenerateFieldValue(faker, meta.FieldName, meta.FieldType, meta.Hints)
}

func isRelationOne(relationType string) bool {
	return strings.HasSuffix(relationType, "one") ||
		relationType == "forone" ||
		relationType == "hasone" ||
		relationType == "foronepoly" ||
		relationType == "hasonepoly"
}

func GetModelDependencies(r *registry.Registry, model yaml.Model) []string {
	var deps []string
	for relatedName, relation := range model.Related {
		relationType := strings.ToLower(relation.Type)

		if strings.Contains(relationType, "poly") {
			continue
		}

		if strings.HasPrefix(relationType, "for") && isRelationOne(relationType) {
			targetModelName := relatedName
			if relation.Aliased != "" {
				targetModelName = relation.Aliased
			}
			if targetModelName != model.Name {
				deps = append(deps, targetModelName)
			}
		}
	}
	return deps
}

func SortModelsByDependency(r *registry.Registry) ([]yaml.Model, error) {
	allModels := r.GetAllModels()

	deps := make([]modelDependency, 0, len(allModels))
	for modelName, model := range allModels {
		deps = append(deps, modelDependency{
			Name:         modelName,
			Dependencies: GetModelDependencies(r, model),
		})
	}

	sorted, err := topologicalSortModels(deps)
	if err != nil {
		return nil, err
	}

	result := make([]yaml.Model, 0, len(sorted))
	for _, name := range sorted {
		result = append(result, allModels[name])
	}
	return result, nil
}

func topologicalSortModels(deps []modelDependency) ([]string, error) {
	graph := make(map[string][]string)
	allNodes := make(map[string]bool)

	for _, d := range deps {
		allNodes[d.Name] = true
		graph[d.Name] = d.Dependencies
		for _, dep := range d.Dependencies {
			allNodes[dep] = true
		}
	}

	inDegree := make(map[string]int)
	for node := range graph {
		if _, exists := inDegree[node]; !exists {
			inDegree[node] = 0
		}
		for _, dep := range graph[node] {
			if _, exists := graph[dep]; exists {
				inDegree[node]++
			}
		}
	}

	var queue []string
	for node := range graph {
		if inDegree[node] == 0 {
			queue = append(queue, node)
		}
	}

	sortStrings(queue)

	var result []string
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		var newQueue []string
		for node, nodeDeps := range graph {
			for _, dep := range nodeDeps {
				if dep == current {
					inDegree[node]--
					if inDegree[node] == 0 {
						newQueue = append(newQueue, node)
					}
				}
			}
		}
		sortStrings(newQueue)
		queue = append(queue, newQueue...)
	}

	if len(result) != len(graph) {
		return nil, fmt.Errorf("circular dependency detected in model definitions")
	}

	return result, nil
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}

func triggerModelFailure(hooks hook.CompileMorpheModel, config cfg.SeedConfig, model yaml.Model, err error) error {
	if hooks.OnCompileMorpheModelFailure != nil {
		return hooks.OnCompileMorpheModelFailure(config, model, err)
	}
	return err
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
