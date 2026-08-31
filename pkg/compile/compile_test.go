package compile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/go-util/assertfile"
	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-seed-data/internal/testutils"
	"github.com/kalo-build/plugin-morphe-seed-data/pkg/compile"
	"github.com/kalo-build/plugin-morphe-seed-data/pkg/compile/cfg"
)

type CompileTestSuite struct {
	assertfile.FileSuite

	TestDirPath            string
	TestGroundTruthDirPath string

	ModelsDirPath string
	EnumsDirPath  string
}

func TestCompileTestSuite(t *testing.T) {
	suite.Run(t, new(CompileTestSuite))
}

func (suite *CompileTestSuite) SetupTest() {
	suite.TestDirPath = testutils.GetTestDirPath()
	suite.TestGroundTruthDirPath = filepath.Join(suite.TestDirPath, "ground-truth", "compile-minimal")

	suite.ModelsDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "models")
	suite.EnumsDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "enums")
}

func (suite *CompileTestSuite) TearDownTest() {
	suite.TestDirPath = ""
}

func (suite *CompileTestSuite) TestMorpheToSeedSQL() {
	workingDirPath := filepath.Join(suite.TestDirPath, "working")
	suite.Nil(os.Mkdir(workingDirPath, 0755))
	defer os.RemoveAll(workingDirPath)

	config := compile.MorpheCompileConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryEnumsDirPath:      suite.EnumsDirPath,
			RegistryModelsDirPath:     suite.ModelsDirPath,
			RegistryStructuresDirPath: filepath.Join(suite.TestDirPath, "registry", "minimal", "structures"),
			RegistryEntitiesDirPath:   filepath.Join(suite.TestDirPath, "registry", "minimal", "entities"),
		},
		SeedConfig: cfg.SeedConfig{
			RowCount: 3,
			Seed:     42,
		},
		SeedWriter: &compile.SeedFileWriter{
			TargetDirPath: workingDirPath,
		},
	}

	compileErr := compile.MorpheToSeedSQL(config)
	suite.NoError(compileErr)

	// Enum seed data
	enumPath := filepath.Join(workingDirPath, "001_nationalities.sql")
	gtEnumPath := filepath.Join(suite.TestGroundTruthDirPath, "001_nationalities.sql")
	suite.FileExists(enumPath)
	suite.FileEquals(enumPath, gtEnumPath)

	// Company seed data (no FK deps, comes first among models)
	companyPath := filepath.Join(workingDirPath, "002_companies.sql")
	gtCompanyPath := filepath.Join(suite.TestGroundTruthDirPath, "002_companies.sql")
	suite.FileExists(companyPath)
	suite.FileEquals(companyPath, gtCompanyPath)

	// Person seed data (depends on Company via ForOne + Nationality enum)
	personPath := filepath.Join(workingDirPath, "003_people.sql")
	gtPersonPath := filepath.Join(suite.TestGroundTruthDirPath, "003_people.sql")
	suite.FileExists(personPath)
	suite.FileEquals(personPath, gtPersonPath)
}
