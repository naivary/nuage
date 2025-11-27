package nuagetest

import (
	"io"
	"net/http"
)

func NewRequest(method, url string, body io.Reader) *http.Request {
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	return r
}

func AddHeaders(r *http.Request, headers map[string]string) {
	for key, value := range headers {
		r.Header.Add(key, value)
	}
}

func AddQueryParams(r *http.Request, params map[string]string) {
	q := r.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	r.URL.RawQuery = q.Encode()
}
