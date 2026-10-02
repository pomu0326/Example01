package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example01/backend/internal/todo"
)

func newTestServer(t *testing.T) http.Handler {
	t.Helper()
	s, err := todo.OpenSQLite(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return NewHandler(s)
}

func do(t *testing.T, h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func decode[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var v T
	if err := json.NewDecoder(rec.Body).Decode(&v); err != nil {
		t.Fatalf("decode %q: %v", rec.Body.String(), err)
	}
	return v
}

func TestCRUDFlow(t *testing.T) {
	h := newTestServer(t)

	rec := do(t, h, "POST", "/api/todos", `{"title":"write tests"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, body %s", rec.Code, rec.Body)
	}
	created := decode[todo.Todo](t, rec)
	if created.Title != "write tests" || created.Completed {
		t.Fatalf("created = %+v", created)
	}
	path := "/api/todos/" + jsonInt(created.ID)

	rec = do(t, h, "GET", "/api/todos", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d", rec.Code)
	}
	if list := decode[[]todo.Todo](t, rec); len(list) != 1 || list[0].ID != created.ID {
		t.Fatalf("list = %+v", list)
	}

	rec = do(t, h, "PATCH", path, `{"completed":true}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("patch status = %d, body %s", rec.Code, rec.Body)
	}
	if got := decode[todo.Todo](t, rec); !got.Completed || got.Title != "write tests" {
		t.Fatalf("patched = %+v", got)
	}

	rec = do(t, h, "GET", path, "")
	if rec.Code != http.StatusOK || !decode[todo.Todo](t, rec).Completed {
		t.Fatalf("get after patch status = %d", rec.Code)
	}

	rec = do(t, h, "DELETE", path, "")
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d", rec.Code)
	}

	rec = do(t, h, "GET", path, "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("get after delete status = %d", rec.Code)
	}
}

func TestErrors(t *testing.T) {
	h := newTestServer(t)
	tests := []struct {
		name, method, path, body string
		want                     int
	}{
		{"empty title", "POST", "/api/todos", `{"title":"  "}`, http.StatusBadRequest},
		{"malformed json", "POST", "/api/todos", `{`, http.StatusBadRequest},
		{"unknown field", "POST", "/api/todos", `{"name":"x"}`, http.StatusBadRequest},
		{"non-numeric id", "GET", "/api/todos/abc", "", http.StatusBadRequest},
		{"zero id", "GET", "/api/todos/0", "", http.StatusBadRequest},
		{"get missing", "GET", "/api/todos/42", "", http.StatusNotFound},
		{"patch missing", "PATCH", "/api/todos/42", `{"completed":true}`, http.StatusNotFound},
		{"delete missing", "DELETE", "/api/todos/42", "", http.StatusNotFound},
		{"wrong method", "PUT", "/api/todos/1", `{}`, http.StatusMethodNotAllowed},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := do(t, h, tc.method, tc.path, tc.body)
			if rec.Code != tc.want {
				t.Fatalf("status = %d, want %d (body %s)", rec.Code, tc.want, rec.Body)
			}
			if tc.want != http.StatusMethodNotAllowed {
				if body := decode[map[string]string](t, rec); body["error"] == "" {
					t.Fatalf("missing error message: %v", body)
				}
			}
		})
	}
}

func TestHealthz(t *testing.T) {
	rec := do(t, newTestServer(t), "GET", "/healthz", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("healthz status = %d", rec.Code)
	}
}

func jsonInt(n int64) string {
	b, _ := json.Marshal(n)
	return string(b)
}
