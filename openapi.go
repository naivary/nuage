//go:generate go tool go-enum --marshal --nocomments
package nuage

import (
	"fmt"
	"net/http"

	"github.com/google/jsonschema-go/jsonschema"
)

const _openAPIVersion = "3.1.0"

// ENUM(MIT, Apache-2.0)
type LicenseKeyword string

// ENUM(matrix, label, simple, form, spaceDelim, pipeDelim, deepObject, cookie)
type Style string

// ENUM(path, query, header, cookie)
type ParamIn string

type OpenAPI struct {
	Version           string               `json:"openapi"`
	Self              string               `json:"$self,omitempty"`
	Info              *Info                `json:"info"`
	JSONSchemaDialect string               `json:"jsonSchemaDialect,omitempty"`
	Servers           []*Server            `json:"servers,omitempty"`
	Paths             map[string]*PathItem `json:"paths"`
	Components        *Components          `json:"components,omitempty"`
}

func NewOpenAPI(info *Info) *OpenAPI {
	return &OpenAPI{
		Version: _openAPIVersion,
		Info:    info,
		Paths:   make(map[string]*PathItem),
	}
}

type PathItem struct {
	Summary     string     `json:"summary,omitempty"`
	Description string     `json:"description,omitempty"`
	Get         *Operation `json:"get,omitempty"`
	Put         *Operation `json:"put,omitempty"`
	Post        *Operation `json:"post,omitempty"`
	Delete      *Operation `json:"delete,omitempty"`
	Options     *Operation `json:"options,omitempty"`
	Head        *Operation `json:"head,omitempty"`
	Patch       *Operation `json:"patch,omitempty"`
	Trace       *Operation `json:"trace,omitempty"`
	Query       *Operation `json:"query,omitempty"`
}

func (p *PathItem) AddOperation(method string, op *Operation) error {
	switch method {
	case http.MethodGet:
		if p.Get != nil {
			return fmt.Errorf("add operation: operation for `%s` for method `%s` exists", op.OperationID, method)
		}
		p.Get = op
	case http.MethodPut:
		if p.Put != nil {
			return fmt.Errorf("add operation: operation for `%s` for method `%s` exists", op.OperationID, method)
		}
		p.Put = op
	case http.MethodPost:
		if p.Post != nil {
			return fmt.Errorf("add operation: operation for `%s` for method `%s` exists", op.OperationID, method)
		}
		p.Post = op
	case http.MethodDelete:
		if p.Delete != nil {
			return fmt.Errorf("add operation: operation for `%s` for method `%s` exists", op.OperationID, method)
		}
		p.Delete = op
	case http.MethodOptions:
		if p.Options != nil {
			return fmt.Errorf("add operation: operation for `%s` for method `%s` exists", op.OperationID, method)
		}
		p.Options = op
	case http.MethodHead:
		if p.Head != nil {
			return fmt.Errorf("add operation: operation for `%s` for method `%s` exists", op.OperationID, method)
		}
		p.Head = op
	case http.MethodPatch:
		if p.Patch != nil {
			return fmt.Errorf("add operation: operation for `%s` for method `%s` exists", op.OperationID, method)
		}
		p.Patch = op
	case http.MethodTrace:
		if p.Trace != nil {
			return fmt.Errorf("add operation: operation for `%s` for method `%s` exists", op.OperationID, method)
		}
		p.Trace = op
	}
	return nil
}

func (p *PathItem) OperationFor(method string) *Operation {
	switch method {
	case http.MethodGet:
		return p.Get
	case http.MethodPut:
		return p.Put
	case http.MethodPost:
		return p.Post
	case http.MethodDelete:
		return p.Delete
	case http.MethodOptions:
		return p.Options
	case http.MethodHead:
		return p.Head
	case http.MethodPatch:
		return p.Patch
	case http.MethodTrace:
		return p.Trace
}

type RequestBody struct {
	Description string                `json:"description,omitempty"`
	Required    bool                  `json:"required"`
	Content     map[string]*MediaType `json:"content,omitempty"`
}

type Response struct {
	Ref         string                `json:"$ref,omitempty"`
	Summary     string                `json:"summary,omitempty"`
	Description string                `json:"description,omitempty"`
	Headers     map[string]*Parameter `json:"headers,omitempty"`
	Content     map[string]*MediaType `json:"content,omitempty"`
}

type MediaType struct {
	Schema     *jsonschema.Schema `json:"schema,omitempty,omitzero"`
	ItemSchema *jsonschema.Schema `json:"itemSchema,omitempty"`
	Example    any                `json:"example,omitempty"`
}

type Components struct {
	Schemas         map[string]*jsonschema.Schema `json:"schemas,omitempty"`
	SecuritySchemes map[string]SecurityScheme     `json:"securitySchemes,omitempty"`
}

type Server struct {
	URL         string                     `json:"url"`
	Description string                     `json:"description,omitempty"`
	Name        string                     `json:"name,omitempty"`
	Variables   map[string]*ServerVariable `json:"variables,omitempty"`
}

type ServerVariable struct {
	Default     string   `json:"default"`
	Enum        []string `json:"enum,omitempty"`
	Description string   `json:"description,omitempty"`
}

type Info struct {
	Version        string   `json:"version"`
	Title          string   `json:"title"`
	Summary        string   `json:"summary,omitempty"`
	Description    string   `json:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
}

type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type License struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier,omitempty"`
	URL        string `json:"url,omitempty"`
}
