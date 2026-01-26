package codegen

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/naivary/nuage/internal/openapiutil"
	"github.com/naivary/nuage/openapi"
	"golang.org/x/tools/go/packages"
)

type requestModel struct {
	// Import statments defined by the request model
	Imports []string

	// Identifier
	Ident string

	// Package in which the request model was found
	PkgName string

	// Parameters infered from the fields of the request model
	Parameters []*parameter
}

type parameter struct {
	// Identifier defined by the user in the struct tag
	Ident string

	// Identifier of the field in the struct where the
	// parameter is found.
	FieldIdent string

	// Location of the parameter
	In openapi.ParamIn

	// GoType is always a built-in go type and always contains
	// a non-empty value.
	GoType typeInfo

	// UnderlyingType is always empty for non named types (e.g. type Named T).
	// For named types it contains informatino about the underlying type of the named
	// type.
	UnderlyingType typeInfo
}

type typeInfo struct {
	// Whether the type was a pointer
	IsPointer bool

	// Name of the type found. For a named type it is the
	// correct package type name (e.g. time.Time). For built-in types it is
	// the name of the Go type itself (e.g. int, string etc.)
	Type string

	// Fields contains the type information of all fields if the
	// underlying type is a struct.
	Fields []typeInfo
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
				genDecl, isGenDecl := decl.(*ast.GenDecl)
				if !isGenDecl {
					continue
				}
				if genDecl.Tok != token.TYPE {
					continue
				}
				for _, spec := range genDecl.Specs {
					typeSpec, isTypeSpec := spec.(*ast.TypeSpec)
					if !isTypeSpec {
						continue
					}
					ident := typeSpec.Name.Name
					if !strings.HasSuffix(ident, "Request") {
						continue
					}
					typ := pkg.TypesInfo.TypeOf(typeSpec.Type)
					s, isStruct := typ.(*types.Struct)
					if !isStruct {
						continue
					}

					// From now the decl is considered a valid request model
					// and will be analysed for code generation.
					data, err := genDecoder(pkg, ident, s)
					if err != nil {
						return err
					}
					if data == nil {
						continue
					}

					// render the actual code
					tmpl, err := template.New("decoder.gotmpl").Funcs(FuncsMap).ParseGlob("templates/*.gotmpl")
					if err != nil {
						return err
					}
					if err := tmpl.Execute(os.Stdout, data); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func genDecoder(pkg *packages.Package, ident string, s *types.Struct) (*requestModel, error) {
	r := requestModel{
		PkgName:    pkg.Name,
		Ident:      ident,
		Parameters: make([]*parameter, 0, s.NumFields()),
		Imports:    make([]string, 0, 1),
	}
	for i := range s.NumFields() {
		tag := reflect.StructTag(s.Tag(i))
		field := s.Field(i)
		param := parameter{
			FieldIdent: field.Name(),
			In:         openapiutil.ParamLocation(tag),
		}
		if param.In == "" {
			// field is not a parameter or at an invalid location
			continue
		}
		val := tag.Get(string(param.In))
		if len(val) == 0 {
			return nil, fmt.Errorf("struct field tag is missing the name of the parameter: %s", field.Name())
		}
		param.Ident = strings.Split(val, ",")[0]

		typ := field.Type()
		err := isSupportedParamType(param.In, typ)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", err, field.Name())
		}
		r.Parameters = append(r.Parameters, &param)
	}
	return &r, nil
}

func resolveStructParamType(s *types.Struct) *parameter {
	param := &parameter{}
	for field := range s.Fields() {
	}
	return nil
}
