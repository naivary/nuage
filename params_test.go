package nuage

import (
	"encoding/json"
	"testing"
)

type paramsTestType struct {
	p1 string `path:"p1"`
	h1 string `header:"X-H1"`
	q1 string `query:"q1"`
	c1 string `cookie:"c1"`
}

func TestParamsFor(t *testing.T) {
	params, err := paramsFor[paramsTestType]()
	if err != nil {
		t.Fatalf("params for: %v", err)
	}
	if len(params) != 4 {
		t.Errorf("params count wrong. Got: %d. Want: 4", len(params))
	}
	for _, param := range params {
		data, err := json.Marshal(param)
		if err != nil {
			t.Errorf("json marshal: %v", err)
		}
		t.Logf("param: %s", data)
	}
}
