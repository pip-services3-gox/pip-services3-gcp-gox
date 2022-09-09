package test_logic

import (
	"context"
	"encoding/json"

	ccomand "github.com/pip-services3-gox/pip-services3-commons-gox/commands"
	cconv "github.com/pip-services3-gox/pip-services3-commons-gox/convert"
	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	crun "github.com/pip-services3-gox/pip-services3-commons-gox/run"
	cvalid "github.com/pip-services3-gox/pip-services3-commons-gox/validate"
	tdata "github.com/pip-services3-gox/pip-services3-gcp-gox/test/data"
)

type DummyCommandSet struct {
	ccomand.CommandSet
	controller IDummyController
}

func NewDummyCommandSet(controller IDummyController) *DummyCommandSet {
	c := DummyCommandSet{
		CommandSet: *ccomand.NewCommandSet(),
		controller: controller,
	}

	c.AddCommand(c.makeGetPageByFilterCommand())
	c.AddCommand(c.makeGetOneByIdCommand())
	c.AddCommand(c.makeCreateCommand())
	c.AddCommand(c.makeUpdateCommand())
	c.AddCommand(c.makeDeleteByIdCommand())

	return &c
}

func (c *DummyCommandSet) makeGetPageByFilterCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"get_dummies",
		cvalid.NewObjectSchema().WithOptionalProperty("filter", cvalid.NewFilterParamsSchema()).WithOptionalProperty("paging", cvalid.NewPagingParamsSchema()),
		func(ctx context.Context, correlationId string, args *crun.Parameters) (result any, err error) {
			var filter *cdata.FilterParams
			var paging *cdata.PagingParams

			if _val, ok := args.Get("filter"); ok {
				filter = cdata.NewFilterParamsFromValue(_val)
			}
			if _val, ok := args.Get("paging"); ok {
				paging = cdata.NewPagingParamsFromValue(_val)
			}

			return c.controller.GetPageByFilter(ctx, correlationId, filter, paging)
		},
	)
}

func (c *DummyCommandSet) makeGetOneByIdCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"get_dummy_by_id",
		cvalid.NewObjectSchema().WithRequiredProperty("dummy_id", cconv.String),
		func(ctx context.Context, correlationId string, args *crun.Parameters) (any, error) {
			id := args.GetAsString("dummy_id")
			return c.controller.GetOneById(ctx, correlationId, id)
		},
	)
}

func (c *DummyCommandSet) makeCreateCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"create_dummy",
		cvalid.NewObjectSchema().WithRequiredProperty("dummy", tdata.NewDummySchema()),
		func(ctx context.Context, correlationId string, args *crun.Parameters) (any, error) {
			var entity tdata.Dummy

			if _val, ok := args.Get("dummy"); ok {
				jsonStr, _ := cconv.JsonConverter.ToJson(_val)
				err := json.Unmarshal([]byte(jsonStr), &entity)
				if err != nil {
					return nil, err
				}
			}

			return c.controller.Create(ctx, correlationId, entity)
		},
	)
}

func (c *DummyCommandSet) makeUpdateCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"update_dummy",
		cvalid.NewObjectSchema().WithRequiredProperty("dummy", tdata.NewDummySchema()),
		func(ctx context.Context, correlationId string, args *crun.Parameters) (any, error) {
			var entity tdata.Dummy

			if _val, ok := args.Get("dummy"); ok {
				jsonStr, _ := cconv.JsonConverter.ToJson(_val)
				err := json.Unmarshal([]byte(jsonStr), &entity)
				if err != nil {
					return nil, err
				}
			}

			return c.controller.Update(ctx, correlationId, entity)
		},
	)
}

func (c *DummyCommandSet) makeDeleteByIdCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"delete_dummy",
		cvalid.NewObjectSchema().WithRequiredProperty("dummy_id", cconv.String),
		func(ctx context.Context, correlationId string, args *crun.Parameters) (any, error) {
			id := args.GetAsString("dummy_id")
			return c.controller.DeleteById(ctx, correlationId, id)
		},
	)
}
