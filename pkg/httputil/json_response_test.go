package httputil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	type testCase struct {
		name   string
		assert func(t *testing.T)
	}

	tc := []testCase{
		{
			name: "given no body and no status code response should have status code ok",
			assert: func(t *testing.T) {
				jr := NewJsonResponse()
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/", nil)

				jr.Response(w, r)

				assert.Equal(t, http.StatusOK, w.Code)
				assert.Equal(t, w.Header().Get("Content-Type"), "application/json")
			},
		},
		{
			name: "given a status code response should contain the same code",
			assert: func(t *testing.T) {
				jr := NewJsonResponse(func(jr *JsonResponse) {

				})
				w := httptest.NewRecorder()
				r := httptest.NewRequest("GET", "/", nil)

				jr.Response(w, r)
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(t)
		})
	}
}
