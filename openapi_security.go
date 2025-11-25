package nuage

// ENUM(apiKey, http, mutualTLS, oauth2, openIdConnect)
type SecurityType string

type SecurityRequirement map[string][]string

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
