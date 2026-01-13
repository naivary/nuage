package nuage

import "net/url"

const (
	OpenAPIVersion    = "3.2.0"
	JSONSchemaDialect = "https://json-schema.org/draft/2020-12/json-schema-core.html"
)

type OpenAPI struct {
	Version           string    `json:"openapi"`
	Self              string    `json:"$self"`
	Info              *Info     `json:"info"`
	JSONSchemaDialect string    `json:"jsonSchemaDialect,omitempty"`
	Servers           []*Server `json:"servers,omitempty"`
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
