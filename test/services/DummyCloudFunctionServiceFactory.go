package services_test

import (
	cref "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	cbuild "github.com/pip-services3-gox/pip-services3-components-gox/build"
)

type DummyCloudFunctionServiceFactory struct {
	cbuild.Factory
	Descriptor                *cref.Descriptor
	ControllerDescriptor      *cref.Descriptor
	CloudServiceDescriptor    *cref.Descriptor
	CmdCloudServiceDescriptor *cref.Descriptor
}

func NewDummyCloudFunctionServiceFactory() *DummyCloudFunctionServiceFactory {

	c := DummyCloudFunctionServiceFactory{
		Factory:                   *cbuild.NewFactory(),
		Descriptor:                cref.NewDescriptor("pip-services-dummies", "factory", "default", "default", "1.0"),
		CloudServiceDescriptor:    cref.NewDescriptor("pip-services-dummies", "service", "cloudfunc", "*", "1.0"),
		CmdCloudServiceDescriptor: cref.NewDescriptor("pip-services-dummies", "service", "commandable-cloudfunc", "*", "1.0"),
	}

	c.RegisterType(c.CloudServiceDescriptor, NewDummyCloudFunctionService)
	c.RegisterType(c.CmdCloudServiceDescriptor, NewDummyCommandableCloudFunctionService)
	return &c
}
