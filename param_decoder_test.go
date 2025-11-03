package nuage

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

func TestPathParamDecoder(t *testing.T) {
	type v struct {
		P1 int    `path:"p1,example=1289397"`
		P2 string `path:"p2,example"`
		P3 string `path:"p3,optional,example"`
	}

	r, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	r.SetPathValue("p1", "123")
	r.SetPathValue("p2", "test-string-value")

	input := v{}
	err = DecodePath(r, &input)
	if err != nil {
		t.Fatalf("decode path: %v", err)
	}
	t.Log(input)
}

func TestHeaderParamDecoder(t *testing.T) {
	type v struct {
		H1 int8    `header:"h1"`
		H2 string `header:"h2"`
		H3 string `header:"h3,required"`
	}

	r, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	r.Header.Set("h1", strconv.FormatInt(math.MaxInt8+1, 10))
	r.Header.Set("h2", "value-for-h2")
	r.Header.Set("h3", "value-for-h3")

	input := v{}
	err = DecodeHeader(r, &input)
	if err != nil {
		t.Fatalf("decode header: %v", err)
	}
	t.Log(input)
}

func TestQueryParamDecoder(t *testing.T) {
	type v struct {
		Q1 int    `query:"q1"`
		Q2 string `query:"q2"`
		Q3 string `query:"q3,required"`
	}

	r, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	q := url.Values{}
	q.Set("q1", "12398")
	q.Set("q2", "query-2")
	q.Set("q3", "query-3")
	r.URL.RawQuery = q.Encode()

	input := v{}
	err = DecodeQuery(r, &input)
	if err != nil {
		t.Fatalf("decode query: %v", err)
	}
	t.Log(input)
}

func TestCookieParamDecoder(t *testing.T) {
	type v struct {
		C1 int    `cookie:"c1"`
		C2 string `cookie:"c2"`
		C3 string `cookie:"c3"`
	}

	r, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	r.AddCookie(&http.Cookie{
		Name:  "c3",
		Value: "dasjhdsa",
		Path:  "/",
	})

	input := v{}
	err = DecodeCookie(r, &input)
	if err != nil {
		t.Fatalf("decode cookie: %v", err)
	}
	t.Log(input)
}
