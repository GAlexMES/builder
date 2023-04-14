package builder

import (
	"go/ast"
	"go/token"
	"go/types"

	. "github.com/dave/jennifer/jen"
)

type PkgFile struct {
	astFile  *ast.File
	fset     *token.FileSet
	gendecls []*ast.GenDecl
	FileName string
	PkgName  string
	pkgScope *types.Scope
}

func (file PkgFile) GenerateBuilder() (string, bool) {
	f := NewFile(file.PkgName)

	structs := file.parsePkgStructs()
	if len(structs) == 0 {
		return "", false
	}

	for _, st := range structs {
		st.DefineBuilderStruct(f)
		st.DefineBuilderInitializer(f)
		st.DefineBuilderConstructors(f)
		st.DefineBuildFunc(f)
	}

	return f.GoString(), true
}

func (file PkgFile) GenerateAccessor() (string, bool) {
	f := NewFile(file.PkgName)

	structs := file.parsePkgStructs()
	if len(structs) == 0 {
		return "", false
	}

	for _, st := range structs {
		println(st.name)
		st.DefineAccessors(f)
	}

	return f.GoString(), true
}

func (file PkgFile) parsePkgStructs() (pkgStructs []PkgStruct) {
	for _, decl := range file.gendecls {
		for _, spec := range decl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			st := file.pkgScope.Lookup(typeSpec.Name.Name)
			structMeta, ok := st.Type().Underlying().(*types.Struct)
			if !ok {
				continue
			}

			pkgStruct := PkgStruct{
				fset: file.fset,
				name: typeSpec.Name.Name,
				meta: structMeta,
			}
			pkgStructs = append(pkgStructs, pkgStruct)
		}
	}

	return
}
