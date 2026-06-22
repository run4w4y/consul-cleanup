package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestAuthStaticToken(t *testing.T) {
	tests := []struct {
		name        string
		authHeader  string
		wantStatus  int
		wantHandler bool
	}{
		{
			name:       "missing authorization header",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "wrong scheme",
			authHeader: "Basic secret",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "wrong token",
			authHeader: "Bearer wrong",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:        "valid token",
			authHeader:  "Bearer secret",
			wantHandler: true,
		},
		{
			name:        "case insensitive scheme",
			authHeader:  "bearer secret",
			wantHandler: true,
		},
		{
			name:        "extra whitespace",
			authHeader:  "Bearer   secret",
			wantHandler: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			handlerCalled := false
			handler := AuthStaticToken("secret")(func(c echo.Context) error {
				handlerCalled = true
				return c.NoContent(http.StatusNoContent)
			})

			err := handler(ctx)
			if tt.wantHandler {
				if err != nil {
					t.Fatalf("handler returned unexpected error: %v", err)
				}
				if !handlerCalled {
					t.Fatal("expected wrapped handler to be called")
				}
				if rec.Code != http.StatusNoContent {
					t.Fatalf("response status = %d, want %d", rec.Code, http.StatusNoContent)
				}
				return
			}

			if handlerCalled {
				t.Fatal("wrapped handler should not be called")
			}
			httpErr, ok := err.(*echo.HTTPError)
			if !ok {
				t.Fatalf("error = %T, want *echo.HTTPError", err)
			}
			if httpErr.Code != tt.wantStatus {
				t.Fatalf("error status = %d, want %d", httpErr.Code, tt.wantStatus)
			}
		})
	}
}
