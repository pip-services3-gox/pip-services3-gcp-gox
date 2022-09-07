package test_logic

import (
	cconvert "github.com/pip-services3-gox/pip-services3-commons-gox/convert"
	cvalid "github.com/pip-services3-gox/pip-services3-commons-gox/validate"
)

type DummySchema struct {
	cvalid.ObjectSchema
}

func NewDummySchema() *DummySchema {
	c := DummySchema{
		ObjectSchema: *cvalid.NewObjectSchema(),
	}

	c.WithOptionalProperty("id", cconvert.String)
	c.WithRequiredProperty("key", cconvert.String)
	c.WithOptionalProperty("content", cconvert.String)

	return &c
}
