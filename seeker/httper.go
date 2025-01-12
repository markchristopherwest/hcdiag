package seeker

import (
	"github.com/hashicorp/hcdiag/client"
)

// HTTPer hits APIs.
type HTTPer struct {
	Path   string            `json:"path"`
	Client *client.APIClient `json:"client"`
}

func NewHTTPer(client *client.APIClient, path string) *Seeker {
	return &Seeker{
		Identifier: "GET" + " " + path,
		Runner: HTTPer{
			Client: client,
			Path:   path,
		},
	}
}

// Run executes a GET request to the Path using the Client
func (h HTTPer) Run() (interface{}, Status, error) {
	result, err := h.Client.Get(h.Path)
	if err != nil {
		return result, Unknown, err
	}
	return result, Success, nil
}
