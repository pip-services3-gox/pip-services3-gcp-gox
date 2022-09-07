package services

import (
	"context"
	"net/http"

	ccomand "github.com/pip-services3-gox/pip-services3-commons-gox/commands"
	crun "github.com/pip-services3-gox/pip-services3-commons-gox/run"
	gcputil "github.com/pip-services3-gox/pip-services3-gcp-gox/utils"
	rpcserv "github.com/pip-services3-gox/pip-services3-rpc-gox/services"
)

// Abstract service that receives commands via Google Function protocol
// to operations automatically generated for commands defined in ccomand.ICommandable components.
// Each command is exposed as invoke method that receives command name and parameters.
//
// Commandable services require only 3 lines of code to implement a robust external
// Google Function-based remote interface.
//
// This service is intended to work inside Google Function container that
// exploses registered actions externally.
//
// 	Configuration parameters:
//		- dependencies:
//			- controller:            override for Controller dependency
// 	References
//		- *:logger:*:*:1.0			(optional) ILogger components to pass log messages
//		- *:counters:*:*:1.0		(optional) ICounters components to pass collected measurements
//
// see CloudFunctionService
//
// 	Example:
//		type MyCommandableCloudFunctionService struct {
//			gcpsrv.CommandableCloudFunctionService
//		}
//
//		func NewMyCommandableCloudFunctionService() *MyCommandableCloudFunctionService {
//			c := MyCommandableCloudFunctionService{}
//			c.CommandableCloudFunctionService = *gcpsrv.NewCommandableCloudFunctionService("mydata")
//			c.DependencyResolver.Put(context.Background(), "controller", crefer.NewDescriptor("mygroup", "controller", "default", "*", "*"))
//			return &c
//		}
//
///		...
//
//		service := NewMyCommandableCloudFunctionService()
//		service.SetReferences(crefer.NewReferencesFromTuples(
//			crefer.NewDescriptor("mygroup","controller","default","default","1.0"), controller,
//		))
//		service.Open(ctx, "123")
//		fmt.Println("The Google Function service is running")
//
type CommandableCloudFunctionService struct {
	CloudFunctionService
	commandSet *ccomand.CommandSet
}

// Creates a new instance of the service.
// Parameters:
// 		- name 	a service name.
func NewCommandableCloudFunctionService(name string) *CommandableCloudFunctionService {
	c := CommandableCloudFunctionService{}
	c.CloudFunctionService = *NewCloudFunctionService(name)

	return &c
}

// Returns body from Google Function request.
// This method can be overloaded in child classes
// Parameters:
//		- req	Google Function request
// Returns Parameters from request
func (c *CommandableCloudFunctionService) GetParameters(req *http.Request) *crun.Parameters {
	return gcputil.CloudFunctionRequestHelper.GetParameters(req)
}

// Registers all actions in Google Function.
func (c *CommandableCloudFunctionService) Register() {
	resCtrl, depErr := c.DependencyResolver.GetOneRequired("controller")
	if depErr != nil {
		panic(depErr)
	}

	controller, ok := resCtrl.(ccomand.ICommandable)
	if !ok {
		c.Logger.Error(context.Background(), "CommandableHttpService", nil, "Can't cast Controller to ICommandable")
		return
	}

	c.commandSet = controller.GetCommandSet()
	commands := c.commandSet.Commands()

	for index := 0; index < len(commands); index++ {
		command := commands[index]
		name := command.Name()

		c.RegisterAction(c.name, nil, func(w http.ResponseWriter, r *http.Request) {
			correlationId := c.GetCorrelationId(r)
			args := c.GetParameters(r)
			args.Remove("correlation_id")

			timing := c.Instrument(r.Context(), correlationId, name)
			execRes, execErr := command.Execute(r.Context(), correlationId, args)
			timing.EndTiming(r.Context(), execErr)
			rpcserv.HttpResponseSender.SendResult(w, r, execRes, execErr)
		})
	}
}
