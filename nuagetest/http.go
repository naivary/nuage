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
