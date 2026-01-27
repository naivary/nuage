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

const (
	KindMap    = "map"
	KindSlice  = "slice"
	KindField  = "field"
	KindPtr    = "ptr"
	KindStruct = "struct"
	KindAlias  = "alias"
	KindNamed  = "named"
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

	Type typeInfo
}

type typeInfo struct {
	Kind     string
	Name     string
	Children []*typeInfo
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
		info := ResolveType(typ)
		if info == nil {
			return nil, fmt.Errorf("type information can not be extracted: %s", field.Name())
		}
		r.Parameters = append(r.Parameters, &param)
	}
	return &r, nil
}

func ResolveType(typ types.Type) *typeInfo {
	switch t := typ.(type) {
	case *types.Pointer:
		return &typeInfo{
			Kind: KindPtr,
			Children: []*typeInfo{
				ResolveType(t.Elem()),
			},
		}
	case *types.Alias:
		return &typeInfo{
			Kind: KindAlias,
			Name: t.Obj().Name(),
			Children: []*typeInfo{
				ResolveType(t.Underlying()),
			},
		}
	case *types.Named:
		return &typeInfo{
			Kind: KindNamed,
			Name: t.Obj().Name(),
			Children: []*typeInfo{
				ResolveType(t.Underlying()),
			},
		}
	case *types.Struct:
		fields := make([]*typeInfo, 0, t.NumFields())
		for f := range t.Fields() {
			fields = append(fields, &typeInfo{
				Kind: KindField,
				Name: f.Name(),
				Children: []*typeInfo{
					ResolveType(f.Type()),
				},
			})
		}
		return &typeInfo{
			Kind:     KindStruct,
			Children: fields,
		}

	case *types.Map:
		return &typeInfo{
			Kind: KindMap,
			Children: []*typeInfo{
				{Kind: "key", Children: []*typeInfo{ResolveType(t.Key())}},
				{Kind: "value", Children: []*typeInfo{ResolveType(t.Elem())}},
			},
		}
	case *types.Slice:
		return &typeInfo{
			Kind: KindSlice,
			Children: []*typeInfo{
				ResolveType(t.Elem()),
			},
		}
	case *types.Basic:
		return &typeInfo{
			Kind: t.Name(), // "int", "string", etc.
		}
	default:
		return nil
	}
}
