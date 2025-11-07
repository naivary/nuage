package nuage

import "context"

type Validator interface {
	// Valid checks the object and returns any
	// problems. If len(problems) == 0 then
	// the object is valid.
	Validate(ctx context.Context) (problems map[string]string)
}
