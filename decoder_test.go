package nuage

import (
	"math"
	"net/http"
	"strconv"
	"testing"
)

type pathParamsDecoder struct {
	P1 int8     `path:"p1"`
	P2 []string `path:"p2"`
}

func TestDecode(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	r.SetPathValue("p1", strconv.FormatInt(math.MaxInt8, 10))
	r.SetPathValue("p2", "v1,v2,v3")
	var input pathParamsDecoder
	err = Decode(r, &input)
	if err != nil {
		t.Errorf("decode: %v", err)
	}
	t.Logf("input: %v", input)
}
