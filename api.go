package nuage

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type api struct {
	openAPI *OpenAPI

	mux *http.ServeMux
}

func NewAPI(root *OpenAPI) *api {
	return &api{
		openAPI: root,
		mux:     http.NewServeMux(),
	}
}

func Handle[I, O any](api *api, operation *Operation, handler HandlerFuncErr[I, O]) error {
	if !isStruct[I]() {
		return errors.New("non struct input type")
	}
	if !isStruct[O]() {
		return errors.New("non struct output type")
	}
	_, pattern, found := strings.Cut(operation.Pattern, " ")
	if !found {
		return fmt.Errorf("invalid pattern syntax: %s", operation.Pattern)
	}
	api.openAPI.Paths[pattern] = &pathItem{}
	api.mux.Handle(operation.Pattern, handler)
	return nil
}

func isStruct[T any]() bool {
	rtype := reflect.TypeFor[T]()
	return deref(rtype).Kind() == reflect.Struct
}

func deref(rtype reflect.Type) reflect.Type {
	if rtype.Kind() == reflect.Pointer {
		return rtype.Elem()
	}
	return rtype
}
