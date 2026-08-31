package sqlfile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kalo-build/go-util/strcase"
)

func WriteSQLDefinitionFileWithOrder(dirPath string, definitionName string, contents string, order int) ([]byte, error) {
	definitionFileName := strcase.ToSnakeCaseLower(definitionName)

	if order > 0 {
		definitionFileName = fmt.Sprintf("%03d_%s", order, definitionFileName)
	}

	definitionFilePath := filepath.Join(dirPath, definitionFileName+".sql")
	if _, readErr := os.ReadDir(dirPath); readErr != nil && os.IsNotExist(readErr) {
		mkDirErr := os.MkdirAll(dirPath, 0755)
		if mkDirErr != nil {
			return nil, mkDirErr
		}
	}
	return []byte(contents), os.WriteFile(definitionFilePath, []byte(contents), 0644)
}
