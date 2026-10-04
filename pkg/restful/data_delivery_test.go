// SPDX-License-Identifier: Apache-2.0
package restful

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/a1"
	"gerrit.o-ran-sc.org/r/ric-plt/a1/pkg/resthooks"
	"github.com/stretchr/testify/assert"
)

func TestMalformedDataDeliveryHTTP(t *testing.T) {
	a1.Init()
	r := &Restful{rh: &resthooks.Resthook{}}
	handler := r.setupHandler().Serve(nil)
	for _, body := range []string{
		`{}`, `[]`, `null`, `{"payload":"data"}`,
		`{"job":1,"payload":"data"}`, `{"job":"1"}`,
	} {
		t.Run(body, func(t *testing.T) {
			t.Parallel()
			request := httptest.NewRequest(http.MethodPost, "/data-delivery", strings.NewReader(body))
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			assert.NotPanics(t, func() { handler.ServeHTTP(recorder, request) })
			assert.Equal(t, http.StatusNotFound, recorder.Code)
		})
	}
}
