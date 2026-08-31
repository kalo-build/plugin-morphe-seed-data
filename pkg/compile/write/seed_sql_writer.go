package write

type SeedSQLWriter interface {
	WriteSeedFile(fileName string, contents string, order int) ([]byte, error)
}
