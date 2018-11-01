package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/knakk/sparql"
	"github.com/sul-dlss/rialto-derivatives/derivative"

	"github.com/sul-dlss/rialto-trigger-rebuild/messages"
	"github.com/sul-dlss/rialto-trigger-rebuild/repository"

	"github.com/sul-dlss/rialto-trigger-rebuild/actions"
	"github.com/sul-dlss/rialto-trigger-rebuild/runtime"

	// Added for the postgres driver
	_ "github.com/lib/pq"
)

func main() {
	// repo := derivativeRepo.BuildRepository()
	endpoint := os.Getenv("SPARQL_ENDPOINT")
	neptune := buildSparqlRepo(endpoint)
	writer := repository.NewCombinedWriter(buildSolrClient(), buildDatabase())
	snsConn := buildSNSConn()
	topicArn := os.Getenv("RIALTO_TOPIC_ARN")
	messageService := messages.NewSNSMessageService(snsConn, &topicArn)
	registry := runtime.NewRegistry(neptune, writer, messageService)

	err := actions.NewRebuildAction(registry).Run()
	if err != nil {
		log.Printf("Error rebuilding: %v", err)
		panic(err)
	}
}

func buildSparqlRepo(url string) repository.Reader {
	repo, err := sparql.NewRepo(url,
		sparql.Timeout(time.Second*60),
	)
	if err != nil {
		panic(err)
	}
	return &repository.SparqlReader{Repo: repo}
}

func buildSNSConn() *sns.SNS {
	endpoint := os.Getenv("RIALTO_SNS_ENDPOINT")
	return sns.New(session.New(), aws.NewConfig().
		WithDisableSSL(false).
		WithEndpoint(endpoint))
}

func buildSolrClient() *derivative.SolrClient {
	// indexer := transform.NewCompositeIndexer(repo)
	host := os.Getenv("SOLR_HOST")
	collection := os.Getenv("SOLR_COLLECTION")
	return derivative.NewSolrClient(host, collection, nil)
}

func buildDatabase() *derivative.PostgresClient {
	conf := derivative.NewPostgresConfig().
		WithUser(os.Getenv("RDS_USERNAME")).
		WithPassword(os.Getenv("RDS_PASSWORD")).
		WithDbname(os.Getenv("RDS_DB_NAME")).
		WithHost(os.Getenv("RDS_HOSTNAME")).
		WithPort(os.Getenv("RDS_PORT"))

	db, err := sql.Open("postgres", conf.ToConnString())
	if err != nil {
		panic(err)
	}

	return &derivative.PostgresClient{
		DB: db,
	}
}
