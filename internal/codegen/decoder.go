package codegen

import (
	"errors"
	"flag"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

const _requestModelNameSuffix = "Request"

type Parameter struct {
	Name           string
	GoType         string
	UnderlyingType string
}

type RequestModel struct {
	PkgName    string
	StructName string
	PathParams []Parameter
}

func GenDecoder(args []string) error {
	flagSet := flag.NewFlagSet("gen-decoder", flag.ExitOnError)
	if err := flagSet.Parse(args); err != nil {
		return err
	}

	cfg := &packages.Config{
		Mode: packages.LoadTypes | packages.LoadAllSyntax,
	}
	pkgs, err := packages.Load(cfg, flagSet.Args()...)
	if err != nil {
		return err
	}
	if exitCode := packages.PrintErrors(pkgs); exitCode > 0 {
		return errors.New("genDecoder: error while loading packages")
	}
	// Print the names of the source files
	// for each package listed on the command line.
	for _, pkg := range pkgs {
		err := findRequestStructs(pkg)
		if err != nil {
			return err
		}
	}
	return nil
}

func findRequestStructs(pkg *packages.Package) error {
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				if !strings.HasSuffix(typeSpec.Name.Name, _requestModelNameSuffix) {
					continue
				}
				typ := pkg.TypesInfo.TypeOf(typeSpec.Type)
				s, isStruct := typ.(*types.Struct)
				if !isStruct {
					continue
				}
				err := genDecoderOf(pkg, typeSpec, s)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func genDecoderOf(
	pkg *packages.Package,
	typeSpec *ast.TypeSpec,
	s *types.Struct,
) error {
	return nil
}
