package nuage

import "testing"

type pathParamTest struct {
	f1 int `path:"f1"`

	f2 int `path:"f2,deprecated"`

	f3 int `path:"f3"`
}

func TestPathParams(t *testing.T) {
	params, err := pathParams[pathParamTest]()
	if err != nil {
		t.Fatalf("path params: %v", err)
	}
	t.Log(params[0])
}
