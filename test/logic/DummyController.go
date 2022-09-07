package test_logic

import (
	"context"

	ccomand "github.com/pip-services3-gox/pip-services3-commons-gox/commands"
	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	tdata "github.com/pip-services3-gox/pip-services3-gcp-gox/test/data"
)

type DummyController struct {
	commandSet *DummyCommandSet
	entities   []tdata.Dummy
}

func NewDummyController() *DummyController {
	return &DummyController{
		entities: make([]tdata.Dummy, 0),
	}
}

// GetCommandSet gets a command set with all supported commands and events.
//	see CommandSet
//	Returns: *CommandSet a command set with commands and events.
func (c *DummyController) GetCommandSet() *ccomand.CommandSet {
	if c.commandSet == nil {
		c.commandSet = NewDummyCommandSet(c)
	}

	return &c.commandSet.CommandSet
}

func (c *DummyController) GetPageByFilter(ctx context.Context, correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (*cdata.DataPage[tdata.Dummy], error) {

	if filter == nil {
		filter = cdata.NewEmptyFilterParams()
	}
	var key string = filter.GetAsString("key")

	if paging == nil {
		paging = cdata.NewEmptyPagingParams()
	}
	var skip int64 = paging.GetSkip(0)
	var take int64 = paging.GetTake(100)

	var result []tdata.Dummy

	for i := 0; i < len(c.entities); i++ {
		var entity tdata.Dummy = c.entities[i]
		if key != "" && key != entity.Key {
			continue
		}

		skip--
		if skip >= 0 {
			continue
		}

		take--
		if take < 0 {
			break
		}

		result = append(result, entity)
	}
	var total int64 = (int64)(len(result))
	return cdata.NewDataPage(result, int(total)), nil
}

func (c *DummyController) GetOneById(ctx context.Context, correlationId string, id string) (result tdata.Dummy, err error) {
	for _, entity := range c.entities {
		if entity.Id == id {
			return entity, nil
		}
	}

	return tdata.Dummy{}, nil
}

func (c *DummyController) Create(ctx context.Context, correlationId string, entity tdata.Dummy) (result tdata.Dummy, err error) {
	if entity.Id == "" {
		entity.Id = cdata.IdGenerator.NextLong()
		c.entities = append(c.entities, entity)
	}
	return entity, nil
}

func (c *DummyController) Update(ctx context.Context, correlationId string, newEntity tdata.Dummy) (result tdata.Dummy, err error) {
	for index := 0; index < len(c.entities); index++ {
		var entity tdata.Dummy = c.entities[index]
		if entity.Id == newEntity.Id {
			c.entities[index] = newEntity
			return newEntity, nil

		}
	}
	return tdata.Dummy{}, nil
}

func (c *DummyController) DeleteById(ctx context.Context, correlationId string, id string) (result tdata.Dummy, err error) {
	var entity tdata.Dummy

	for i := 0; i < len(c.entities); {
		entity = c.entities[i]
		if entity.Id == id {
			if i == len(c.entities)-1 {
				c.entities = c.entities[:i]
			} else {
				c.entities = append(c.entities[:i], c.entities[i+1:]...)
			}
			return entity, nil
		} else {
			i++
		}
	}
	return tdata.Dummy{}, nil
}
