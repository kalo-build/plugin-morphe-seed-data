package compile

import (
	"github.com/kalo-build/plugin-morphe-seed-data/pkg/sqlfile"
)

type SeedFileWriter struct {
	TargetDirPath string
}

func (w *SeedFileWriter) WriteSeedFile(fileName string, contents string, order int) ([]byte, error) {
	return sqlfile.WriteSQLDefinitionFileWithOrder(w.TargetDirPath, fileName, contents, order)
}
