package internal

import (
	"errors"
	"net/http"
)

var ErrTestFailed = errors.New("test failed expected parameters")

type RequestMethod string
type ObjectPath string

const (
	MethodGet     RequestMethod = http.MethodGet
	MethodHead    RequestMethod = http.MethodHead
	MethodPost    RequestMethod = http.MethodPost
	MethodPut     RequestMethod = http.MethodPut
	MethodPatch   RequestMethod = http.MethodPatch
	MethodDelete  RequestMethod = http.MethodDelete
	MethodConnect RequestMethod = http.MethodConnect
	MethodOptions RequestMethod = http.MethodOptions
	MethodTrace   RequestMethod = http.MethodTrace

	APIPrefixPathAndVersioningForV1 = "/api/v1"

	InfraWorkspace   ObjectPath = "/infrastructure/workspaces"
	ServiceWorkspace ObjectPath = "/services/workspaces"
	Environment      ObjectPath = "/infrastructure/workspaces/{workspaceName}/environments"
	Proxy            ObjectPath = "/services/workspaces/{workspaceName}/proxies"
	ApiDocs          ObjectPath = "/services/workspaces/{workspaceName}/api-docs"
	Devportals       ObjectPath = "/infrastructure/workspaces/{workspaceName}/devportals"
)

type InfraWorkspaceStruct struct {
	Name string
}
