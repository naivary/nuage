package nuage

import (
	"net/http"
	"testing"
)

type pathParamsDecoder struct {
	PSimplePrimitive string            `path:"pSimplePrimitive"`
	PSimpleArray     []string          `path:"pSimpleArray"`
	PSimpleObject    map[string]string `path:"pSimpleObject"`

	PLabelPrimitive string            `path:"pLabelPrimitive,Label"`
	PLabelArray     []string          `path:"pLabelArray,Label"`
	PLabelObject    map[string]string `path:"pLabelObject,Label"`
}

func TestDecode(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	r.SetPathValue("pSimplePrimitive", "test")
	r.SetPathValue("pSimpleArray", "e1,e2,e3")
	r.SetPathValue("pSimpleObject", "k1,v1,k2,v2")

	r.SetPathValue("pLabelPrimitive", ".t")
	r.SetPathValue("pLabelArray", ".e1,e2,e3")
	r.SetPathValue("pLabelObject", ".k1,v1,k2,v2")

	var input pathParamsDecoder
	err = Decode(r, &input)
	if err != nil {
		t.Errorf("decode: %v", err)
	}
	t.Logf("input: %v", input)
}
