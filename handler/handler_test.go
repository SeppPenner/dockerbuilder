package handler

import (
	"crypto/sha1"
	"github.com/brocaar/dockerbuilder/config"
	"github.com/brocaar/dockerbuilder/worker"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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
	r, _ := http.NewRequest("POST", "http://example.com/github/hook", strings.NewReader(""))
	r.Header.Add("X-Github-Event", "ping")

	handler := &GitHubHandler{config: &config.Configuration{}}
	handler.Hook(w, r)

	if w.Code != 200 {
		t.Errorf("expected: 200, got: %d", w.Code)
	}
}

// Test GitHub create tag event
func TestGitHubHandlerCreateTag(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "http://example.com/github/hook", strings.NewReader(githubCreateTagPayload))
	r.Header.Add("X-Github-Event", "create")

	config := &config.Configuration{
		DockerIndexNamespace: "dockeruser",
		CleanupContainer:     true,
	}
	var taskQueue = make(worker.TaskQueue, 1)

	handler := NewGitHubHandler(taskQueue, config)
	handler.Hook(w, r)

	if w.Code != 200 {
		t.Errorf("exected: 200, got: %d", w.Code)
	}

	select {
	case workerTask := <-taskQueue:
		if workerTask.Revision != "v0.1.2" {
			t.Errorf("expected revision: v0.1.2, got: %s", workerTask.Revision)
		}

		if workerTask.DockerIndexNamespace != "dockeruser" {
			t.Errorf("expected docker index namespace: dockeruser, got: %s", workerTask.DockerIndexNamespace)
		}

		if workerTask.Repository.Name != "dockerbuilder" {
			t.Errorf("expected repo name: dockerbuilder, got: %s", workerTask.Repository.Name)
		}

		if workerTask.Repository.Owner != "brocaar" {
			t.Errorf("expected repo owner: brocaar, got: %s", workerTask.Repository.Owner)
		}

		if workerTask.CleanupContainer != true {
			t.Errorf("expected cleanup container: true, got: %t", workerTask.CleanupContainer)
		}
	case <-time.After(1):
		t.Error("exepected a worker-task in the channel, but it was empty!")
	}
}

var githubCreateTagPayload = `{
	"ref": "v0.1.2",
	"ref_type": "tag",
	"repository": {
		"name": "dockerbuilder",
		"full_name": "brocaar/dockerbuilder",
		"owner": {
			"login": "brocaar"
		}
	}
}`
