package nuage

// HandlerFuncErr represents the primary request handler function signature
// used by the framework to implement REST API endpoints.
//
// It is a generic handler that:
//   - Accepts a strongly-typed request model (I)
//   - Returns a strongly-typed response model (O)
type HandlerFuncErr[I, O any] func(ctx *Context, requestModel I) (outputModel O, err error)
