package build_test

import (
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	cbuild "github.com/pip-services3-gox/pip-services3-components-gox/build"
	tlogic "github.com/pip-services3-gox/pip-services3-gcp-gox/test/logic"
)

type DummyFactory struct {
	cbuild.Factory
	Descriptor           *cref.Descriptor
	ControllerDescriptor *cref.Descriptor
}

// NewDefaultRpcFactory creates a new instance of the factory.
func NewDummyFactory() *DummyFactory {
	c := DummyFactory{
		Factory:              *cbuild.NewFactory(),
		Descriptor:           cref.NewDescriptor("pip-services-dummies", "factory", "default", "default", "1.0"),
		ControllerDescriptor: cref.NewDescriptor("pip-services-dummies", "controller", "default", "*", "1.0"),
	}

	c.RegisterType(c.ControllerDescriptor, tlogic.NewDummyController)
	return &c
}
