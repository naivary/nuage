package nuage

import (
	"encoding/json"
	"io"
)

type Formater interface {
	Decode(r io.Reader, value any) error
	Encode(w io.Writer, value any) error
}

var _ Formater = (*jsonFormater)(nil)

type jsonFormater struct{}

func (j jsonFormater) Decode(r io.Reader, value any) error {
	return json.NewDecoder(r).Decode(&value)
}

func (j jsonFormater) Encode(w io.Writer, value any) error {
	return json.NewEncoder(w).Encode(&value)
}
