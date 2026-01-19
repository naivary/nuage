package nuage

import (
	"context"
	"time"
)

var _ context.Context = (*Context)(nil)

type Context struct {
	ctx context.Context
}

func NewCtx() *Context {
	return &Context{
		ctx: context.Background(),
	}
}

func (c *Context) Deadline() (time.Time, bool) {
	return c.ctx.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Context) Err() error {
	return c.ctx.Err()
}

func (c *Context) Value(key any) any {
	return c.ctx.Value(key)
}
