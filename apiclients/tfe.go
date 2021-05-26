package apiclients

// https://www.terraform.io/docs/cloud/api/index.html

import (
	"os"
)

const DefaultTFEAddr = "https://127.0.0.1"

// NewTFEAPI returns an APIClient for TFE.
func NewTFEAPI() *APIClient {
	addr := os.Getenv("TFE_HTTP_ADDR")
	if addr == "" {
		addr = DefaultTFEAddr
	}

	headers := map[string]string{}
	token := os.Getenv("TFE_TOKEN")
	if token != "" {
		headers["Authorization"] = ("Bearer " + token)
	}

	return NewAPIClient("tfe", addr, headers)
}