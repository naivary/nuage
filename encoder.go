package nuage

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

const ContentTypeJSON = "application/json"

const (
	_headerKeySetCookie   = "Set-Cookie"
	_headerKeyContentType = "Content-Type"
)

func encode[O any](w http.ResponseWriter, status int, output O) error {
	if !isStruct[O]() {
		return fmt.Errorf("encode: only structs can be used as response")
	}
	rvalue := derefValue(reflect.ValueOf(output))
	rtype := deref(rvalue.Type())
	for i := range rtype.NumField() {
		field := rtype.Field(i)
		if !field.IsExported() {
			// unexported fields will be ignored for decoding
			continue
		}
		fieldValue := rvalue.Field(i)
		cookie, isCookie := fieldValue.Interface().(*http.Cookie)
		if isCookie {
			jsonTagValue := field.Tag.Get("json")
			if jsonTagValue != "-" {
				return fmt.Errorf(`encode: make sure to add json:"-" to your cookie response value`)
			}
			if err := cookie.Valid(); err != nil {
				return fmt.Errorf("encode: invalid cookie %v", cookie)
			}
			w.Header().Set(_headerKeySetCookie, cookie.String())
			continue
		}
		opts, err := parseParamTagOpts(field)
		if errors.Is(err, errTagNotFound) {
			continue
		}
		if err != nil {
			return err
		}
		switch opts.tagKey {
		case _tagKeyHeader:
			w.Header().Set(opts.name, fieldValue.String())
		}
	}
	w.Header().Set(_headerKeyContentType, ContentTypeJSON)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(output); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}
