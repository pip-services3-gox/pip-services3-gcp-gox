package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	cconf "github.com/pip-services3-gox/pip-services3-commons-gox/config"
	"github.com/pip-services3-gox/pip-services3-commons-gox/convert"
	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	cerr "github.com/pip-services3-gox/pip-services3-commons-gox/errors"
	crefer "github.com/pip-services3-gox/pip-services3-commons-gox/refer"
	ccount "github.com/pip-services3-gox/pip-services3-components-gox/count"
	clog "github.com/pip-services3-gox/pip-services3-components-gox/log"
	ctrace "github.com/pip-services3-gox/pip-services3-components-gox/trace"
	gcpconn "github.com/pip-services3-gox/pip-services3-gcp-gox/connect"
	rpcsrv "github.com/pip-services3-gox/pip-services3-rpc-gox/services"
)

// Abstract client that calls Google Functions.
//
// When making calls "cmd" parameter determines which what action shall be called, while
// other parameters are passed to the action itself.
//
//	Configuration parameters
//		- connections:
//			- uri:           full connection uri with specific app and function name
//			- protocol:      connection protocol
//			- project_id:    is your Google Cloud Platform project ID
//			- region:        is the region where your function is deployed
//			- function:      is the name of the HTTP function you deployed
//			- org_id:        organization name
//		- options:
//			- retries:               number of retries (default: 3)
//			- connect_timeout:       connection timeout in milliseconds (default: 10 sec)
//			- timeout:               invocation timeout in milliseconds (default: 10 sec)
//		- credentials:
//			- account: the service account name
//			- auth_token:    Google-generated ID token or null if using custom auth (IAM)
//
//	References
//		- *:logger:*:*:1.0				(optional) ILogger components to pass log messages
//		- *:counters:*:*:1.0			(optional) ICounters components to pass collected measurements
//		- *:discovery:*:*:1.0			(optional) IDiscovery services to resolve connection
//		- *:credential-store:*:*:1.0	(optional) Credential stores to resolve credentials
//
// see CloudFunction, CommandableGoogleClient
//
//	Example:
//		type MyCloudFunctionClient struct {
//			clients.CloudFunctionClient
//		}
//
//		func NewMyCloudFunctionClient() *MyCloudFunctionClient {
//			return &MyCloudFunctionClient{
//				CloudFunctionClient: *gcpclient.NewCloudFunctionClient(),
//			}
//		}
//
//		func (c *MyCloudFunctionClient) GetData(ctx context.Context, correlationId string, id string) MyData {
//			timing := c.Instrument(ctx, correlationId, "myclient.get_data")
//
//			response, err := c.Call(ctx, "get_data", correlationId, data.NewAnyValueMapFromTuples("id", dummyId))
//
//			defer timing.EndTiming(ctx, err)
//			return rpcclients.HandleHttpResponse[MyData](response, correlationId)
//		}
//
//		...
//
//		client := NewMyCloudFunctionClient()
//		client.Configure(config.NewConfigParamsFromTuples(
//			"connection.uri", "http://region-id.cloudfunctions.net/myfunction",
//			"connection.protocol", "http",
//			"connection.region", "region",
//			"connection.function", "myfunction",
//			"connection.project_id", "id",
//			"credential.auth_token", "XXX",
//		))
//		result := client.GetData("123", "1")
//
type CloudFunctionClient struct {
	//The HTTP client.
	Client *http.Client
	// The Google Function connection parameters
	Connection *gcpconn.GcpConnectionParams
	//The number of retries.
	Retries int
	// The default headers to be added to every request.
	Headers *cdata.StringValueMap
	//The connection timeout in milliseconds.
	ConnectTimeout int
	//The invocation timeout in milliseconds.
	Timeout int
	//The remote service uri which is calculated on open.
	Uri string
	//The connection resolver.
	ConnectionResolver *gcpconn.GcpConnectionResolver
	//The dependency resolver.
	DependencyResolver *crefer.DependencyResolver

	//The logger.
	Logger *clog.CompositeLogger
	//The performance counters.
	Counters *ccount.CompositeCounters
	// The tracer.
	Tracer *ctrace.CompositeTracer

	defaultConfig *cconf.ConfigParams
}

const (
	DefaultConnectTimeout = 10000
	DefaultTimeout        = 10000
	DefaultRetriesCount   = 3
)

// Creates new instance of CloudFunctionClient
func NewCloudFunctionClient() *CloudFunctionClient {
	c := CloudFunctionClient{}

	c.defaultConfig = cconf.NewConfigParamsFromTuples(
		"options.connectTimeout", DefaultConnectTimeout,
		"options.timeout", DefaultTimeout,
		"options.retries", DefaultRetriesCount,
	)

	c.ConnectionResolver = gcpconn.NewGcpConnectionResolver()
	c.DependencyResolver = crefer.NewDependencyResolver()
	c.Logger = clog.NewCompositeLogger()
	c.Counters = ccount.NewCompositeCounters()
	c.Tracer = ctrace.NewCompositeTracer(context.Background(), nil)
	c.Headers = cdata.NewEmptyStringValueMap()

	return &c
}

// Configure object by passing configuration parameters.
//	Parameters:
//		- ctx context.Context
//		- config: ConfigParams configuration parameters to be set.
func (c *CloudFunctionClient) Configure(ctx context.Context, config *cconf.ConfigParams) {
	c.ConnectionResolver.Configure(ctx, config)
	c.DependencyResolver.Configure(ctx, config)

	c.Retries = config.GetAsIntegerWithDefault("options.retries", DefaultRetriesCount)
	c.ConnectTimeout = config.GetAsIntegerWithDefault("options.connectTimeout", DefaultConnectTimeout)
	c.Timeout = config.GetAsIntegerWithDefault("options.timeout", DefaultTimeout)
}

// SetReferences sets references to dependent components.
//	see IReferences
//	Parameters:
//		- ctx context.Context
//		- references IReferences references to locate the component dependencies.
func (c *CloudFunctionClient) SetReferences(ctx context.Context, references crefer.IReferences) {
	c.Logger.SetReferences(ctx, references)
	c.Counters.SetReferences(ctx, references)
	c.ConnectionResolver.SetReferences(ctx, references)
	c.DependencyResolver.SetReferences(ctx, references)
}

// Instrument method are adds instrumentation to log calls and measure call time.
// It returns a services.InstrumentTiming object that is used to end the time measurement.
//	Parameters:
//		- ctx context.Context
//		- correlationId string (optional) transaction id to trace execution through call chain.
//		- name string a method name.
//	Returns: services.InstrumentTiming object to end the time measurement.
func (c *CloudFunctionClient) Instrument(ctx context.Context, correlationId string, name string) *rpcsrv.InstrumentTiming {
	c.Logger.Trace(ctx, correlationId, "Calling %s method", name)
	c.Counters.IncrementOne(ctx, name+".call_count")
	counterTiming := c.Counters.BeginTiming(ctx, name+".call_time")
	traceTiming := c.Tracer.BeginTrace(ctx, correlationId, name, "")
	return rpcsrv.NewInstrumentTiming(correlationId, name, "call",
		c.Logger, c.Counters, counterTiming, traceTiming)
}

// IsOpen Checks if the component is opened.
//	Returns: bool true if the component has been opened and false otherwise.
func (c *CloudFunctionClient) IsOpen() bool {
	return c.Client != nil
}

// Open opens the component.
//	Parameters:
//		- ctx context.Context
//		- correlationId: string transaction id to trace execution through call chain.
//	Return: error
func (c *CloudFunctionClient) Open(ctx context.Context, correlationId string) error {
	if c.IsOpen() {
		return nil
	}

	connection, err := c.ConnectionResolver.Resolve(correlationId)
	if err != nil {
		return err
	}

	c.Uri, _ = connection.Uri()
	c.Client = &http.Client{
		// Timeout includes connection time, any redirects, and reading the response body
		Timeout: time.Duration(c.Timeout+c.ConnectTimeout) * time.Millisecond,
	}

	if c.Client == nil {
		return cerr.NewConnectionError(
			correlationId,
			"CANNOT_CONNECT",
			"Connection to Google function service failed",
		).WithDetails("url", c.Uri)
	}

	c.Logger.Debug(ctx, correlationId, "Google function client connected to %s", c.Uri)

	return nil
}

// Closes component and frees used resources.
// Parameters:
//		-correlationId	(optional) transaction id to trace execution through call chain.
func (c *CloudFunctionClient) Close(ctx context.Context, correlationId string) error {
	if c.Client != nil {
		c.Logger.Debug(ctx, correlationId, "Closed Google function service at %s", c.Uri)
		c.Client = nil
		c.Uri = ""
	}
	return nil
}

// Performs Google Function invocation.
// Parameters:
//		- cmd	an action name to be called.
//		- correlationId	(optional) transaction id to trace execution through call chain.
//		- args	action arguments
// Returns action result.
func (c *CloudFunctionClient) Call(ctx context.Context, cmd string, correlationId string,
	args *cdata.AnyValueMap) (*http.Response, error) {
	if cmd == "" {
		cerr.NewUnknownError(correlationId, "NO_COMMAND", "Missing command")
	}

	if !c.IsOpen() {
		return nil, nil
	}

	if correlationId == "" {
		correlationId = cdata.IdGenerator.NextShort()
	}

	if args == nil {
		args = cdata.NewEmptyAnyValueMap()
	}
	args.Put("cmd", cmd)
	args.Put("correlation_id", correlationId)

	var jsonStr string
	if args != nil {
		jsonStr, _ = convert.JsonConverter.ToJson(args.Value())
	}

	retries := c.Retries
	var response *http.Response

	for retries > 0 {
		req, err := c.prepareRequest(ctx, correlationId, http.MethodPost, c.Uri, []byte(jsonStr))
		if err != nil {
			return nil, err
		}

		response, err = c.Client.Do(req)
		if err != nil {
			retries--
			if retries == 0 {
				return nil, cerr.NewUnknownError(
					correlationId,
					"COMMUNICATION_ERROR",
					"Unknown communication problem on GCP client",
				).
					WithCause(err)
			}

			err = c.waitForRetry(ctx, correlationId, retries)
			if err != nil {
				return nil, err
			}
			continue
		}
		break
	}

	if response.StatusCode == 204 {
		_ = response.Body.Close()
		return nil, nil
	}

	if response.StatusCode >= 400 {
		defer response.Body.Close()
		return nil, c.handleResponseError(response, correlationId)
	}

	return response, nil
}

// AddFilterParams method are adds filter parameters (with the same name as they defined)
// to invocation parameter map.
//	Parameters:
//		- params  *cdata.StringValueMap      invocation parameters.
//		- filter  *cdata.FilterParams     (optional) filter parameters
//	Returns: invocation parameters with added filter parameters.
func (c *CloudFunctionClient) AddFilterParams(params *cdata.StringValueMap, filter *cdata.FilterParams) *cdata.StringValueMap {

	if params == nil {
		params = cdata.NewEmptyStringValueMap()
	}
	if filter != nil {
		for k, v := range filter.Value() {
			params.Put(k, v)
		}
	}
	return params
}

// AddPagingParams method are adds paging parameters (skip, take, total) to invocation parameter map.
// Parameters:
//    - params        invocation parameters.
//    - paging        (optional) paging parameters
// Return invocation parameters with added paging parameters.
func (c *CloudFunctionClient) AddPagingParams(params *cdata.StringValueMap, paging *cdata.PagingParams) *cdata.StringValueMap {
	if params == nil {
		params = cdata.NewEmptyStringValueMap()
	}

	if paging != nil {
		params.Put("total", paging.Total)
		if paging.Skip >= 0 {
			params.Put("skip", paging.Skip)
		}
		if paging.Take >= 0 {
			params.Put("take", paging.Take)
		}
	}

	return params
}

func (c *CloudFunctionClient) waitForRetry(ctx context.Context, correlationId string, retries int) error {
	waitTime := c.Timeout * int(math.Pow(float64(c.Retries-retries), 2))

	select {
	case <-time.After(time.Duration(waitTime) * time.Millisecond):
		return nil
	case <-ctx.Done():
		return cerr.ApplicationErrorFactory.Create(
			&cerr.ErrorDescription{
				Type:          "Application",
				Category:      "Application",
				Code:          "CONTEXT_CANCELLED",
				Message:       "request canceled by parent context",
				CorrelationId: correlationId,
			},
		)
	}
}

func (c *CloudFunctionClient) prepareRequest(ctx context.Context, correlationId string,
	method string, url string, body []byte) (*http.Request, error) {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, cerr.NewUnknownError(
			correlationId,
			"UNSUPPORTED_METHOD",
			"Method is not supported by GCP client",
		).
			WithDetails("verb", method).
			WithCause(err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	for k, v := range c.Headers.Value() {
		req.Header.Set(k, v)
	}

	return req, nil
}

func (c *CloudFunctionClient) handleResponseError(response *http.Response, correlationId string) error {
	r, rErr := ioutil.ReadAll(response.Body)
	if rErr != nil {
		eDesct := cerr.ErrorDescription{
			Type:          "Application",
			Category:      "Application",
			Status:        response.StatusCode,
			Code:          "",
			Message:       rErr.Error(),
			CorrelationId: correlationId,
		}
		return cerr.ApplicationErrorFactory.Create(&eDesct).WithCause(rErr)
	}

	appErr := cerr.ApplicationError{}
	_ = json.Unmarshal(r, &appErr)
	if appErr.Status == 0 && len(r) > 0 { // not standart Pip.Services error
		values := make(map[string]any)
		decodeErr := json.Unmarshal(r, &values)
		if decodeErr != nil { // not json response
			appErr.Message = (string)(r)
		}
		appErr.Details = values
	}
	appErr.Status = response.StatusCode
	return &appErr
}
