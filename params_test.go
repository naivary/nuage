package nuage

import "testing"

type pathParamTest struct {
	f1 int `path:"f1"`

	f2 int `path:"f2,deprecated"`

	// f3 is this and that and this should be the description of the parameter.
	f3 int `path:"f3,Cookie"`
}

func TestPathParams(t *testing.T) {
	params, err := pathParams[pathParamTest]()
	if err != nil {
		t.Fatalf("path params: %v", err)
	}
	for _, param := range params {
		t.Log(param)
	}
}
