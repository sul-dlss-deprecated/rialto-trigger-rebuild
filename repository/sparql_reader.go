package repository

import (
	"fmt"
	"log"

	"github.com/knakk/sparql"
)

// Reader queries a list of resources
type Reader interface {
	QueryResources(f func(*sparql.Results) error) error
}

// SparqlReader is the functions that interact with the SPARQL repository
type SparqlReader struct {
	Repo        *sparql.Repo
	TripleLimit int
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
		ORDER BY ?s OFFSET %v LIMIT %v`, offset, r.TripleLimit)
		}, f)
}

// Calls sparqlForOffset once for each page to generate the query
// Calls f() on each page of results
func (r *SparqlReader) queryPage(sparqlForOffset func(offset int) string, f func(*sparql.Results) error) error {
	page := 0
	for {
		offset := page * r.TripleLimit
		query := sparqlForOffset(offset)
		log.Printf("[SPARQL] %s", query)
		results, err := r.Repo.Query(query)
		if err != nil {
			return err
		}
		resultCount := len(results.Solutions())
		log.Printf("[RESULTS] %v", resultCount)
		if resultCount == 0 {
			break
		}
		if err = f(results); err != nil {
			return err
		}
		page++
	}
	return nil
}
