package nuage

import (
	"net/http"
	"net/url"
	"testing"
)

type decodeTest struct {
	Q1 *string `query:"q1" json:"-"`
	Q2 string  `query:"q2" json:"-"`
}

func TestDecode(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "", nil)
	r.URL.RawQuery = url.Values{
		"q1": []string{"e1"},
		"q2": []string{"e2"},
	}.Encode()
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	input := decodeTest{}
	err = Decode(r, &input)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	t.Log(input)
}
