package api

import (
	"encoding/json"
	"example.com/pz3-http/internal/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateTask(t *testing.T) {
	cases := []struct {
		name, body, contentType string
		want                    int
	}{
		{"valid", `{"title":"Купить молоко"}`, "application/json", 201},
		{"charset", `{"title":"Test"}`, "application/json; charset=utf-8", 201},
		{"bad json", `{`, "application/json", 400},
		{"empty", `{}`, "application/json", 422},
		{"short", `{"title":"аб"}`, "application/json", 422},
		{"three characters", `{"title":"абв"}`, "application/json", 201},
		{"140 characters", `{"title":"` + strings.Repeat("я", 140) + `"}`, "application/json", 201},
		{"141 characters", `{"title":"` + strings.Repeat("я", 141) + `"}`, "application/json", 422},
		{"wrong content type", `{"title":"Test"}`, "text/plain", 400},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := storage.NewMemoryStore()
			h := NewHandlers(store)
			req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", tc.contentType)
			rec := httptest.NewRecorder()
			h.CreateTask(rec, req)
			if rec.Code != tc.want {
				t.Fatalf("status=%d want=%d body=%s", rec.Code, tc.want, rec.Body.String())
			}
			if !strings.Contains(rec.Header().Get("Content-Type"), "application/json") {
				t.Fatal("missing JSON content type")
			}
			if tc.want == 201 {
				var got storage.Task
				if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
					t.Fatal(err)
				}
				if got.ID != 1 || got.Done || got.Title == "" {
					t.Fatalf("unexpected task: %+v", got)
				}
				if len(store.List()) != 1 {
					t.Fatal("task not stored")
				}
			} else if len(store.List()) != 0 {
				t.Fatal("invalid task stored")
			}
		})
	}
}
func TestGetTask(t *testing.T) {
	store := storage.NewMemoryStore()
	created := store.Create("Test task")
	h := NewHandlers(store)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks/{id}", h.GetTask)
	for _, tc := range []struct {
		path string
		want int
	}{
		{"/tasks/1", 200}, {"/tasks/9999", 404}, {"/tasks/abc", 400}, {"/tasks/0", 400},
	} {
		t.Run(tc.path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if rec.Code != tc.want {
				t.Fatalf("status=%d want=%d", rec.Code, tc.want)
			}
			if tc.want == 200 {
				var got storage.Task
				if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
					t.Fatal(err)
				}
				if got != *created {
					t.Fatalf("got %+v want %+v", got, *created)
				}
			}
		})
	}
}
