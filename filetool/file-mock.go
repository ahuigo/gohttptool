package filetool

/*
import (
	"github.com/spf13/afero"
)

func CreateMockFile(content []byte) error {
	var fs = afero.NewMemMapFs()
	file, err := fs.Create("test.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	return nil
}
*/
