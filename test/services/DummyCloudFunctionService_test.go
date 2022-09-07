package services_test

import (
	"context"
	"testing"

	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	"github.com/stretchr/testify/assert"
)

type DummyCloudFunctionServiceTest struct {
	fixture       *DummyCloudFunctionFixture
	funcContainer *DummyCloudFunction
}

func newDummyCloudFunctionServiceTest() *DummyCloudFunctionServiceTest {
	return &DummyCloudFunctionServiceTest{}
}

func (c *DummyCloudFunctionServiceTest) setup(t *testing.T) {
	config := cconf.NewConfigParamsFromTuples(
		"logger.descriptor", "pip-services:logger:console:default:1.0",
		"service.descriptor", "pip-services-dummies:service:gcp-function:default:1.0",
	)

	ctx := context.Background()

	c.funcContainer = NewDummyCloudFunction()
	c.funcContainer.Configure(ctx, config)
	err := c.funcContainer.Open(ctx, "")
	assert.Nil(t, err)

	c.fixture = NewDummyCloudFunctionFixture(c.funcContainer.GetHandler())
}

func (c *DummyCloudFunctionServiceTest) teardown(t *testing.T) {
	err := c.funcContainer.Close(context.Background(), "")
	assert.Nil(t, err)
}

func TestCrudOperationsCloudService(t *testing.T) {
	c := newDummyCloudFunctionServiceTest()
	if c == nil {
		return
	}

	c.setup(t)
	t.Run("CRUD Operations", c.fixture.TestCrudOperations)
	c.teardown(t)
}
