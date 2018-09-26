package runtime

import (
	"github.com/sul-dlss-labs/rialto-derivatives/derivative"
	"github.com/sul-dlss-labs/rialto-trigger-rebuild/messages"
	"github.com/sul-dlss-labs/rialto-trigger-rebuild/repository"
)

// Registry is the object that holds all the external services
type Registry struct {
	Canonical  repository.Reader
	Derivative derivative.Writer
	Topic      messages.MessageService
}

// NewRegistry creates a new instance of the service registry
func NewRegistry(repo repository.Reader, writer derivative.Writer, conn messages.MessageService) *Registry {
	return &Registry{
		Canonical:  repo,
		Derivative: writer,
		Topic:      conn,
	}
}
