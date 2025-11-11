package nuage

import (
	"reflect"
	"testing"
)

type paramTagOptsTest struct {
	Q1 string `query:"q1,explode,deprecated,required" paramexample:"e1" paramstyle:"DeepObject"`
}

func TestParseParamTagOpts(t *testing.T) {
	rvalue := reflect.TypeFor[paramTagOptsTest]()
	for i := range rvalue.NumField() {
		field := rvalue.Field(i)
		opts, err := parseParamTagOpts(_tagKeyQuery, field)
		if err != nil {
			t.Fatalf("parse param tag opts: %v", err)
		}
		t.Log(opts)
	}
}
