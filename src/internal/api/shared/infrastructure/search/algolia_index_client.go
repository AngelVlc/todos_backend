package search

import (
	"fmt"
	"log"

	sharedApp "github.com/AngelVlc/todos_backend/src/internal/api/shared/application"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

type AlgoliaIndexClient struct {
	client    *search.Client
	indexName string
	index     *search.Index
	cfgSvr    sharedApp.ConfigurationService
}

func NewAlgoliaIndexClient(cfgSvr sharedApp.ConfigurationService, objectName string) *AlgoliaIndexClient {
	client := search.NewClient(cfgSvr.GetAlgoliaAppId(), cfgSvr.GetAlgoliaApiKey())
	indexName := fmt.Sprintf("%v-%v", cfgSvr.GetEnvironment(), objectName)
	index := client.InitIndex(indexName)

	return &AlgoliaIndexClient{
		client:    client,
		index:     index,
		cfgSvr:    cfgSvr,
		indexName: indexName,
	}
}

func (c *AlgoliaIndexClient) SaveObjects(objects interface{}) error {
	res, err := c.index.SaveObjects(objects)

	for _, r := range res.Responses {
		log.Printf("Indexed documents in %v with IDs %v in task %v", c.indexName, r.ObjectIDs, r.TaskID)
	}

	return err
}

func (c *AlgoliaIndexClient) DeleteObject(objectID string) error {
	res, err := c.index.DeleteObject(objectID)

	log.Printf("Deleted indexed document in %v with ID %v in task %v", c.indexName, objectID, res.TaskID)

	return err
}
