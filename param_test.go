package nuage

import (
	"encoding/json"
	"testing"
)

type paramsTestType struct {
	P1 string `path:"p1" json:"-"`

	H1 string `header:"X-H1" json:"-"`

	Q1 string `query:"q1" json:"-"`

	C1 string `cookie:"c1" json:"-"`
}

func TestParamsFor(t *testing.T) {
	params, err := ParamSpecsFor[paramsTestType]()
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
