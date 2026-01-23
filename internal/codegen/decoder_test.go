package codegen_test

import (
	"testing"

	"github.com/naivary/nuage/internal/codegen"
)

func TestGenDecoder(t *testing.T) {
	err := codegen.GenDecoder([]string{"./testdata"})
	if err != nil {
		t.Errorf("codegen: %v", err)
	}
}
