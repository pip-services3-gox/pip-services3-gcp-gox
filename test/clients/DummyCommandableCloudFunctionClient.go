package clients_test

import (
	"context"

	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	gcpclient "github.com/pip-services3-gox/pip-services3-gcp-gox/clients"
	tdata "github.com/pip-services3-gox/pip-services3-gcp-gox/test/data"
	rpcclient "github.com/pip-services3-gox/pip-services3-rpc-gox/clients"
)

type DummyCommandableCloudFunctionClient struct {
	*gcpclient.CommandableCloudFunctionClient
}

func NewDummyCommandableCloudFunctionClient() *DummyCommandableCloudFunctionClient {
	return &DummyCommandableCloudFunctionClient{
		CommandableCloudFunctionClient: gcpclient.NewCommandableCloudFunctionClient("dummies"),
	}
}

func (c *DummyCommandableCloudFunctionClient) GetDummies(ctx context.Context, correlationId string, filter cdata.FilterParams, paging cdata.PagingParams) (result cdata.DataPage[tdata.Dummy], err error) {
	params := cdata.NewEmptyStringValueMap()
	c.AddFilterParams(params, &filter)
	c.AddPagingParams(params, &paging)

	response, err := c.CallCommand(ctx, "dummies.get_dummies", correlationId, cdata.NewAnyValueMapFromValue(params.Value()))
	if err != nil {
		return *cdata.NewEmptyDataPage[tdata.Dummy](), err
	}

	return rpcclient.HandleHttpResponse[cdata.DataPage[tdata.Dummy]](response, correlationId)
}

func (c *DummyCommandableCloudFunctionClient) GetDummyById(ctx context.Context, correlationId string, dummyId string) (result tdata.Dummy, err error) {
	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy_id", dummyId)

	response, err := c.CallCommand(ctx, "dummies.get_dummy_by_id", correlationId, params)
	if err != nil {
		return tdata.Dummy{}, err
	}

	return rpcclient.HandleHttpResponse[tdata.Dummy](response, correlationId)
}

func (c *DummyCommandableCloudFunctionClient) CreateDummy(ctx context.Context, correlationId string, dummy tdata.Dummy) (result tdata.Dummy, err error) {
	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy", dummy)

	response, err := c.CallCommand(ctx, "dummies.create_dummy", correlationId, params)
	if err != nil {
		return tdata.Dummy{}, err
	}

	return rpcclient.HandleHttpResponse[tdata.Dummy](response, correlationId)
}

func (c *DummyCommandableCloudFunctionClient) UpdateDummy(ctx context.Context, correlationId string, dummy tdata.Dummy) (result tdata.Dummy, err error) {
	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy", dummy)

	response, err := c.CallCommand(ctx, "dummies.update_dummy", correlationId, params)
	if err != nil {
		return tdata.Dummy{}, err
	}

	return rpcclient.HandleHttpResponse[tdata.Dummy](response, correlationId)
}

func (c *DummyCommandableCloudFunctionClient) DeleteDummy(ctx context.Context, correlationId string, dummyId string) (result tdata.Dummy, err error) {
	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy_id", dummyId)

	response, err := c.CallCommand(ctx, "dummies.delete_dummy", correlationId, params)
	if err != nil {
		return tdata.Dummy{}, err
	}

	return rpcclient.HandleHttpResponse[tdata.Dummy](response, correlationId)
}
