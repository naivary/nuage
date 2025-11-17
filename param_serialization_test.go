package nuage

import (
	"reflect"
	"slices"
	"testing"

	"github.com/naivary/nuage/openapi"
)

func TestSerializePathParam(t *testing.T) {
	m := reflect.TypeFor[map[string]string]()
	s := reflect.TypeFor[[]string]()
	tests := []struct {
		name     string
		v        string
		style    openapi.Style
		explode  bool
		typ      reflect.Type
		expected []string
	}{
		{
			name:     "simple primitive",
			v:        "param",
			style:    openapi.StyleSimple,
			typ:      reflect.TypeFor[string](),
			expected: []string{"param"},
		},
		{
			name:     "simple primitive explode",
			v:        "param",
			style:    openapi.StyleSimple,
			typ:      reflect.TypeFor[string](),
			explode:  true,
			expected: []string{"param"},
		},
		{
			name:     "simple array",
			v:        "e1,e2,e3",
			style:    openapi.StyleSimple,
			typ:      s,
			expected: []string{"e1", "e2", "e3"},
		},
		{
			name:     "simple array explode",
			v:        "e1,e2,e3",
			style:    openapi.StyleSimple,
			typ:      s,
			explode:  true,
			expected: []string{"e1", "e2", "e3"},
		},
		{
			name:     "simple map",
			v:        "k1,v1,k2,v2",
			style:    openapi.StyleSimple,
			typ:      m,
			expected: []string{"k1", "v1", "k2", "v2"},
		},
		{
			name:     "simple map explode",
			v:        "k1=v1,k2=v2",
			explode:  true,
			style:    openapi.StyleSimple,
			typ:      m,
			expected: []string{"k1", "v1", "k2", "v2"},
		},
		{
			name:     "label primitive",
			v:        ".p",
			style:    openapi.StyleLabel,
			typ:      reflect.TypeFor[string](),
			expected: []string{"p"},
		},
		{
			name:     "label primitive explode",
			v:        ".p",
			style:    openapi.StyleLabel,
			explode:  true,
			typ:      reflect.TypeFor[string](),
			expected: []string{"p"},
		},
		{
			name:     "label array",
			v:        ".e1,e2,e3",
			style:    openapi.StyleLabel,
			typ:      s,
			expected: []string{"e1", "e2", "e3"},
		},
		{
			name:     "label array explode",
			v:        ".e1.e2.e3",
			style:    openapi.StyleLabel,
			explode:  true,
			typ:      s,
			expected: []string{"e1", "e2", "e3"},
		},
		{
			name:     "label map",
			v:        ".k1,v1,k2,v2",
			style:    openapi.StyleLabel,
			typ:      m,
			expected: []string{"k1", "v1", "k2", "v2"},
		},
		{
			name:     "label map explode",
			v:        ".k1=v1.k2=v2",
			style:    openapi.StyleLabel,
			explode:  true,
			typ:      m,
			expected: []string{"k1", "v1", "k2", "v2"},
		},
		{
			name:     "matrix primitive",
			v:        ";p1=t",
			style:    openapi.StyleMatrix,
			typ:      reflect.TypeFor[string](),
			expected: []string{"t"},
		},
		{
			name:     "matrix primitive explode",
			v:        ";p1=t",
			style:    openapi.StyleMatrix,
			explode:  true,
			typ:      reflect.TypeFor[string](),
			expected: []string{"t"},
		},
		{
			name:     "matrix array",
			v:        ";p1=e1,e2,e3",
			style:    openapi.StyleMatrix,
			typ:      s,
			expected: []string{"e1", "e2", "e3"},
		},
		{
			name:     "matrix array explode",
			v:        ";p1=e1;p2=e2;p3=e3",
			style:    openapi.StyleMatrix,
			explode:  true,
			typ:      s,
			expected: []string{"e1", "e2", "e3"},
		},
		{
			name:     "matrix map",
			v:        ";p1=k1,v1,k2,v2",
			style:    openapi.StyleMatrix,
			typ:      m,
			expected: []string{"k1", "v1", "k2", "v2"},
		},
		{
			name:     "matrix map explode",
			v:        ";k1=v1;k2=v2",
			style:    openapi.StyleMatrix,
			explode:  true,
			typ:      m,
			expected: []string{"k1", "v1", "k2", "v2"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			values, err := serializePathParam(tc.v, tc.typ, tc.style, tc.explode)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			if !slices.Equal(values, tc.expected) {
				t.Errorf("values slice not equal. Got: %v. Want: %v", values, tc.expected)
			}
			t.Logf("values: %v", values)
		})
	}
}
