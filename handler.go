package nuage

import "context"

type HandlerFuncErr[I, O any] func(ctx context.Context, input I) (O, error)
