package fileoperator

import (
	"fmt"
	"os"
	"strings"

	"github.com/arabian9ts/builder/pkg/builder"
)

func filterBuilderFile(info os.FileInfo) bool {
	if info.IsDir() {
		return false
	}

	if idx := strings.Index(info.Name(), "_builder"); 0 < idx {
		return false
	}

	if idx := strings.Index(info.Name(), "_accessor"); 0 < idx {
		return false
	}

	return true
}

func filterNonBuilderFile(info os.FileInfo) bool {
	if info.IsDir() {
		return false
	}

	if idx := strings.Index(info.Name(), "_builder"); 0 < idx {
		return true
	}

	if idx := strings.Index(info.Name(), "_accessor"); 0 < idx {
		return true
	}

	return false
}

func CleanBuilder(targetPkg string) error {
	pkg, err := builder.LoadPackage(targetPkg, filterNonBuilderFile)
	if err != nil {
		return err
	}

	files := pkg.ParsePkgFiles()
	for _, file := range files {
		pos := strings.LastIndex(file.FileName, ".")
		fileName := fmt.Sprintf("%s.go", file.FileName[:pos])
		err := os.Remove(fileName)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateBuilder(targetPkg string) error {
	pkg, err := builder.LoadPackage(targetPkg, filterBuilderFile)
	if err != nil {
		return err
	}

	files := pkg.ParsePkgFiles()
	for _, file := range files {
		code, hasBuilder := file.GenerateBuilder()
		if !hasBuilder {
			continue
		}

		pos := strings.LastIndex(file.FileName, ".")
		fileName := fmt.Sprintf("%s_builder.go", file.FileName[:pos])
		fp, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return err
		}
		defer fp.Close()

		fp.WriteString(code)
	}

	return nil
}

func CreateAccessor(targetPkg string) error {
	pkg, err := builder.LoadPackage(targetPkg, filterBuilderFile)
	if err != nil {
		return err
	}

	files := pkg.ParsePkgFiles()
	for _, file := range files {
		code, hasAccessor := file.GenerateAccessor()
		if !hasAccessor {
			continue
		}

		pos := strings.LastIndex(file.FileName, ".")
		fileName := fmt.Sprintf("%s_accessor.go", file.FileName[:pos])
		fp, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return err
		}
		defer fp.Close()

		fp.WriteString(code)
	}

	return nil
}
