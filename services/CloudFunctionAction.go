package services

import (
	"net/http"

	cvalid "github.com/pip-services3-gox/pip-services3-commons-gox/validate"
)

type CloudFunctionAction struct {
	// Command to call the action
	Cmd string
	// Schema to validate action parameters
	Schema *cvalid.Schema
	// Action to be executed
	Action http.HandlerFunc
}
