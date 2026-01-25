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

		isSupported := isSupportedParamType(param.In, field.Type())
		if isSupported != nil {
			return nil, isSupported
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

// findParamType is trying to find the correct GoType and UnderlyingType
// while checking if the types are compatible with the infered parameter
// type. The following support matrix is validated:
// path parameter: int[8,16,32,64], uint[8,16,32,64], float[32,64], string, []T and map[string]string
// header parameter: int[8,16,32,64], uint[8,16,32,64], float[32,64], string, time.Time, time.Duration, []T
// cookie parameter: *http.Cookie
// query parameter: int[8,16,32,64], uint[8,16,32,64], float[32,64], string,
// map[string]string, time.Time, time.Duration, []T, struct{...}
func isSupportedParamType(in openapi.ParamIn, typ types.Type) error {
	ptr, isPtr := typ.(*types.Pointer)
	if isPtr {
		return isSupportedParamType(in, ptr.Elem())
	}
	alias, isAlias := typ.(*types.Alias)
	if isAlias {
		return isSupportedParamType(in, alias.Underlying())
	}
	switch in {
	case openapi.ParamInPath:
		switch t := typ.(type) {
		case *types.Named:
			return isSupportedParamType(in, t.Underlying())
		case *types.Struct:
			return errors.New("structs are not supported for path parameters")
		case *types.Basic:
			if t.Kind() == types.Bool {
				return errors.New("booleans are not supported for path parameters")
			}
		case *types.Map:
			key, isBasic := t.Key().(*types.Basic)
			if !isBasic {
				return errors.New("maps for path parameter can only have primitive types as key")
			}
			if key.Kind() != types.String {
				return errors.New("maps for path parameter can only have string as a key")
			}
			val, isBasic := t.Elem().(*types.Basic)
			if !isBasic {
				return errors.New("maps for path parameter can only have primitive types as values")
			}
			if val.Kind() != types.String {
				return errors.New("maps for path parameter can only have string as a value")
			}
		case *types.Slice:
			return isSupportedParamType(in, t.Elem())
		}
	case openapi.ParamInHeader:
		switch t := typ.(type) {
		case *types.Named:
			if t.String() != "time.Time" && t.String() != "time.Duration" {
				return isSupportedParamType(in, t.Underlying())
			}
		case *types.Struct:
			return errors.New("structs are not supported for header parameters")
		}
	case openapi.ParamInCookie:
		named, isNamed := typ.(*types.Named)
		if !isNamed {
			return errors.New("cookie parameter can only be stored in http.Cookie")
		}
		if named.String() != "http.Cookie" {
			return errors.New("cookie parameter can only be stored in http.Cookie")
		}
	}
	return nil
}
