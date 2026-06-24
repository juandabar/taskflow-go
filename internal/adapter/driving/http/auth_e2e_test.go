package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	httpAdapter "github.com/juandabar/taskflow-go/internal/adapter/driving/http"
	"github.com/juandabar/taskflow-go/internal/infrastructure/container"
	"github.com/juandabar/taskflow-go/internal/infrastructure/database"
)

const testJWTSecret = "test-secret-key-for-e2e-tests-32chars"

func setupRouter(t *testing.T) http.Handler {
	t.Helper()

	db, err := database.NewSQLiteConnection(":memory:")
	if err != nil {
		t.Fatalf("failed to setup database: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	c := container.NewContainer(db, testJWTSecret)
	return httpAdapter.NewRouter(c.AuthHandler, c.UserHandler, testJWTSecret)
}

func TestRegisterEndpoint(t *testing.T) {
	router := setupRouter(t)

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name: "success",
			body: map[string]string{
				"name":     "Juan",
				"email":    "juan@example.com",
				"password": "securepassword123",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "missing fields",
			body:       map[string]string{"name": "Juan"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "duplicate email",
			body: map[string]string{
				"name":     "Juan",
				"email":    "juan@example.com",
				"password": "securepassword123",
			},
			wantStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			r := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(bodyBytes))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d - body: %s", tt.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}
