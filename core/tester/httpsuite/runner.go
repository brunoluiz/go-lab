package httpsuite

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func RunSequence(ctx context.Context, t *testing.T, h http.Handler, steps []Step) {
	kv := NewKV()
	for _, s := range steps {
		vars := map[string]any{}
		if s.PreHook != nil {
			vars = s.PreHook(t, kv)
		}

		payload, err := FromJSON(s.RequestTemplatePath, vars)
		require.NoError(t, err)
		t.Run(payload.Name, func(_ *testing.T) {
			req, err := payload.Request(ctx)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)
			s.PostHook(t, kv, w)
			t.Logf("status: %d", w.Code)
			t.Logf("response: %s", w.Body.String())
		})
	}
}
