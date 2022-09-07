package clients_test

import (
	"context"
	"os"
	"testing"

	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
)

func TestDummyCloudFunctionClient(t *testing.T) {

	functionName := os.Getenv("GCP_FUNCTION_NAME")
	protocol := os.Getenv("GCP_FUNCTION_PROTOCOL")
	region := os.Getenv("GCP_FUNCTION_REGION")
	projectId := os.Getenv("GCP_PROJECT_ID")
	uri := os.Getenv("GCP_FUNCTION_URI")

	// if uri == "" {
	// 	uri = "http://localhost:3000" // uncomment for local test
	// }

	if uri == "" && (region == "" || functionName == "" || protocol == "" || projectId == "") {
		t.Skip("No credentials set, skip TestDummyCloudFunctionClient")
	}

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.uri", uri,
		"connection.protocol", protocol,
		"connection.region", region,
		"connection.function", functionName,
		"connection.project_id", projectId,
	)

	client := NewDummyCloudFunctionClient()
	fixture := NewDummyClientFixture(client)

	client.Configure(context.Background(), restConfig)
	client.SetReferences(context.Background(), cref.NewEmptyReferences())
	client.Open(context.Background(), "")
	defer client.Close(context.Background(), "DummyCloudFunctionClient")

	t.Run("DummyCloudFunctionClient.CrudOperations", fixture.TestCrudOperations)
}
