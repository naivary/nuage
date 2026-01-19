package nuage

import "net/http"

type Nuage struct {
	mux *http.ServeMux
}

func New() (*Nuage, error) {
	return &Nuage{
		mux: http.NewServeMux(),
	}, nil
}
