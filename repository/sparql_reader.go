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
	QueryResources(resoureType string, f func(*sparql.Results) error) error
}

// SparqlReader is the functions that interact with the SPARQL repository
type SparqlReader struct {
	Repo *sparql.Repo
}

// QueryResources calls the function with a list of resources (in managable sized chunks) populated
// by querying for everything in the triplestore
func (r *SparqlReader) QueryResources(resourceType string, f func(*sparql.Results) error) error {
	log.Printf("[RESOURCE TYPE] %s", resourceType)
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
		results, err := r.Repo.Query(query)
		if err != nil {
			return err
		}
		resultCount := len(results.Solutions())
		log.Printf("[SPARQL RESULTS] %v results for page %v", resultCount, offset)
		if resultCount == 0 {
			break
		}
		err = f(results)
		if err != nil {
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
