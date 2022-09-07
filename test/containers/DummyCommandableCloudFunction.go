package containers_test

import (
	"context"

	crefer "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	gcpcont "github.com/pip-services3-gox/pip-services3-gcp-gox/containers"
	tbuild "github.com/pip-services3-gox/pip-services3-gcp-gox/test/build"
)

type DummyCommandableCloudFunction struct {
	gcpcont.CommandableCloudFunction
}

func NewDummyCommandableCloudFunction() *DummyCommandableCloudFunction {
	c := DummyCommandableCloudFunction{}
	c.CommandableCloudFunction = *gcpcont.NewCommandableCloudFunctionWithParams("dummy", "Dummy commandable cloud function")
	c.DependencyResolver.Put(context.Background(), "controller", crefer.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))

	c.AddFactory(tbuild.NewDummyFactory())

	return &c
}
