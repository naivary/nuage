package openapi

import (
	"net/url"

	"github.com/google/jsonschema-go/jsonschema"
)

const (
	// OpenAPIVersion is the OpenAPI Specification version implemented
	// and supported by this framework.
	//
	// This value is used when generating OpenAPI documents to indicate
	// the exact specification version the output conforms to.
	OpenAPIVersion = "3.2.0"

	// JSONSchemaDialect is the default JSON Schema dialect URI used by
	// the framework when producing OpenAPI schemas.
	//
	// It identifies the JSON Schema specification that schema definitions
	// are expected to follow, and is typically referenced from the
	// OpenAPI document's `jsonSchemaDialect` field.
	JSONSchemaDialect = "https://json-schema.org/draft/2020-12/json-schema-core.html"
)

type ParamIn string

const (
	ParamInPath   ParamIn = "path"
	ParamInQuery  ParamIn = "query"
	ParamInHeader ParamIn = "header"
	ParamInCookie ParamIn = "cookie"
)

type ParamStyle string

const (
	ParamStyleMatrix     ParamStyle = "matrix"
	ParamStyleLabel      ParamStyle = "label"
	ParamStyleSimple     ParamStyle = "simple"
	ParamStyleForm       ParamStyle = "form"
	ParamStyleSpaceDelim ParamStyle = "spaceDelim"
	ParamStylePipeDelim  ParamStyle = "pipeDelim"
	ParamStyleDeepObject ParamStyle = "deepObject"
	ParamStyleCookie     ParamStyle = "deepObject"
)

type SecurityRequirement map[string][]string

type SecurityType string

const (
	SecurityTypeAPIKey        SecurityType = "apiKey"
	SecurityTypeHTTP          SecurityType = "http"
	SecurityTypeMutualTLS     SecurityType = "mutualTLS"
	SecurityTypeOAuth2        SecurityType = "oauth2"
	SecurityTypeOpenIDConnect SecurityType = "openIdConnect"
)

type OpenAPI struct {
	Version           string               `json:"openapi"`
	Self              string               `json:"$self"`
	Info              *Info                `json:"info"`
	JSONSchemaDialect string               `json:"jsonSchemaDialect,omitempty"`
	Servers           []*Server            `json:"servers,omitempty"`
	Paths             map[string]*PathItem `json:"paths"`
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
type Server struct {
	URL         string                     `json:"url"`
	Description string                     `json:"description,omitempty"`
	Name        string                     `json:"name,omitempty"`
	Variables   map[string]*ServerVariable `json:"variables,omitempty"`
}

type ServerVariable struct {
	Enum        []string `json:"enum,omitempty"`
	Default     string   `json:"default"`
	Description string   `json:"description,omitempty"`
}
type Contact struct {
	Name  string  `json:"name,omitempty"`
	URL   url.URL `json:"url"`
	Email string  `json:"email,omitempty"`
}

type License struct {
	Name       string  `json:"name"`
	Identifier string  `json:"identifier,omitempty"`
	URL        url.URL `json:"url"`
}

type Components struct {
	Schemas         map[string]*jsonschema.Schema
	Responses       map[string]*Response
	SecuritySchemes map[string]*SecurityScheme
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

type Operation struct {
	Tags        []string              `json:"tags,omitempty"`
	Summary     string                `json:"summary,omitempty"`
	Description string                `json:"description,omitempty"`
	OperationID string                `json:"operationId,omitempty"`
	Deprecated  bool                  `json:"deprecated,omitempty"`
	Security    []SecurityRequirement `json:"security,omitempty"`

	Parameters  []*Parameter         `json:"parameters,omitempty"`
	RequestBody *RequestBody         `json:"requestBody,omitempty"`
	Responses   map[string]*Response `json:"responses,omitempty"`

	// nuage specific
	Pattern string `json:"-"`

	ResponseContentType string `json:"-"`
	ResponseDesc        string `json:"-"`
	ResponseStatusCode  int    `json:"-"`

	RequestContentType    string `json:"-"`
	RequestDesc           string `json:"-"`
	IsRequestBodyRequired *bool  `json:"-"`
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

type Parameter struct {
	Name        string             `json:"name,omitempty"`
	ParamIn     ParamIn            `json:"in,omitempty"`
	Description string             `json:"description,omitempty"`
	Required    bool               `json:"required,omitempty"`
	Deprecated  bool               `json:"deprecated,omitempty"`
	Example     any                `json:"example,omitempty"`
	Schema      *jsonschema.Schema `json:"schema,omitempty"`
	Style       ParamStyle         `json:"style,omitempty"`
	Explode     bool               `json:"explode,omitempty"`
}

type MediaType struct {
	Schema     *jsonschema.Schema `json:"schema,omitempty,omitzero"`
	ItemSchema *jsonschema.Schema `json:"itemSchema,omitempty"`
	Example    any                `json:"example,omitempty"`
}

type SecurityScheme struct {
	Type            SecurityType `json:"type"`
	Description     string       `json:"description"`
	Name            string       `json:"name"`
	In              ParamIn      `json:"in"`
	Scheme          string       `json:"scheme"`
	BearerFormat    string       `json:"bearerFormat"`
	Flows           *OAuthFlows  `json:"flows"`
	OpenIDConnecURL string       `json:"openIdConnectUrl"`
}

type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit"`
	Password          *OAuthFlow `json:"password"`
	ClientCredentials *OAuthFlow `json:"clientCredentials"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode"`
}

type OAuthFlow struct {
	AuhtorizationURL string            `json:"authorizationUrl"`
	TokenURL         string            `json:"tokenUrl"`
	RefreshURL       string            `json:"refreshUrl"`
	Scopes           map[string]string `json:"scopes"`
}
