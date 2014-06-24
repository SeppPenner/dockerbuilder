package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test GitHub ping event returns 200.
func TestGitHubHandlerPing(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "http://example.com/github/hook", nil)
	if err != nil {
		t.Errorf("creating request failed: %s", err)
	}
	r.Header.Add("X-Github-Event", "ping")

	handler := &GitHubHandler{}
	handler.Hook(w, r)

	if w.Code != 200 {
		t.Errorf("expected: 200, got: %d", w.Code)
	}
}
