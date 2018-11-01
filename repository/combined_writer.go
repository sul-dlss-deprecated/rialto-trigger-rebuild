package repository

import (
	"errors"

	"github.com/sul-dlss/rialto-derivatives/derivative"
	"github.com/sul-dlss/rialto-derivatives/models"
)

// CombinedWriter writes both to PostgreSQL and Solr and fulfills the derivative.Writer interface
type CombinedWriter struct {
	solr *derivative.SolrClient
	db   *derivative.PostgresClient
}

// NewCombinedWriter creates a new CombinedWriter instance
func NewCombinedWriter(solr *derivative.SolrClient, db *derivative.PostgresClient) derivative.Writer {
	return &CombinedWriter{
		solr: solr,
		db:   db,
	}
}

// Add is not implemented, but exists only to fulfil the interface
func (d *CombinedWriter) Add([]models.Resource) error {
	return errors.New("Add is not implemented in rialto-trigger-rebuild")
}

// RemoveAll removes all records from Solr and all rows from PostgreSQL
func (d *CombinedWriter) RemoveAll() error {
	err := d.solr.RemoveAll()
	if err != nil {
		return err
	}

	return d.db.RemoveAll()
}
