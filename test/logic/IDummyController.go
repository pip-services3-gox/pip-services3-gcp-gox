package test_logic

import (
	"context"

	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	tdata "github.com/pip-services3-gox/pip-services3-gcp-gox/test/data"
)

type IDummyController interface {
	GetPageByFilter(ctx context.Context, correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *cdata.DataPage[tdata.Dummy], err error)
	GetOneById(ctx context.Context, correlationId string, id string) (result tdata.Dummy, err error)
	Create(ctx context.Context, correlationId string, entity tdata.Dummy) (result tdata.Dummy, err error)
	Update(ctx context.Context, correlationId string, entity tdata.Dummy) (result tdata.Dummy, err error)
	DeleteById(ctx context.Context, correlationId string, id string) (result tdata.Dummy, err error)
}
