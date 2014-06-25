// Package handler contains the handlers for the webhook requests.
package handler

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"github.com/brocaar/dockerbuilder/worker"
	"hash"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type GitHubHandler struct {
	taskQueue worker.TaskQueue
	secret    []byte
}

// NewGitHubHandler returns a new instance of GitHubHandler.
func NewGitHubHandler(taskQueue worker.TaskQueue, secret []byte) *GitHubHandler {
	return &GitHubHandler{
		taskQueue: taskQueue,
		secret:    secret,
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
	if len(h.secret) > 0 && checkGitHubMac(h.secret, r.Header.Get("X-Hub-Signature"), requestBody) != true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// github ping
	if r.Header.Get("X-Github-Event") == "ping" {
		w.WriteHeader(http.StatusOK)
		return
	}
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
