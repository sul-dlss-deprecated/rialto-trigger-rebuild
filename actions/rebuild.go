package actions

import (
	"github.com/knakk/sparql"
	"github.com/sul-dlss/rialto-trigger-rebuild/runtime"
)

// Action is an abstract type of work this can do.
type Action interface {
	Run() error
}

// RebuildAction drops the repository and then rebuilds it
type RebuildAction struct {
	registry *runtime.Registry
}

// NewRebuildAction creates a Rebuild ("rebuild") action
func NewRebuildAction(registry *runtime.Registry) Action {
	return &RebuildAction{registry: registry}
}

// Run does the rebuilding
func (r *RebuildAction) Run() error {
	err := r.registry.Derivative.RemoveAll()
	if err != nil {
		return err
	}
	var resourceTypes = []string{"http://xmlns.com/foaf/0.1/Organization",
		"http://xmlns.com/foaf/0.1/Person",
		"http://vivoweb.org/ontology/core#Grant",
		"http://www.w3.org/2004/02/skos/core#Concept",
		"http://purl.org/ontology/bibo/Document",
		"http://xmlns.com/foaf/0.1/Project"}
	for _, resourceType := range resourceTypes {
		err = r.runType(resourceType)
		if err != nil {
			break
		}
	}
	return err
}

func (r *RebuildAction) runType(resourceType string) error {
	err := r.registry.Canonical.QueryResources(resourceType, func(results *sparql.Results) error {
		innerErr := r.registry.Topic.Publish(r.solutionsToEntities(results))

		return innerErr
	})
	return err
}

func (r *RebuildAction) solutionsToEntities(results *sparql.Results) []string {
	entities := []string{}
	for _, solution := range results.Solutions() {
		entities = append(entities, solution["s"].String())
	}
	return entities
}
