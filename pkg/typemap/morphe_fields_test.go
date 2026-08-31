package typemap_test

import (
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/plugin-morphe-seed-data/pkg/typemap"
)

type MorpheFieldsTestSuite struct {
	suite.Suite
	faker *gofakeit.Faker
}

func TestMorpheFieldsTestSuite(t *testing.T) {
	suite.Run(t, new(MorpheFieldsTestSuite))
}

func (suite *MorpheFieldsTestSuite) SetupTest() {
	suite.faker = gofakeit.New(42)
}

func (suite *MorpheFieldsTestSuite) TestGenerateFieldValue_UUID() {
	val := typemap.GenerateFieldValue(suite.faker, "ID", yaml.ModelFieldTypeUUID, typemap.FormatHints{})
	str, ok := val.(string)
	suite.True(ok)
	suite.Len(str, 36)
	suite.Contains(str, "-")
}

func (suite *MorpheFieldsTestSuite) TestGenerateFieldValue_AutoIncrement() {
	val := typemap.GenerateFieldValue(suite.faker, "ID", yaml.ModelFieldTypeAutoIncrement, typemap.FormatHints{})
	suite.IsType("", val)
}

func (suite *MorpheFieldsTestSuite) TestGenerateFieldValue_String_Default() {
	val := typemap.GenerateFieldValue(suite.faker, "SomeField", yaml.ModelFieldTypeString, typemap.FormatHints{})
	str, ok := val.(string)
	suite.True(ok)
	suite.NotEmpty(str)
}

func (suite *MorpheFieldsTestSuite) TestGenerateFieldValue_String_FormatEmail() {
	val := typemap.GenerateFieldValue(suite.faker, "Email", yaml.ModelFieldTypeString, typemap.FormatHints{Format: "email"})
	str, ok := val.(string)
	suite.True(ok)
	suite.Contains(str, "@")
}

func (suite *MorpheFieldsTestSuite) TestGenerateFieldValue_String_FormatFirstName() {
	val := typemap.GenerateFieldValue(suite.faker, "Name", yaml.ModelFieldTypeString, typemap.FormatHints{Format: "firstName"})
	str, ok := val.(string)
	suite.True(ok)
	suite.NotEmpty(str)
	suite.False(strings.Contains(str, " "))
}

func (suite *MorpheFieldsTestSuite) TestGenerateFieldValue_String_MaxLength() {
	val := typemap.GenerateFieldValue(suite.faker, "Code", yaml.ModelFieldTypeString, typemap.FormatHints{MaxLength: 5})
	str, ok := val.(string)
	suite.True(ok)
	suite.LessOrEqual(len(str), 5)
}

func (suite *MorpheFieldsTestSuite) TestGenerateFieldValue_String_MinLength() {
	val := typemap.GenerateFieldValue(suite.faker, "Code", yaml.ModelFieldTypeString, typemap.FormatHints{MinLength: 50})
	str, ok := val.(string)
	suite.True(ok)
	suite.GreaterOrEqual(len(str), 50)
}

func (suite *MorpheFieldsTestSuite) TestGenerateFieldValue_String_Regex() {
	val := typemap.GenerateFieldValue(suite.faker, "TaxID", yaml.ModelFieldTypeString, typemap.FormatHints{Regex: `^\d{2}-\d{7}$`})
	str, ok := val.(string)
	suite.True(ok)
	suite.Regexp(`^\d{2}-\d{7}$`, str)
}

func (suite *MorpheFieldsTestSuite) TestGenerateFieldValue_Integer() {
	val := typemap.GenerateFieldValue(suite.faker, "Count", yaml.ModelFieldTypeInteger, typemap.FormatHints{})
	_, ok := val.(int)
	suite.True(ok)
}

func (suite *MorpheFieldsTestSuite) TestGenerateFieldValue_Integer_MinMax() {
	for i := 0; i < 20; i++ {
		val := typemap.GenerateFieldValue(suite.faker, "Age", yaml.ModelFieldTypeInteger, typemap.FormatHints{HasMin: true, Min: 18, HasMax: true, Max: 65})
		v, ok := val.(int)
		suite.True(ok)
		suite.GreaterOrEqual(v, 18)
		suite.LessOrEqual(v, 65)
	}
}

func (suite *MorpheFieldsTestSuite) TestGenerateFieldValue_Float() {
	val := typemap.GenerateFieldValue(suite.faker, "Amount", yaml.ModelFieldTypeFloat, typemap.FormatHints{})
	_, ok := val.(float64)
	suite.True(ok)
}

func (suite *MorpheFieldsTestSuite) TestGenerateFieldValue_Boolean() {
	val := typemap.GenerateFieldValue(suite.faker, "Active", yaml.ModelFieldTypeBoolean, typemap.FormatHints{})
	_, ok := val.(bool)
	suite.True(ok)
}

func (suite *MorpheFieldsTestSuite) TestGenerateFieldValue_Time() {
	val := typemap.GenerateFieldValue(suite.faker, "CreatedAt", yaml.ModelFieldTypeTime, typemap.FormatHints{})
	str, ok := val.(string)
	suite.True(ok)
	suite.Contains(str, "T")
	suite.Contains(str, "Z")
}

func (suite *MorpheFieldsTestSuite) TestGenerateFieldValue_Date() {
	val := typemap.GenerateFieldValue(suite.faker, "BirthDate", yaml.ModelFieldTypeDate, typemap.FormatHints{})
	str, ok := val.(string)
	suite.True(ok)
	suite.Regexp(`^\d{4}-\d{2}-\d{2}$`, str)
}

func (suite *MorpheFieldsTestSuite) TestGenerateFieldValue_Protected() {
	val := typemap.GenerateFieldValue(suite.faker, "SSN", yaml.ModelFieldTypeProtected, typemap.FormatHints{})
	str, ok := val.(string)
	suite.True(ok)
	suite.NotEmpty(str)
}

func (suite *MorpheFieldsTestSuite) TestGenerateFieldValue_Sealed() {
	val := typemap.GenerateFieldValue(suite.faker, "Password", yaml.ModelFieldTypeSealed, typemap.FormatHints{})
	str, ok := val.(string)
	suite.True(ok)
	suite.True(strings.HasPrefix(str, "$2a$10$"))
}

func (suite *MorpheFieldsTestSuite) TestGenerateFieldValue_String_InferredEmail() {
	f := gofakeit.New(100)
	val := typemap.GenerateFieldValue(f, "Email", yaml.ModelFieldTypeString, typemap.FormatHints{})
	str, ok := val.(string)
	suite.True(ok)
	suite.Contains(str, "@")
}

func (suite *MorpheFieldsTestSuite) TestGenerateFieldValue_String_InferredName() {
	f := gofakeit.New(100)
	val := typemap.GenerateFieldValue(f, "FirstName", yaml.ModelFieldTypeString, typemap.FormatHints{})
	str, ok := val.(string)
	suite.True(ok)
	suite.NotEmpty(str)
}

func (suite *MorpheFieldsTestSuite) TestGenerateFieldValue_Deterministic() {
	f1 := gofakeit.New(42)
	f2 := gofakeit.New(42)
	val1 := typemap.GenerateFieldValue(f1, "Name", yaml.ModelFieldTypeString, typemap.FormatHints{Format: "name"})
	val2 := typemap.GenerateFieldValue(f2, "Name", yaml.ModelFieldTypeString, typemap.FormatHints{Format: "name"})
	suite.Equal(val1, val2)
}
