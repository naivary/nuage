package nuage

type HandlerFuncErr[I, O any] func(ctx *Context, requestModel I) (O, error)
