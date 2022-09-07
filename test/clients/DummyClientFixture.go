package clients_test

import (
	"context"
	"testing"

	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	tdata "github.com/pip-services3-gox/pip-services3-gcp-gox/test/data"
	"github.com/stretchr/testify/assert"
)

type DummyClientFixture struct {
	client IDummyClient
}

func NewDummyClientFixture(client IDummyClient) *DummyClientFixture {
	dcf := DummyClientFixture{client: client}
	return &dcf
}

func (c *DummyClientFixture) TestCrudOperations(t *testing.T) {
	ctx := context.Background()
	dummy1 := tdata.Dummy{Id: "", Key: "Key 1", Content: "Content 1"}
	dummy2 := tdata.Dummy{Id: "", Key: "Key 2", Content: "Content 2"}

	// Create one dummy
	dummy, err := c.client.CreateDummy(ctx, "ClientFixture", dummy1)
	assert.Nil(t, err)
	assert.Equal(t, dummy.Content, dummy1.Content)
	assert.Equal(t, dummy.Key, dummy1.Key)
	dummy1 = dummy

	// Create another dummy
	dummy, err = c.client.CreateDummy(ctx, "ClientFixture", dummy2)
	assert.Nil(t, err)
	assert.Equal(t, dummy.Content, dummy2.Content)
	assert.Equal(t, dummy.Key, dummy2.Key)
	dummy2 = dummy

	// Get all dummies
	dummies, err := c.client.GetDummies(ctx, "ClientFixture", *cdata.NewEmptyFilterParams(), *cdata.NewPagingParams(0, 5, false))
	assert.Nil(t, err)
	assert.NotNil(t, dummies)
	assert.Len(t, dummies.Data, 2)

	// Update the dummy
	dummy1.Content = "Updated Content 1"
	dummy, err = c.client.UpdateDummy(ctx, "ClientFixture", dummy1)
	assert.Nil(t, err)
	assert.Equal(t, dummy.Content, "Updated Content 1")
	assert.Equal(t, dummy.Key, dummy1.Key)
	dummy1 = dummy

	// Delete dummy
	dummy, err = c.client.DeleteDummy(ctx, "ClientFixture", dummy1.Id)
	assert.Nil(t, err)

	// Try to get delete dummy
	dummy, err = c.client.GetDummyById(ctx, "ClientFixture", dummy1.Id)
	assert.Nil(t, err)
	assert.Equal(t, tdata.Dummy{}, dummy)
}
