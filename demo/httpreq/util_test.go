package httpreq

import (
	"os"
	"path/filepath"
)

func getTestDataPath(filename string) string {
	pwd, _ := os.Getwd()
	return filepath.Join(pwd, "./testdata", filename)
}
