package nuage

const (
	// ContentTypeJSON represents the standard JSON media type as defined in RFC 8259.
	//
	// This content type is used for typical JSON request and response bodies,
	// most commonly with HTTP methods such as POST and PUT where the full
	// resource representation is sent.
	ContentTypeJSON = "application/json"

	// ContentTypeMergePatch represents the JSON Merge Patch media type
	// as defined in RFC 7386.
	//
	// This content type is used for partial updates of a JSON document.
	// The patch document is a JSON object that is merged into the target
	// resource using specific merge semantics.
	ContentTypeMergePatch = "application/merge-patch+json"

	// ContentTypeJSONPatch represents the JSON Patch media type
	// as defined in RFC 6902.
	//
	// This content type is used for applying a sequence of operations
	// to a JSON document. The request body must be a JSON array of
	// patch operations.
	ContentTypeJSONPatch = "application/json-patch+json"

	// ContentTypeHTTPError defines the MIME type used for HTTP error responses
	// that follow the "Problem Details for HTTP APIs" specification (RFC 9457).
	//
	// This content type, "application/problem+json", indicates that the
	// response body contains a standardized JSON object describing the
	// error, including fields such as "type", "title", "status", "detail",
	// and optionally "instance". It is intended to provide clients with
	// machine-readable error information in a consistent format, which can
	// be programmatically processed or displayed to users.
	ContentTypeHTTPError = "application/problem+json"
)
