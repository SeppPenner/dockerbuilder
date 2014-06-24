// Package handler contains the handlers for the webhook requests.
package handler

import (
	"github.com/brocaar/dockerbuilder/worker"
	"net/http"
)

type GitHubHandler struct {
	taskQueue worker.TaskQueue
}

// NewGitHubHandler returns a new instance of GitHubHandler.
func NewGitHubHandler(taskQueue worker.TaskQueue) *GitHubHandler {
	return &GitHubHandler{
		taskQueue: taskQueue,
	}
}

// Hook is a HTTP handler for webhook requests by GitHub.
func (h *GitHubHandler) Hook(w http.ResponseWriter, r *http.Request) {
	// github ping
	if r.Header.Get("X-Github-Event") == "ping" {
		w.WriteHeader(http.StatusOK)
	}
}
