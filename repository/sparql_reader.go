package repository

import (
	"fmt"
	"log"

	"github.com/knakk/sparql"
)

// Number of records to pull per requst.  If this is too large, then the SPARQL proxy
// lambda will hit a 6mb limit. See: https://github.com/sul-dlss-labs/sparql-loader/issues/44
// Also consider that this is the size of the Entities field in the message.  Too many entities
// Can cause the SNS message to be too large
const tripleLimit = 1000

type Reader interface {
	QueryResources(f func(*sparql.Results) error) error
}

// SparqlReader is the functions that interact with the SPARQL repository
type SparqlReader struct {
	Repo *sparql.Repo
}

// QueryResources calls the function with a list of resources (in managable sized chunks) populated
// by querying for everything in the triplestore
func (r *SparqlReader) QueryResources(f func(*sparql.Results) error) error {
	return r.queryPage(
		func(offset int) string {
			return fmt.Sprintf(`SELECT DISTINCT ?s
	WHERE {
	  ?s a ?type .
	}
		ORDER BY ?s OFFSET %v LIMIT %v`, offset, tripleLimit)
		}, f)
}

// Calls sparqlForOffset once for each page to generate the query
// Calls f() on each page of results
func (r *SparqlReader) queryPage(sparqlForOffset func(offset int) string, f func(*sparql.Results) error) error {
	page := 0
	for {
		offset := page * tripleLimit
		query := sparqlForOffset(offset)
		log.Printf("[SPARQL] %s", query)
		results, err := r.Repo.Query(query)
		if err != nil {
			return err
		}
		if resultCount := len(results.Solutions()); resultCount == 0 {
			break
		}
		if err = f(results); err != nil {
			return err
		}
		page++
	}
	return nil
}
