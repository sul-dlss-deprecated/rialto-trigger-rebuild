package repository

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/knakk/sparql"
)

// Reader is the interface to the sparql reader
type Reader interface {
	QueryResources(f func(*sparql.Results) error) error
}

// SparqlReader is the functions that interact with the SPARQL repository
type SparqlReader struct {
	Repo *sparql.Repo
}

// QueryResources calls the function with a list of resources (in managable sized chunks) populated
// by querying for everything in the triplestore
func (r *SparqlReader) QueryResources(resourceType string, f func(*sparql.Results) error) error {
	return r.queryPage(
		func(offset int) string {
			return fmt.Sprintf(`SELECT ?s
	WHERE {
	  ?s a <%v> .
	}
		ORDER BY ?s OFFSET %v LIMIT %v`, resourceType, offset, tripleLimit())
		}, f)
}

// Calls sparqlForOffset once for each page to generate the query
// Calls f() on each page of results
func (r *SparqlReader) queryPage(sparqlForOffset func(offset int) string, f func(*sparql.Results) error) error {
	page := 0
	for {
		offset := page * tripleLimit()
		query := sparqlForOffset(offset)
		log.Printf("[SPARQL] %s", query)
		results, err := r.Repo.Query(query)
		log.Printf("Returned from query")
		if err != nil {
			log.Printf("Returning error: %s", err)
			return err
		}
		log.Printf("Counting results")
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

// Number of records to pull per requst.  If this is too large, then the SPARQL proxy
// lambda will hit a 6mb limit. See: https://github.com/sul-dlss/sparql-loader/issues/44
func tripleLimit() int {
	i, _ := strconv.Atoi(os.Getenv("SPARQL_LIMIT"))
	return i
}
