package containers

import (
	"context"
	"net/http"

	ccomand "github.com/pip-services3-gox/pip-services3-commons-gox/commands"
	crun "github.com/pip-services3-gox/pip-services3-commons-gox/run"
	gcputil "github.com/pip-services3-gox/pip-services3-gcp-gox/utils"
	rpcserv "github.com/pip-services3-gox/pip-services3-rpc-gox/services"
)

// Abstract Google Function function, that acts as a container to instantiate and run components
// and expose them via external entry point. All actions are automatically generated for commands
// defined in ICommandable components. Each command is exposed as an action defined by "cmd" parameter.
//
// Container configuration for this Google Function is stored in "./config/config.yml" file.
// But this path can be overridden by <code>CONFIG_PATH</code> environment variable.
//
//	References
//		- *:logger:*:*:1.0							(optional) ILogger components to pass log messages
//		- *:counters:*:*:1.0						(optional) ICounters components to pass collected measurements
//		- *:service:gcp-function:*:1.0       		(optional) ICloudFunctionService services to handle action requests
//		- *:service:commandable-gcp-function:*:1.0	(optional) ICloudFunctionService services to handle action requests
//
//	Example:
//		type MyCloudFunction struct {
//			*containers.CommandableCloudFunction
//			controller IMyController
//		}
//
//		func NewMyCloudFunction() *MyCloudFunction {
//			c := MyCloudFunction{}
//			c.CloudFunction = containers.NewCommandableCloudFunctionWithParams("mygroup", "MyGroup CloudFunction")
//
//			return &c
//		}
//
//		...
//
//		cloudFunction := NewMyCloudFunction()
//		cloudFunction.Run(ctx)
//		fmt.Println("MyCloudFunction is started")
//
// Deprecated: This component has been deprecated. Use CloudFunctionService instead.
type CommandableCloudFunction struct {
	*CloudFunction
}

// Creates a new instance of this Google Function.
func NewCommandableCloudFunction() *CommandableCloudFunction {
	c := CommandableCloudFunction{}
	c.CloudFunction = InheritCloudFunction(&c)
	return &c
}

// Creates a new instance of this Google Function.
// Parameters:
//		- name	(optional) a container name (accessible via ContextInfo)
//		- description	(optional) a container description (accessible via ContextInfo)
func NewCommandableCloudFunctionWithParams(name string, description string) *CommandableCloudFunction {
	c := CommandableCloudFunction{}
	c.CloudFunction = InheritCloudFunctionWithParams(&c, name, description)
	return &c
}

// Returns body from Google Function request.
// This method can be overloaded in child classes
// Parameters:
//		- req	Googl Function request
// Returns Parameters from request
func (c *CommandableCloudFunction) GetParameters(req *http.Request) *crun.Parameters {
	return gcputil.CloudFunctionRequestHelper.GetParameters(req)
}

func (c *CommandableCloudFunction) registerCommandSet(commandSet *ccomand.CommandSet) {
	commands := commandSet.Commands()
	for index := 0; index < len(commands); index++ {
		command := commands[index]

		c.RegisterAction(command.Name(), nil, func(w http.ResponseWriter, r *http.Request) {
			correlationId := c.GetCorrelationId(r)
			args := c.GetParameters(r)

			timing := c.Instrument(r.Context(), correlationId, command.Name())
			execRes, execErr := command.Execute(r.Context(), correlationId, args)
			timing.EndTiming(r.Context(), execErr)

			rpcserv.HttpResponseSender.SendResult(w, r, execRes, execErr)
		})
	}
}

// Registers all actions in this Google Function.
//
// Deprecated: Overloading of this method has been deprecated. Use CloudFunctionService instead.
func (c *CommandableCloudFunction) Register() {
	resCtrl, depErr := c.DependencyResolver.GetOneRequired("controller")
	if depErr != nil {
		panic(depErr)
	}

	controller, ok := resCtrl.(ccomand.ICommandable)
	if !ok {
		c.Logger().Error(context.Background(), "CommandableHttpService", nil, "Can't cast Controller to ICommandable")
		return
	}

	commandSet := controller.GetCommandSet()
	c.registerCommandSet(commandSet)
}
