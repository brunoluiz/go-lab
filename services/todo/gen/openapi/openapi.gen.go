// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package openapi

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
	strictgin "github.com/oapi-codegen/runtime/strictmiddleware/gin"
)

// AddListJSONBody defines parameters for AddList.
type AddListJSONBody struct {
	Title string `json:"title"`
}

// UpdateListJSONBody defines parameters for UpdateList.
type UpdateListJSONBody struct {
	Title string `json:"title"`
}

// AddListJSONRequestBody defines body for AddList for application/json ContentType.
type AddListJSONRequestBody AddListJSONBody

// UpdateListJSONRequestBody defines body for UpdateList for application/json ContentType.
type UpdateListJSONRequestBody UpdateListJSONBody

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// AddListWithBody request with any body
	AddListWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	AddList(ctx context.Context, body AddListJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteList request
	DeleteList(ctx context.Context, listId string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetListById request
	GetListById(ctx context.Context, listId string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateListWithBody request with any body
	UpdateListWithBody(ctx context.Context, listId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateList(ctx context.Context, listId string, body UpdateListJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) AddListWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAddListRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AddList(ctx context.Context, body AddListJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAddListRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteList(ctx context.Context, listId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteListRequest(c.Server, listId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetListById(ctx context.Context, listId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetListByIdRequest(c.Server, listId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateListWithBody(ctx context.Context, listId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateListRequestWithBody(c.Server, listId, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateList(ctx context.Context, listId string, body UpdateListJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateListRequest(c.Server, listId, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewAddListRequest calls the generic AddList builder with application/json body
func NewAddListRequest(server string, body AddListJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewAddListRequestWithBody(server, "application/json", bodyReader)
}

// NewAddListRequestWithBody generates requests for AddList with any type of body
func NewAddListRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v1/lists")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewDeleteListRequest generates requests for DeleteList
func NewDeleteListRequest(server string, listId string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "list_id", runtime.ParamLocationPath, listId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v1/lists/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetListByIdRequest generates requests for GetListById
func NewGetListByIdRequest(server string, listId string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "list_id", runtime.ParamLocationPath, listId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v1/lists/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewUpdateListRequest calls the generic UpdateList builder with application/json body
func NewUpdateListRequest(server string, listId string, body UpdateListJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpdateListRequestWithBody(server, listId, "application/json", bodyReader)
}

// NewUpdateListRequestWithBody generates requests for UpdateList with any type of body
func NewUpdateListRequestWithBody(server string, listId string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "list_id", runtime.ParamLocationPath, listId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v1/lists/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// AddListWithBodyWithResponse request with any body
	AddListWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*AddListResponse, error)

	AddListWithResponse(ctx context.Context, body AddListJSONRequestBody, reqEditors ...RequestEditorFn) (*AddListResponse, error)

	// DeleteListWithResponse request
	DeleteListWithResponse(ctx context.Context, listId string, reqEditors ...RequestEditorFn) (*DeleteListResponse, error)

	// GetListByIdWithResponse request
	GetListByIdWithResponse(ctx context.Context, listId string, reqEditors ...RequestEditorFn) (*GetListByIdResponse, error)

	// UpdateListWithBodyWithResponse request with any body
	UpdateListWithBodyWithResponse(ctx context.Context, listId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateListResponse, error)

	UpdateListWithResponse(ctx context.Context, listId string, body UpdateListJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateListResponse, error)
}

type AddListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *struct {
		CreatedAt time.Time `json:"created_at"`
		Tasks     *[]string `json:"tasks,omitempty"`
		Title     string    `json:"title"`
		UniqId    string    `json:"uniq_id"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	JSON400 *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
}

// Status returns HTTPResponse.Status
func (r AddListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r AddListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Success bool `json:"success"`
	}
	JSON404 *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
}

// Status returns HTTPResponse.Status
func (r DeleteListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetListByIdResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		CreatedAt time.Time `json:"created_at"`
		Tasks     *[]string `json:"tasks,omitempty"`
		Title     string    `json:"title"`
		UniqId    string    `json:"uniq_id"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	JSON400 *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	JSON404 *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
}

// Status returns HTTPResponse.Status
func (r GetListByIdResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetListByIdResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateListResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		CreatedAt time.Time `json:"created_at"`
		Tasks     *[]string `json:"tasks,omitempty"`
		Title     string    `json:"title"`
		UniqId    string    `json:"uniq_id"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	JSON400 *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	JSON404 *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
}

// Status returns HTTPResponse.Status
func (r UpdateListResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateListResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// AddListWithBodyWithResponse request with arbitrary body returning *AddListResponse
func (c *ClientWithResponses) AddListWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*AddListResponse, error) {
	rsp, err := c.AddListWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAddListResponse(rsp)
}

func (c *ClientWithResponses) AddListWithResponse(ctx context.Context, body AddListJSONRequestBody, reqEditors ...RequestEditorFn) (*AddListResponse, error) {
	rsp, err := c.AddList(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAddListResponse(rsp)
}

// DeleteListWithResponse request returning *DeleteListResponse
func (c *ClientWithResponses) DeleteListWithResponse(ctx context.Context, listId string, reqEditors ...RequestEditorFn) (*DeleteListResponse, error) {
	rsp, err := c.DeleteList(ctx, listId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteListResponse(rsp)
}

// GetListByIdWithResponse request returning *GetListByIdResponse
func (c *ClientWithResponses) GetListByIdWithResponse(ctx context.Context, listId string, reqEditors ...RequestEditorFn) (*GetListByIdResponse, error) {
	rsp, err := c.GetListById(ctx, listId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetListByIdResponse(rsp)
}

// UpdateListWithBodyWithResponse request with arbitrary body returning *UpdateListResponse
func (c *ClientWithResponses) UpdateListWithBodyWithResponse(ctx context.Context, listId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateListResponse, error) {
	rsp, err := c.UpdateListWithBody(ctx, listId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateListResponse(rsp)
}

func (c *ClientWithResponses) UpdateListWithResponse(ctx context.Context, listId string, body UpdateListJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateListResponse, error) {
	rsp, err := c.UpdateList(ctx, listId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateListResponse(rsp)
}

// ParseAddListResponse parses an HTTP response from a AddListWithResponse call
func ParseAddListResponse(rsp *http.Response) (*AddListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &AddListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest struct {
			CreatedAt time.Time `json:"created_at"`
			Tasks     *[]string `json:"tasks,omitempty"`
			Title     string    `json:"title"`
			UniqId    string    `json:"uniq_id"`
			UpdatedAt time.Time `json:"updated_at"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	}

	return response, nil
}

// ParseDeleteListResponse parses an HTTP response from a DeleteListWithResponse call
func ParseDeleteListResponse(rsp *http.Response) (*DeleteListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Success bool `json:"success"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

// ParseGetListByIdResponse parses an HTTP response from a GetListByIdWithResponse call
func ParseGetListByIdResponse(rsp *http.Response) (*GetListByIdResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetListByIdResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			CreatedAt time.Time `json:"created_at"`
			Tasks     *[]string `json:"tasks,omitempty"`
			Title     string    `json:"title"`
			UniqId    string    `json:"uniq_id"`
			UpdatedAt time.Time `json:"updated_at"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

// ParseUpdateListResponse parses an HTTP response from a UpdateListWithResponse call
func ParseUpdateListResponse(rsp *http.Response) (*UpdateListResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpdateListResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			CreatedAt time.Time `json:"created_at"`
			Tasks     *[]string `json:"tasks,omitempty"`
			Title     string    `json:"title"`
			UniqId    string    `json:"uniq_id"`
			UpdatedAt time.Time `json:"updated_at"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Add a new list to the store
	// (POST /api/v1/lists)
	AddList(c *gin.Context)
	// Deletes a list
	// (DELETE /api/v1/lists/{list_id})
	DeleteList(c *gin.Context, listId string)
	// Find list by ID
	// (GET /api/v1/lists/{list_id})
	GetListById(c *gin.Context, listId string)
	// Update existing list
	// (PUT /api/v1/lists/{list_id})
	UpdateList(c *gin.Context, listId string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// AddList operation middleware
func (siw *ServerInterfaceWrapper) AddList(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.AddList(c)
}

// DeleteList operation middleware
func (siw *ServerInterfaceWrapper) DeleteList(c *gin.Context) {

	var err error

	// ------------- Path parameter "list_id" -------------
	var listId string

	err = runtime.BindStyledParameterWithOptions("simple", "list_id", c.Param("list_id"), &listId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter list_id: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.DeleteList(c, listId)
}

// GetListById operation middleware
func (siw *ServerInterfaceWrapper) GetListById(c *gin.Context) {

	var err error

	// ------------- Path parameter "list_id" -------------
	var listId string

	err = runtime.BindStyledParameterWithOptions("simple", "list_id", c.Param("list_id"), &listId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter list_id: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetListById(c, listId)
}

// UpdateList operation middleware
func (siw *ServerInterfaceWrapper) UpdateList(c *gin.Context) {

	var err error

	// ------------- Path parameter "list_id" -------------
	var listId string

	err = runtime.BindStyledParameterWithOptions("simple", "list_id", c.Param("list_id"), &listId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter list_id: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.UpdateList(c, listId)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.POST(options.BaseURL+"/api/v1/lists", wrapper.AddList)
	router.DELETE(options.BaseURL+"/api/v1/lists/:list_id", wrapper.DeleteList)
	router.GET(options.BaseURL+"/api/v1/lists/:list_id", wrapper.GetListById)
	router.PUT(options.BaseURL+"/api/v1/lists/:list_id", wrapper.UpdateList)
}

type AddListRequestObject struct {
	Body *AddListJSONRequestBody
}

type AddListResponseObject interface {
	VisitAddListResponse(w http.ResponseWriter) error
}

type AddList201JSONResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Tasks     *[]string `json:"tasks,omitempty"`
	Title     string    `json:"title"`
	UniqId    string    `json:"uniq_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (response AddList201JSONResponse) VisitAddListResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	return json.NewEncoder(w).Encode(response)
}

type AddList400JSONResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (response AddList400JSONResponse) VisitAddListResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type DeleteListRequestObject struct {
	ListId string `json:"list_id"`
}

type DeleteListResponseObject interface {
	VisitDeleteListResponse(w http.ResponseWriter) error
}

type DeleteList200JSONResponse struct {
	Success bool `json:"success"`
}

func (response DeleteList200JSONResponse) VisitDeleteListResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type DeleteList404JSONResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (response DeleteList404JSONResponse) VisitDeleteListResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type GetListByIdRequestObject struct {
	ListId string `json:"list_id"`
}

type GetListByIdResponseObject interface {
	VisitGetListByIdResponse(w http.ResponseWriter) error
}

type GetListById200JSONResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Tasks     *[]string `json:"tasks,omitempty"`
	Title     string    `json:"title"`
	UniqId    string    `json:"uniq_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (response GetListById200JSONResponse) VisitGetListByIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetListById400JSONResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (response GetListById400JSONResponse) VisitGetListByIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type GetListById404JSONResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (response GetListById404JSONResponse) VisitGetListByIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type UpdateListRequestObject struct {
	ListId string `json:"list_id"`
	Body   *UpdateListJSONRequestBody
}

type UpdateListResponseObject interface {
	VisitUpdateListResponse(w http.ResponseWriter) error
}

type UpdateList200JSONResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Tasks     *[]string `json:"tasks,omitempty"`
	Title     string    `json:"title"`
	UniqId    string    `json:"uniq_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (response UpdateList200JSONResponse) VisitUpdateListResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type UpdateList400JSONResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (response UpdateList400JSONResponse) VisitUpdateListResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type UpdateList404JSONResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (response UpdateList404JSONResponse) VisitUpdateListResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Add a new list to the store
	// (POST /api/v1/lists)
	AddList(ctx context.Context, request AddListRequestObject) (AddListResponseObject, error)
	// Deletes a list
	// (DELETE /api/v1/lists/{list_id})
	DeleteList(ctx context.Context, request DeleteListRequestObject) (DeleteListResponseObject, error)
	// Find list by ID
	// (GET /api/v1/lists/{list_id})
	GetListById(ctx context.Context, request GetListByIdRequestObject) (GetListByIdResponseObject, error)
	// Update existing list
	// (PUT /api/v1/lists/{list_id})
	UpdateList(ctx context.Context, request UpdateListRequestObject) (UpdateListResponseObject, error)
}

type StrictHandlerFunc = strictgin.StrictGinHandlerFunc
type StrictMiddlewareFunc = strictgin.StrictGinMiddlewareFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// AddList operation middleware
func (sh *strictHandler) AddList(ctx *gin.Context) {
	var request AddListRequestObject

	var body AddListJSONRequestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Status(http.StatusBadRequest)
		ctx.Error(err)
		return
	}
	request.Body = &body

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.AddList(ctx, request.(AddListRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "AddList")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(AddListResponseObject); ok {
		if err := validResponse.VisitAddListResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}

// DeleteList operation middleware
func (sh *strictHandler) DeleteList(ctx *gin.Context, listId string) {
	var request DeleteListRequestObject

	request.ListId = listId

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.DeleteList(ctx, request.(DeleteListRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "DeleteList")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(DeleteListResponseObject); ok {
		if err := validResponse.VisitDeleteListResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetListById operation middleware
func (sh *strictHandler) GetListById(ctx *gin.Context, listId string) {
	var request GetListByIdRequestObject

	request.ListId = listId

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetListById(ctx, request.(GetListByIdRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetListById")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(GetListByIdResponseObject); ok {
		if err := validResponse.VisitGetListByIdResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}

// UpdateList operation middleware
func (sh *strictHandler) UpdateList(ctx *gin.Context, listId string) {
	var request UpdateListRequestObject

	request.ListId = listId

	var body UpdateListJSONRequestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Status(http.StatusBadRequest)
		ctx.Error(err)
		return
	}
	request.Body = &body

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.UpdateList(ctx, request.(UpdateListRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "UpdateList")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(UpdateListResponseObject); ok {
		if err := validResponse.VisitUpdateListResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xX34vbRhD+V5Z5akGxfUkeip7acCSYtBykF2gox7HWju1xpd3Nzugc1+h/L7uSrrYs",
	"iMOF0hA/WV7Nj29mvvnW3kPhKu8sWmHI900GZJcO8nhuRRcSH7HSVELeH/28CLV1ZU1/TywKNBkY5CKQ",
	"F3IWcrh1z4xTJbGoSlu9wgqtQAZCUmL/HjJ4wMCtx9VkNpnFQM6j1Z4ghxeT2eQFZOC1rCMymGpP04er",
	"aYybDrzjhO44+S/GKK0sblsA4pSsUbG4gJDiBx0t56a1/ZU4Qgv4sUaWV87s+toj5nwP2vuSiuQz3XBM",
	"sQcu1ljphCHEiEKYEHUF7gE/6cqnWn/bKXlsR+zBzsdjlkB2BU3TpqaABvI/uwB3j2ZuscFCoDm2k1Bj",
	"OmDvLLepn8+ungC8CKgFzb1OnksXqvgERgs+E6rwFHgGovmv5EyCVVv+iUl7oEPQu/T9i/uTQW3p4z2Z",
	"Y69ZvWHe3q7nb27ff6iev/3p9ebmZvmHfLh6OxrDmy8sbzCXHkRfQnbYsaP447M7pujvdVEgc0T2cjZ7",
	"yticwdHGV8isV2PvBoX1hlkb6xz0t2tU3bqoJAhkWZF90CUZRdbXwikN11Wlw+6zKyl6xRFKIsBd9Dza",
	"9Ok+ftyTadpdL1HwdOuv07nSPY2O97x9262610FXKBhi1mGcaKPIRIxdqiiIkCcZggysrmJ3Okgw3Mrs",
	"YFTDxt+dbOxTRs8dh85ZpoVzJWp7Mvw+xjlTv+kbqpZkiddo1JZkrboY6gfrVFfKjy2xX36LxGaPBS0J",
	"jQrIrg4Fqq1mZZ2opautGTC7JRb/y7sBmTNY4cgd9Q6lDja6MdlVieOsfYMS6fhqNzefo+38Wrnl43aF",
	"FP7/ytzLXXO5a77qXfO9iM1rsqbd8cVOza/H1MbXI2rzPjFG4SdiErKrcbVprc65I9/1YBOGr6kx3/xv",
	"8IsuXnTxoov/sS4e6NuhvA3/VzTNPwEAAP//7naf7OwQAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
