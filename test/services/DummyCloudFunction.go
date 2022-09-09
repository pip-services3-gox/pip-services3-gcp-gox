package services_test

import (
	gcpsrv "github.com/pip-services3-gox/pip-services3-gcp-gox/containers"
	tbuild "github.com/pip-services3-gox/pip-services3-gcp-gox/test/build"
)

type DummyCloudFunction struct {
	*gcpsrv.CloudFunction
}

func NewDummyCloudFunction() *DummyCloudFunction {
	c := DummyCloudFunction{CloudFunction: gcpsrv.NewCloudFunctionWithParams("dummy", "Dummy cloud function")}
	c.AddFactory(tbuild.NewDummyFactory())
	c.AddFactory(NewDummyCloudFunctionServiceFactory())

	return &c
}
