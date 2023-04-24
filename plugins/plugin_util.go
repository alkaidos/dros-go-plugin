package plugins

import (
	"errors"
	"github.com/isyscore/isc-gobase/file"
	"os"
	"path"
	"strings"
)

func GetFileInfo(fileName string) ([]byte, error) {
	dir, _ := os.Getwd()
	pkg := strings.Replace(dir, "\\", "/", -1)
	fileName = path.Join(pkg, "", fileName)
	if !file.FileExists(fileName) {
		return nil, errors.New("找不到对应的swagger文件，pathName：" + fileName)
	}
	content, err := os.ReadFile(fileName)
	return content, err
}
