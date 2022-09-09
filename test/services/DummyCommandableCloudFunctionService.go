package services_test

import (
	"context"

	crefer "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	gcpserv "github.com/pip-services3-gox/pip-services3-gcp-gox/services"
)

type DummyCommandableCloudFunctionService struct {
	*gcpserv.CommandableCloudFunctionService
}

func NewDummyCommandableCloudFunctionService() *DummyCommandableCloudFunctionService {
	c := DummyCommandableCloudFunctionService{}
	c.CommandableCloudFunctionService = gcpserv.NewCommandableCloudFunctionService("dummies")
	c.DependencyResolver.Put(context.Background(), "controller", crefer.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	return &c
}
