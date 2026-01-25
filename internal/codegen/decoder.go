package codegen

import (
	"errors"
	"flag"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path"
	"reflect"
	"slices"
	"strings"
	"text/template"

	"github.com/naivary/nuage/openapi"
	"golang.org/x/tools/go/packages"
)

type requestModel struct {
	PkgName string
	Ident   string
	Params  []*parameter
}

type parameter struct {
	// Name of the parameter
	Name string
	// Location of the parameter
	In openapi.ParamIn

	// Custom or built-in go type
	GoType         string
	UnderlyingType string
	IsPointer      bool
}

type decoderData struct {
	// All required imports (types etc.)
	Imports      []string
	RequestModel *requestModel
}

func (d *decoderData) addImport(pkg string) {
	if slices.Contains(d.Imports, pkg) {
		return
	}
	d.Imports = append(d.Imports, pkg)
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
					data, err := genDecoder(pkg, typeSpec)
					if err != nil {
						return err
					}
					if data == nil {
						continue
					}

					// render the actual code
					tmpl, err := template.New("decoder.gotmpl").Funcs(template.FuncMap{
						"BitSize": BitSize,
					}).ParseGlob("templates/*.gotmpl")
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

func genDecoder(pkg *packages.Package, typeSpec *ast.TypeSpec) (*decoderData, error) {
	data := decoderData{
		Imports: make([]string, 0),
	}
	typ := pkg.TypesInfo.TypeOf(typeSpec.Type)
	s, isStructType := typ.(*types.Struct)
	if !isStructType {
		return nil, nil
	}
	ident := typeSpec.Name.Name
	if !strings.HasSuffix(ident, "Request") {
		return nil, nil
	}

	// found struct is a valid struct and will be considered for generation
	params := make([]*parameter, 0, s.NumFields())
	for i := range s.NumFields() {
		tag := reflect.StructTag(s.Tag(i))
		field := s.Field(i)
		param := parameter{
			In:   openapi.ParamLocation(tag),
			Name: field.Name(),
		}
		if param.In == "" {
			// field is not a parameter or is located at a invalid location
			continue
		}

		// TODO(naivary): This can be solved with recurision but rn
		// its better in the flow of the code like this
		fieldType := field.Type()
		ptr, isPtr := fieldType.(*types.Pointer)
		if isPtr {
			param.IsPointer = true
			fieldType = ptr.Elem()
		}
		// we need to identify if the data type used for the parameter is built-in or custom
		// TODO: check that the type used is compatiable with the param category. For example
		// path parameters canno tbe of type boolean in any way.
		switch t := fieldType.(type) {
		case *types.Named:
			param.GoType = path.Base(t.String())
			param.UnderlyingType = t.Underlying().String()
			data.addImport(importPath(t.String()))
		default:
			param.GoType = t.String()
			param.UnderlyingType = t.String()
		}
		params = append(params, &param)
	}
	data.RequestModel = &requestModel{
		PkgName: pkg.Name,
		Ident:   ident,
		Params:  params,
	}
	return &data, nil
}

func importPath(symbol string) string {
	if i := strings.LastIndex(symbol, "."); i != -1 {
		return symbol[:i]
	}
	return symbol
}

func isValidParamType(param openapi.ParamIn, underlyingType types.Type) error {
    // TODO: Should structs be allowed for deep objects?
	switch underlyingType.(type) {
	case *types.Struct, *types.Chan, *types.Signature:
		return errors.New("not seriliasable")
	}
	// TODO: disallow general types like func, chan, interface, (struct)
	switch param {
	case openapi.ParamInPath:
		basic, isBasic := underlyingType.(*types.Basic)
		if !isBasic {
			return nil
		}
		if basic.Kind() == types.Bool {
			// TODO: detailed error message for the user
			return errors.New("path parameter cannot be of type boolean")
		}
	}
	return nil
}
