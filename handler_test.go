package nuage

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type requestTest struct {
	P1 string `path:"p1" json:"-"`
}

type responseTest struct {
	Name string `json:"name"`

	H1 string `header:"X-H1" json:"-"`

	C1 *http.Cookie `json:"-"`
}

func TestHandlerFuncErr(t *testing.T) {
	hl := HandlerFuncErr[requestTest, responseTest](func(r *http.Request, w http.ResponseWriter, input *requestTest) (*responseTest, error) {
		return &responseTest{
			Name: "test",
			H1:   "something",
			C1: &http.Cookie{
				Name: "c1",
				Value: "testvalue",
			},
		}, nil
	})
	r, err := http.NewRequest(http.MethodGet, "", nil)
	r.SetPathValue("p1", "testp1")
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	w := httptest.NewRecorder()
	hl.ServeHTTP(w, r)
	t.Log(w.Body.String())
	t.Log(w.Header())
}
