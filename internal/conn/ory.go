package conn

import (
	ory "github.com/ory/client-go"
)

var Ory *ory.APIClient

func NewOryClient() *ory.APIClient {
	config := ory.NewConfiguration()
	return ory.NewAPIClient(config)
}

func InitOry(client *ory.APIClient) {
	Ory = client
}
