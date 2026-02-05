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
	kindPtr      = "ptr"
	kindMap      = "map"
	kindMapKey   = "map.key"
	kindMapValue = "map.value"
	kindSlice    = "slice"
	kindStruct   = "struct"
	kindAlias    = "alias"
	kindNamed    = "named"
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

	TypeInfo *typeInfo

	Opts *openapiutil.ParamOpts
}

type typeInfo struct {
	Kind string

	// Ident is the identifier of the named type or if its of kind `field`
	// the identifier of the field in the struct.
	Ident string

	// Package in which the type is defined. If the type is defined in the
	// same package as the request model it will be empty.
	Pkg      string
	Children []*typeInfo
}

func GenDecoder(args []string) error {
	fs := flag.NewFlagSet("decoder", flag.ExitOnError)
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
					ident, s := isRequestModel(pkg, spec)
					if s == nil {
						continue
					}
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
					if err := tmpl.ExecuteTemplate(os.Stdout, "decoder", data); err != nil {
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
		opts, err := openapiutil.ParseParamOpts(tag)
		if err != nil {
			return nil, err
		}
		param.Ident = opts.Name
		param.Opts = opts

		typ := field.Type()
		err = isSupportedParamType(opts, typ)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", err, field.Name())
		}
		info := resolveType(typ)
		if info == nil {
			return nil, fmt.Errorf("type information can not be extracted: %s", field.Name())
		}
		param.TypeInfo = info
		r.Imports = append(r.Imports, resolveImports(pkg, info)...)
		r.Parameters = append(r.Parameters, &param)
	}
	return &r, nil
}

// isRequestModel reports whether `spec` is a request model in the context
// of nuage and should be considered for generation of code.
func isRequestModel(pkg *packages.Package, spec ast.Spec) (string, *types.Struct) {
	typeSpec, isTypeSpec := spec.(*ast.TypeSpec)
	if !isTypeSpec {
		return "", nil
	}
	ident := typeSpec.Name.Name
	if !strings.HasSuffix(ident, "Request") {
		return "", nil
	}
	typ := pkg.TypesInfo.TypeOf(typeSpec.Type)
	s, isStruct := typ.(*types.Struct)
	if !isStruct {
		return "", nil
	}
	return ident, s
}

func resolveType(typ types.Type) *typeInfo {
	switch t := typ.(type) {
	case *types.Pointer:
		return &typeInfo{
			Kind: kindPtr,
			Children: []*typeInfo{
				resolveType(t.Elem()),
			},
		}
	case *types.Alias:
		return &typeInfo{
			Kind:  kindAlias,
			Ident: t.Obj().Name(),
			Children: []*typeInfo{
				resolveType(t.Underlying()),
			},
		}
	case *types.Named:
		return &typeInfo{
			Kind:  kindNamed,
			Ident: t.Obj().Name(),
			Pkg:   t.Obj().Pkg().Name(),
			Children: []*typeInfo{
				resolveType(t.Underlying()),
			},
		}
	case *types.Struct:
		switch t.String() {
		case _timeTypeName, _cookieTypeName:
			return nil
		}
		fields := make([]*typeInfo, 0, t.NumFields())
		for f := range t.Fields() {
			info := resolveType(f.Type())
			info.Ident = f.Name()
			fields = append(fields, info)
		}
		return &typeInfo{
			Kind:     kindStruct,
			Children: fields,
		}
	case *types.Map:
		return &typeInfo{
			Kind: kindMap,
			Children: []*typeInfo{
				{Kind: kindMapKey, Children: []*typeInfo{resolveType(t.Key())}},
				{Kind: kindMapValue, Children: []*typeInfo{resolveType(t.Elem())}},
			},
		}
	case *types.Slice:
		return &typeInfo{
			Kind: kindSlice,
			Children: []*typeInfo{
				resolveType(t.Elem()),
			},
		}
	case *types.Basic:
		return &typeInfo{
			Kind: t.Name(),
		}
	default:
		return nil
	}
}

func resolveImports(pkg *packages.Package, info *typeInfo) []string {
	if info.Kind != kindNamed {
		return nil
	}
	imports := make([]string, 0, 1)
	if info.Pkg != pkg.Name {
		imports = append(imports, info.Pkg)
	}
	for _, i := range info.Children {
		imports = append(imports, resolveImports(pkg, i)...)
	}
	return imports
}
