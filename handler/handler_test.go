package handler

import (
	"crypto/sha1"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetMacString(t *testing.T) {
	mac := getMacString(sha1.New, []byte("verysecret"), []byte("hello world"))
	if mac != "3c24d84f47692d0b4aef2a5002e0bb6beb25dc75" {
		t.Errorf("expected: 3c24d84f47692d0b4aef2a5002e0bb6beb25dc75, got: %s", mac)
	}
}

func TestCheckGitHubMac(t *testing.T) {
	testTable := []struct {
		secret    string
		signature string
		body      string
		expected  bool
	}{
		{
			// valid
			"verysecret",
			"sha1=3c24d84f47692d0b4aef2a5002e0bb6beb25dc75",
			"hello world",
			true,
		}, {
			// body invalid
			"verysecret",
			"sha1=3c24d84f47692d0b4aef2a5002e0bb6beb25dc75",
			"Hello World!",
			false,
		}, {
			// token invalid
			"verysecret",
			"sha1=3c24d84f47692d0b4aef2a5002e0bb6beb25dc76",
			"hello world",
			false,
		}, {
			// algorithm invalid
			"verysecret",
			"sha512=3c24d84f47692d0b4aef2a5002e0bb6beb25dc75",
			"hello world",
			false,
		}, {
			// secret invalid
			"VerySecret",
			"sha1=3c24d84f47692d0b4aef2a5002e0bb6beb25dc75",
			"hello world",
			false,
		},
	}

	for i, test := range testTable {
		testResult := checkGitHubMac([]byte(test.secret), test.signature, []byte(test.body))
		if testResult != test.expected {
			t.Errorf("test %d: expected: %t, got: %t", i, test.expected, testResult)
		}
	}
}

// Test GitHub ping event returns 200.
func TestGitHubHandlerPing(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("POST", "http://example.com/github/hook", strings.NewReader(""))
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
