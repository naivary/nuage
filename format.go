package nuage

import (
	"encoding/json"
	"io"
)

type Formater interface {
	Decode(r io.Reader, v any) error
	Encode(w io.Writer, v any) error
}

var _ Formater = (*jsonFormater)(nil)

type jsonFormater struct{}

func (j jsonFormater) Decode(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(&v)
}

func (j jsonFormater) Encode(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(&v)
}
