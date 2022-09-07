package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	crun "github.com/pip-services3-gox/pip-services3-commons-gox/run"
)

var CloudFunctionRequestHelper = _TCloudFunctionRequestHelper{}

type _TCloudFunctionRequestHelper struct {
}

func (c *_TCloudFunctionRequestHelper) GetCorrelationId(req *http.Request) string {
	correlationId := req.URL.Query().Get("correlation_id")
	if correlationId == "" {
		correlationId = req.Header.Get("correlation_id")
	}
	return correlationId
}

func (c *_TCloudFunctionRequestHelper) GetCommand(req *http.Request) (string, error) {
	cmd := req.URL.Query().Get("cmd")

	if cmd == "" {
		var body map[string]any

		err := c.DecodeBody(req, &body)

		if err != nil {
			return "", err
		}

		if val, ok := body["cmd"].(string); ok {
			cmd = val
		}
	}

	return cmd, nil
}

func (c *_TCloudFunctionRequestHelper) DecodeBody(req *http.Request, target any) error {
	bodyBytes, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, target)

	if err != nil {
		return err
	}

	_ = req.Body.Close()
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return nil
}

func (c *_TCloudFunctionRequestHelper) GetParameters(req *http.Request) *crun.Parameters {
	var params map[string]any

	_ = c.DecodeBody(req, &params) // Ignore the error

	return crun.NewParametersFromValue(params)
}
