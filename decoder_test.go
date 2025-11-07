package nuage

import (
	"math"
	"net/http"
	"strconv"
	"testing"
)

type paramsDecodeT1 struct {
	P1 int8 `path:"p1"`
}

func TestDecode(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	r.SetPathValue("p1", strconv.FormatInt(math.MaxInt8+1, 10))
	var input paramsDecodeT1
	err = Decode(r, &input)
	if err != nil {
		t.Errorf("decode: %v", err)
	}
	t.Logf("input: %v", input)
}
