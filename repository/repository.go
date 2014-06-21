// Package repository contains tools related to code repositories.
package repository

import (
	"fmt"
)

// Constants related to the CSM type.
const (
	ScmGit = "git"
)

// Constants related to the CSM host.
const (
	HostGitHub = "github.com"
)

const (
	repoPatternGitHub = "git@github.com:%s/%s.git"
)

// Repository holds all the repository related data.
type Repository struct {
	Owner string
	Name  string
	Host  string
	SCM   string
	Url   string
}

// NewRepository returns a new instance of the Repository type.
func NewRepository(host, owner, name, scm string) *Repository {
	var url string

	if host == HostGitHub {
		url = fmt.Sprintf(repoPatternGitHub, owner, name)
	}

	return &Repository{
		Host:  host,
		Owner: owner,
		Name:  name,
		SCM:   scm,
		Url:   url,
	}
}
