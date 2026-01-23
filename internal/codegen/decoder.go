package codegen

import (
	"errors"
	"flag"
	"go/ast"
	"go/types"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/naivary/nuage/openapi"
	"golang.org/x/tools/go/packages"
)

type paramInfo struct {
	GoType         string
	UnderlyingType string
	Name           string
	In             openapi.ParamIn
}

type requestModelInfo struct {
	PkgName     string
	StructIdent string
	Params      []*paramInfo
}

func GenDecoder(args []string) error {
	fs := flag.NewFlagSet("decoders", flag.ExitOnError)
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	cfg := &packages.Config{
		Mode: packages.LoadTypes | packages.LoadAllSyntax,
	}
	pkgs, err := packages.Load(cfg, fs.Args()...)
	if err != nil {
		return err
	}
	if exitCode := packages.PrintErrors(pkgs); exitCode > 0 {
		return errors.New("GenDecoder: error while loading packages")
	}
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				if err := genDecoder(pkg, file, decl); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func genDecoder(pkg *packages.Package, file *ast.File, decl ast.Decl) error {
	genDecl, isGenDecl := decl.(*ast.GenDecl)
	if !isGenDecl {
		return nil
	}
	requestModelInfos := make([]*requestModelInfo, 0, len(genDecl.Specs))
	for _, spec := range genDecl.Specs {
		typeSpec, isTypeSpec := spec.(*ast.TypeSpec)
		if !isTypeSpec {
			return nil
		}
		typ := pkg.TypesInfo.TypeOf(typeSpec.Type)
		s, isStructType := typ.(*types.Struct)
		if !isStructType {
			return nil
		}
		structIdent := typeSpec.Name.Name
		if !strings.HasSuffix(structIdent, "Request") {
			return nil
		}

		paramInfos := make([]*paramInfo, 0, s.NumFields())
		for i := range s.NumFields() {
			tag := reflect.StructTag(s.Tag(i))
			field := s.Field(i)
			paramIn := openapi.ParamLocation(tag)
			if paramIn == "" {
				// field is not a parameter or located at a invalid location
				continue
			}
			fieldInfo := &paramInfo{
				In:   paramIn,
				Name: field.Name(),
			}
			switch t := field.Type().(type) {
			case *types.Named:
				fieldInfo.GoType = t.String()
				fieldInfo.UnderlyingType = t.Underlying().String()
			default:
				fieldInfo.GoType = t.String()
				fieldInfo.UnderlyingType = t.String()
			}
			paramInfos = append(paramInfos, fieldInfo)
		}

		reqModelInfo := &requestModelInfo{
			PkgName:     pkg.Name,
			StructIdent: structIdent,
			Params:      paramInfos,
		}
		requestModelInfos = append(requestModelInfos, reqModelInfo)
	}

	tmpl, err := template.ParseGlob("templates/*.gotmpl")
	if err != nil {
		return err
	}
	for _, reqModelInfo := range requestModelInfos {
		err = tmpl.Execute(os.Stdout, &reqModelInfo)
		if err != nil {
			return err
		}
	}
	return nil
}
