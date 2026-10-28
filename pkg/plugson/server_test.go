package plugson

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlugsonServer(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create Plugson server: %v", err)
	}

	req := httptest.NewRequest("GET", "/api/config", nil)
	w := httptest.NewRecorder()

	server.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected HTTP 200 for /api/config, got %d", w.Code)
	}

	reqHTML := httptest.NewRequest("GET", "/", nil)
	wHTML := httptest.NewRecorder()

	server.ServeHTTP(wHTML, reqHTML)

	if wHTML.Code != http.StatusOK {
		t.Errorf("expected HTTP 200 for index.html, got %d", wHTML.Code)
	}
}
