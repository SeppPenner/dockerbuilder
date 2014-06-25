// Package handler contains the handlers for the webhook requests.
package handler

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"github.com/brocaar/dockerbuilder/config"
	"github.com/brocaar/dockerbuilder/repository"
	"github.com/brocaar/dockerbuilder/worker"
	"hash"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type gitHubRepositoryOwner struct {
	Login string `json:"login"`
}

type gitHubRepository struct {
	Name     string                 `json:"name"`
	FullName string                 `json:"full_name"`
	Owner    *gitHubRepositoryOwner `json:"owner"`
}

type githubCreateEvent struct {
	Ref        string            `json:"ref"`
	RefType    string            `json:"ref_type"`
	Repository *gitHubRepository `json:"repository"`
}

type GitHubHandler struct {
	taskQueue worker.TaskQueue
	config    *config.Configuration
}

// NewGitHubHandler returns a new instance of GitHubHandler.
func NewGitHubHandler(taskQueue worker.TaskQueue, c *config.Configuration) *GitHubHandler {
	return &GitHubHandler{
		taskQueue: taskQueue,
		config:    c,
	}
}

// Hook is a HTTP handler for webhook requests by GitHub.
func (h *GitHubHandler) Hook(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("could not read request body: %s\n", err)
		return
	}

	// check the signature
	if len(h.config.GitHubSecret) > 0 && checkGitHubMac(h.config.GitHubSecret, r.Header.Get("X-Hub-Signature"), requestBody) != true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// handle ping event
	if r.Header.Get("X-Github-Event") == "ping" {
		h.handlePing(w, r)
	}

	// handle create event
	if r.Header.Get("X-Github-Event") == "create" {
		h.handleCreate(w, r, requestBody)
	}
}

// handlePing handles a GitHub ping event
func (h *GitHubHandler) handlePing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// handleCreate handles a GitHub create event
func (h *GitHubHandler) handleCreate(w http.ResponseWriter, r *http.Request, requestBody []byte) {
	// unmarshall JSON requestBody
	var eventPayload githubCreateEvent
	err := json.Unmarshal(requestBody, &eventPayload)
	if err != nil {
		log.Printf("could not parse JSON payload. Error: %s, JSON body: %s\n", err, requestBody)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// currently we ignore events other than "tag" and we just return a 200
	if eventPayload.RefType != "tag" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// create new worker task
	h.taskQueue <- &worker.WorkerTask{
		Revision:             eventPayload.Ref,
		DockerIndexNamespace: h.config.DockerIndexNamespace,
		Repository:           repository.NewRepository(repository.HostGitHub, eventPayload.Repository.Owner.Login, eventPayload.Repository.Name, repository.ScmGit),
		CleanupContainer:     h.config.CleanupContainer,
	}
	w.WriteHeader(http.StatusOK)
}

func getMacString(hashAlgo func() hash.Hash, key, message []byte) string {
	mac := hmac.New(hashAlgo, key)
	mac.Write(message)
	macBytes := mac.Sum(nil)
	return hex.EncodeToString(macBytes)
}

func checkGitHubMac(key []byte, gitHubSignature string, message []byte) bool {
	signatureParts := strings.Split(gitHubSignature, "=")
	var algo func() hash.Hash

	// more algorithms can be added later if needed
	if signatureParts[0] == "sha1" {
		algo = sha1.New
	} else {
		return false
	}

	expectedMac := getMacString(algo, key, message)
	if expectedMac == signatureParts[1] {
		return true
	}

	return false
}
