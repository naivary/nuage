package nuage

import (
	"log/slog"
	"net/http"

	"github.com/naivary/nuage/openapi"
)

const _headerKeyContentType = "Content-Type"

type Responser interface {
	StatusCode() int
	Description() string
}

type HandlerFuncErr[I, O any] func(r *http.Request, input *I) (*O, error)

type endpoint[I, O any] struct {
	handler HandlerFuncErr[I, O]
	doc     *openapi.Operation
	logger  *slog.Logger
	formats map[string]Formater
}

func (e endpoint[I, O]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	format := r.Header.Get(_headerKeyContentType)
	formater, isSupportedFormat := e.formats[format]
	if !isSupportedFormat {
		// unssuported format
		e.logger.Error("format not suypported", "format", format)
		return
	}
	var input I
	err := decodeParams(r, &input)
	if err != nil {
		e.logger.Error(err.Error())
		return
	}
	if err := formater.Decode(r.Body, &input); err != nil {
		// bad request internal error of decoding format
		e.logger.Error(err.Error())
		return
	}
}
