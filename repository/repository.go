// Package repository contains tools related to code repositories.
package repository

// Constants related to the CSM type.
const (
	ScmGit = "git"
)

// Constants related to the CSM host.
const (
	HostGitHub = "github.com"
)

// Repository holds all the repository related data.
type Repository struct {
	Owner string
	Name  string
	Host  string
	SCM   string
	URL   string
}

// NewRepository returns a new instance of the Repository type.
func NewRepository(host, name, scm string) *Repository {
	return &Repository{
		Name:  name,
		Owner: name,
		Host:  host,
		SCM:   scm,
	}
}
