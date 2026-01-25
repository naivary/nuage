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
	"slices"
	"strings"
	"text/template"

	"github.com/naivary/nuage/internal/typesutil"
	"github.com/naivary/nuage/openapi"
	"golang.org/x/tools/go/packages"
)

const (
	_timeTypeName   = "time.Time"
	_cookieTypeName = "http.Cookie"
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
					// and will be analysed for code generation
					data, err := genDecoder(pkg, ident, s)
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

func genDecoder(pkg *packages.Package, ident string, s *types.Struct) (*decoderData, error) {
	data := decoderData{
		Imports: make([]string, 0),
	}
	params := make([]*parameter, 0, s.NumFields())
	for i := range s.NumFields() {
		tag := reflect.StructTag(s.Tag(i))
		field := s.Field(i)
		param := parameter{
			Name: field.Name(),
			In:   openapi.ParamLocation(tag),
		}
		if param.In == "" {
			// field is not a parameter or is located at a invalid location
			continue
		}
		if err := isSupportedParamType(param.In, field); err != nil {
			return nil, err
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

func isSupportedParamType(in openapi.ParamIn, field *types.Var) error {
	typ := field.Type()
	ptr, isPtr := typ.(*types.Pointer)
	if isPtr {
		typ = ptr.Elem()
	}

	// the following types are not serializable and thus not supported
	// by any of the parameter types.
	switch t := typ.(type) {
	case *types.Signature, *types.Chan:
		return fmt.Errorf("functions or channels are not supported as parameter types because they are not serializable: %s", field.Name())
	case *types.Basic:
		kind := t.Kind()
		if kind == types.Uintptr || kind == types.UnsafePointer {
			return errors.New("uintptr and unsafe pointer are not supported as parameters of any kind")
		}
	case *types.Array:
		return errors.New("arrays are a valid type parameter in general but are hard to support for nuage. For this version it is not supported. Use slices instead")
	}

	switch in {
	case openapi.ParamInPath:
		return isSupportedPathParamType(field, field.Type(), false)
	case openapi.ParamInHeader:
		return isSupportedHeaderParamType(field, field.Type(), false)
	case openapi.ParamInQuery:
		return isSupportedQueryType(field, field.Type(), false)
	case openapi.ParamInCookie:
		return isSupportedCookieType(field, field.Type())
	}
	return nil
}

func isSupportedPathParamType(field *types.Var, typ types.Type, isSlice bool) error {
	switch t := typ.(type) {
	case *types.Alias:
		return isSupportedPathParamType(field, t.Underlying(), false)
	case *types.Pointer:
		return isSupportedPathParamType(field, t.Elem(), false)
	case *types.Basic:
		kind := t.Kind()
		if kind == types.Bool {
			return fmt.Errorf("path parameters cannot be of type boolean: %s", field.Name())
		}
		if typesutil.IsComplex(kind) {
			return fmt.Errorf("path parameters cannot be of type complex: %s", field.Name())
		}
		if typesutil.IsFloat(kind) {
			return fmt.Errorf("path parameters cannot be of type float: %s", field.Name())
		}
	case *types.Named:
		return isSupportedPathParamType(field, t.Underlying(), false)
	case *types.Slice:
		if isSlice {
			return fmt.Errorf("path parameters cannot be nested slices: %s", field.Name())
		}
		return isSupportedPathParamType(field, t.Elem(), true)
	default:
		return fmt.Errorf("type `%s` is not supported for path parameters", typ.String())
	}
	return nil
}

func isSupportedHeaderParamType(field *types.Var, typ types.Type, isSlice bool) error {
	switch t := typ.(type) {
	case *types.Alias:
		return isSupportedHeaderParamType(field, t.Underlying(), false)
	case *types.Pointer:
		return isSupportedHeaderParamType(field, t.Elem(), false)
	case *types.Basic:
		kind := t.Kind()
		if kind == types.Bool {
			return fmt.Errorf("header parameters cannot be of type boolean: %s", field.Name())
		}
		if typesutil.IsComplex(kind) {
			return fmt.Errorf("header parameters cannot be of type complex: %s", field.Name())
		}
		if typesutil.IsFloat(kind) {
			return fmt.Errorf("header parameters cannot be of type float: %s", field.Name())
		}
	case *types.Named:
		name := t.String()
		if name == _timeTypeName {
			return nil
		}
		return isSupportedHeaderParamType(field, t.Underlying(), false)
	case *types.Slice:
		if isSlice {
			return fmt.Errorf("header parameters cannot be nested slices: %s", field.Name())
		}
		return isSupportedHeaderParamType(field, t.Elem(), true)
	default:
		return fmt.Errorf("type `%s` is not supported for header parameters", typ.String())
	}
	return nil
}

func isSupportedQueryType(field *types.Var, typ types.Type, isSlice bool) error {
	switch t := typ.(type) {
	case *types.Alias:
		return isSupportedQueryType(field, t.Underlying(), false)
	case *types.Pointer:
		return isSupportedQueryType(field, t.Elem(), false)
	case *types.Basic:
		kind := t.Kind()
		if typesutil.IsFloat(kind) {
			return fmt.Errorf("query parameters cannot be of type float: %s", field.Name())
		}
		if typesutil.IsComplex(kind) {
			return fmt.Errorf("query parameters cannot be of type complex: %s", field.Name())
		}
	case *types.Named:
		name := t.String()
		if name == _timeTypeName {
			return nil
		}
		return isSupportedQueryType(field, t.Underlying(), false)
	case *types.Slice:
		if isSlice {
			return fmt.Errorf("query parameters cannot be nested slices: %s", field.Name())
		}
		return isSupportedQueryType(field, t.Elem(), true)
	case *types.Map:
		isKeyTypeBasic := typesutil.IsBasic(t.Key(), true)
		isValTypeBasic := typesutil.IsBasic(t.Elem(), true)
		if !isKeyTypeBasic || !isValTypeBasic {
			return fmt.Errorf("map types for query parameters can only be of form map[string]string")
		}
		key := t.Key().(*types.Basic)
		if key.Kind() != types.String {
			return fmt.Errorf("the map key type of a query parameters has to be string")
		}
		val := t.Elem().(*types.Basic)
		if val.Kind() != types.String {
			return fmt.Errorf("the map value type of a query parameters has to be string")
		}
	case *types.Struct:
		for field := range t.Fields() {
			err := isSupportedQueryType(field, field.Type(), false)
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("type `%s` is not supported for query parameters", typ.String())
	}
	return nil
}

func isSupportedCookieType(field *types.Var, typ types.Type) error {
	switch t := typ.(type) {
	case *types.Alias:
		return isSupportedCookieType(field, t.Underlying())
	case *types.Pointer:
		return isSupportedCookieType(field, t.Elem())
	case *types.Named:
		name := t.String()
		if name == _cookieTypeName {
			return nil
		}
		return isSupportedCookieType(field, t.Underlying())
	default:
		return fmt.Errorf("type `%s` is not supported for cookie parameters", typ.String())
	}
}
