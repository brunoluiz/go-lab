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

// Radar defines model for Radar.
type Radar struct {
	CreatedAt time.Time       `json:"created_at"`
	Id        string          `json:"id"`
	Items     []RadarItem     `json:"items"`
	Quadrants []RadarQuadrant `json:"quadrants"`
	Title     string          `json:"title"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// RadarItem defines model for RadarItem.
type RadarItem struct {
	CreatedAt   time.Time     `json:"created_at"`
	Description string        `json:"description"`
	Id          string        `json:"id"`
	Name        string        `json:"name"`
	Quadrant    RadarQuadrant `json:"quadrant"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// RadarQuadrant defines model for RadarQuadrant.
type RadarQuadrant struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// AddRadarJSONBody defines parameters for AddRadar.
type AddRadarJSONBody struct {
	Title string `json:"title"`
}

// UpdateRadarJSONBody defines parameters for UpdateRadar.
type UpdateRadarJSONBody struct {
	Title string `json:"title"`
}

// AddRadarItemJSONBody defines parameters for AddRadarItem.
type AddRadarItemJSONBody struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	QuadrantId  string `json:"quadrant_id"`
}

// UpdateRadarItemJSONBody defines parameters for UpdateRadarItem.
type UpdateRadarItemJSONBody struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	QuadrantId  string `json:"quadrant_id"`
}

// AddRadarJSONRequestBody defines body for AddRadar for application/json ContentType.
type AddRadarJSONRequestBody AddRadarJSONBody

// UpdateRadarJSONRequestBody defines body for UpdateRadar for application/json ContentType.
type UpdateRadarJSONRequestBody UpdateRadarJSONBody

// AddRadarItemJSONRequestBody defines body for AddRadarItem for application/json ContentType.
type AddRadarItemJSONRequestBody AddRadarItemJSONBody

// UpdateRadarItemJSONRequestBody defines body for UpdateRadarItem for application/json ContentType.
type UpdateRadarItemJSONRequestBody UpdateRadarItemJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /api/v1/radars)
	AddRadar(c *gin.Context)

	// (DELETE /api/v1/radars/{radar_id})
	DeleteRadar(c *gin.Context, radarId string)

	// (GET /api/v1/radars/{radar_id})
	GetRadarById(c *gin.Context, radarId string)

	// (PUT /api/v1/radars/{radar_id})
	UpdateRadar(c *gin.Context, radarId string)

	// (GET /api/v1/radars/{radar_id}/items)
	GetRadarItems(c *gin.Context, radarId string)

	// (POST /api/v1/radars/{radar_id}/items)
	AddRadarItem(c *gin.Context, radarId string)

	// (PUT /api/v1/radars/{radar_id}/items/{radar_item_id})
	UpdateRadarItem(c *gin.Context, radarId string, radarItemId string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// AddRadar operation middleware
func (siw *ServerInterfaceWrapper) AddRadar(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.AddRadar(c)
}

// DeleteRadar operation middleware
func (siw *ServerInterfaceWrapper) DeleteRadar(c *gin.Context) {

	var err error

	// ------------- Path parameter "radar_id" -------------
	var radarId string

	err = runtime.BindStyledParameterWithOptions("simple", "radar_id", c.Param("radar_id"), &radarId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter radar_id: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.DeleteRadar(c, radarId)
}

// GetRadarById operation middleware
func (siw *ServerInterfaceWrapper) GetRadarById(c *gin.Context) {

	var err error

	// ------------- Path parameter "radar_id" -------------
	var radarId string

	err = runtime.BindStyledParameterWithOptions("simple", "radar_id", c.Param("radar_id"), &radarId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter radar_id: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetRadarById(c, radarId)
}

// UpdateRadar operation middleware
func (siw *ServerInterfaceWrapper) UpdateRadar(c *gin.Context) {

	var err error

	// ------------- Path parameter "radar_id" -------------
	var radarId string

	err = runtime.BindStyledParameterWithOptions("simple", "radar_id", c.Param("radar_id"), &radarId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter radar_id: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.UpdateRadar(c, radarId)
}

// GetRadarItems operation middleware
func (siw *ServerInterfaceWrapper) GetRadarItems(c *gin.Context) {

	var err error

	// ------------- Path parameter "radar_id" -------------
	var radarId string

	err = runtime.BindStyledParameterWithOptions("simple", "radar_id", c.Param("radar_id"), &radarId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter radar_id: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetRadarItems(c, radarId)
}

// AddRadarItem operation middleware
func (siw *ServerInterfaceWrapper) AddRadarItem(c *gin.Context) {

	var err error

	// ------------- Path parameter "radar_id" -------------
	var radarId string

	err = runtime.BindStyledParameterWithOptions("simple", "radar_id", c.Param("radar_id"), &radarId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter radar_id: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.AddRadarItem(c, radarId)
}

// UpdateRadarItem operation middleware
func (siw *ServerInterfaceWrapper) UpdateRadarItem(c *gin.Context) {

	var err error

	// ------------- Path parameter "radar_id" -------------
	var radarId string

	err = runtime.BindStyledParameterWithOptions("simple", "radar_id", c.Param("radar_id"), &radarId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter radar_id: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Path parameter "radar_item_id" -------------
	var radarItemId string

	err = runtime.BindStyledParameterWithOptions("simple", "radar_item_id", c.Param("radar_item_id"), &radarItemId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter radar_item_id: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.UpdateRadarItem(c, radarId, radarItemId)
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

	router.POST(options.BaseURL+"/api/v1/radars", wrapper.AddRadar)
	router.DELETE(options.BaseURL+"/api/v1/radars/:radar_id", wrapper.DeleteRadar)
	router.GET(options.BaseURL+"/api/v1/radars/:radar_id", wrapper.GetRadarById)
	router.PUT(options.BaseURL+"/api/v1/radars/:radar_id", wrapper.UpdateRadar)
	router.GET(options.BaseURL+"/api/v1/radars/:radar_id/items", wrapper.GetRadarItems)
	router.POST(options.BaseURL+"/api/v1/radars/:radar_id/items", wrapper.AddRadarItem)
	router.PUT(options.BaseURL+"/api/v1/radars/:radar_id/items/:radar_item_id", wrapper.UpdateRadarItem)
}

type AddRadarRequestObject struct {
	Body *AddRadarJSONRequestBody
}

type AddRadarResponseObject interface {
	VisitAddRadarResponse(w http.ResponseWriter) error
}

type AddRadar201JSONResponse Radar

func (response AddRadar201JSONResponse) VisitAddRadarResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	return json.NewEncoder(w).Encode(response)
}

type AddRadar400JSONResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (response AddRadar400JSONResponse) VisitAddRadarResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type DeleteRadarRequestObject struct {
	RadarId string `json:"radar_id"`
}

type DeleteRadarResponseObject interface {
	VisitDeleteRadarResponse(w http.ResponseWriter) error
}

type DeleteRadar200JSONResponse struct {
	Success bool `json:"success"`
}

func (response DeleteRadar200JSONResponse) VisitDeleteRadarResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type DeleteRadar404JSONResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (response DeleteRadar404JSONResponse) VisitDeleteRadarResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type GetRadarByIdRequestObject struct {
	RadarId string `json:"radar_id"`
}

type GetRadarByIdResponseObject interface {
	VisitGetRadarByIdResponse(w http.ResponseWriter) error
}

type GetRadarById200JSONResponse Radar

func (response GetRadarById200JSONResponse) VisitGetRadarByIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetRadarById400JSONResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (response GetRadarById400JSONResponse) VisitGetRadarByIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type GetRadarById404JSONResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (response GetRadarById404JSONResponse) VisitGetRadarByIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type UpdateRadarRequestObject struct {
	RadarId string `json:"radar_id"`
	Body    *UpdateRadarJSONRequestBody
}

type UpdateRadarResponseObject interface {
	VisitUpdateRadarResponse(w http.ResponseWriter) error
}

type UpdateRadar200JSONResponse Radar

func (response UpdateRadar200JSONResponse) VisitUpdateRadarResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type UpdateRadar400JSONResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (response UpdateRadar400JSONResponse) VisitUpdateRadarResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type UpdateRadar404JSONResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (response UpdateRadar404JSONResponse) VisitUpdateRadarResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type GetRadarItemsRequestObject struct {
	RadarId string `json:"radar_id"`
}

type GetRadarItemsResponseObject interface {
	VisitGetRadarItemsResponse(w http.ResponseWriter) error
}

type GetRadarItems200JSONResponse RadarItem

func (response GetRadarItems200JSONResponse) VisitGetRadarItemsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetRadarItems400JSONResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (response GetRadarItems400JSONResponse) VisitGetRadarItemsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type AddRadarItemRequestObject struct {
	RadarId string `json:"radar_id"`
	Body    *AddRadarItemJSONRequestBody
}

type AddRadarItemResponseObject interface {
	VisitAddRadarItemResponse(w http.ResponseWriter) error
}

type AddRadarItem201JSONResponse RadarItem

func (response AddRadarItem201JSONResponse) VisitAddRadarItemResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	return json.NewEncoder(w).Encode(response)
}

type AddRadarItem400JSONResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (response AddRadarItem400JSONResponse) VisitAddRadarItemResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type UpdateRadarItemRequestObject struct {
	RadarId     string `json:"radar_id"`
	RadarItemId string `json:"radar_item_id"`
	Body        *UpdateRadarItemJSONRequestBody
}

type UpdateRadarItemResponseObject interface {
	VisitUpdateRadarItemResponse(w http.ResponseWriter) error
}

type UpdateRadarItem200JSONResponse RadarItem

func (response UpdateRadarItem200JSONResponse) VisitUpdateRadarItemResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type UpdateRadarItem400JSONResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (response UpdateRadarItem400JSONResponse) VisitUpdateRadarItemResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type UpdateRadarItem404JSONResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (response UpdateRadarItem404JSONResponse) VisitUpdateRadarItemResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {

	// (POST /api/v1/radars)
	AddRadar(ctx context.Context, request AddRadarRequestObject) (AddRadarResponseObject, error)

	// (DELETE /api/v1/radars/{radar_id})
	DeleteRadar(ctx context.Context, request DeleteRadarRequestObject) (DeleteRadarResponseObject, error)

	// (GET /api/v1/radars/{radar_id})
	GetRadarById(ctx context.Context, request GetRadarByIdRequestObject) (GetRadarByIdResponseObject, error)

	// (PUT /api/v1/radars/{radar_id})
	UpdateRadar(ctx context.Context, request UpdateRadarRequestObject) (UpdateRadarResponseObject, error)

	// (GET /api/v1/radars/{radar_id}/items)
	GetRadarItems(ctx context.Context, request GetRadarItemsRequestObject) (GetRadarItemsResponseObject, error)

	// (POST /api/v1/radars/{radar_id}/items)
	AddRadarItem(ctx context.Context, request AddRadarItemRequestObject) (AddRadarItemResponseObject, error)

	// (PUT /api/v1/radars/{radar_id}/items/{radar_item_id})
	UpdateRadarItem(ctx context.Context, request UpdateRadarItemRequestObject) (UpdateRadarItemResponseObject, error)
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

// AddRadar operation middleware
func (sh *strictHandler) AddRadar(ctx *gin.Context) {
	var request AddRadarRequestObject

	var body AddRadarJSONRequestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Status(http.StatusBadRequest)
		ctx.Error(err)
		return
	}
	request.Body = &body

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.AddRadar(ctx, request.(AddRadarRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "AddRadar")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(AddRadarResponseObject); ok {
		if err := validResponse.VisitAddRadarResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}

// DeleteRadar operation middleware
func (sh *strictHandler) DeleteRadar(ctx *gin.Context, radarId string) {
	var request DeleteRadarRequestObject

	request.RadarId = radarId

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.DeleteRadar(ctx, request.(DeleteRadarRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "DeleteRadar")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(DeleteRadarResponseObject); ok {
		if err := validResponse.VisitDeleteRadarResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetRadarById operation middleware
func (sh *strictHandler) GetRadarById(ctx *gin.Context, radarId string) {
	var request GetRadarByIdRequestObject

	request.RadarId = radarId

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetRadarById(ctx, request.(GetRadarByIdRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetRadarById")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(GetRadarByIdResponseObject); ok {
		if err := validResponse.VisitGetRadarByIdResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}

// UpdateRadar operation middleware
func (sh *strictHandler) UpdateRadar(ctx *gin.Context, radarId string) {
	var request UpdateRadarRequestObject

	request.RadarId = radarId

	var body UpdateRadarJSONRequestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Status(http.StatusBadRequest)
		ctx.Error(err)
		return
	}
	request.Body = &body

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.UpdateRadar(ctx, request.(UpdateRadarRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "UpdateRadar")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(UpdateRadarResponseObject); ok {
		if err := validResponse.VisitUpdateRadarResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetRadarItems operation middleware
func (sh *strictHandler) GetRadarItems(ctx *gin.Context, radarId string) {
	var request GetRadarItemsRequestObject

	request.RadarId = radarId

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetRadarItems(ctx, request.(GetRadarItemsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetRadarItems")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(GetRadarItemsResponseObject); ok {
		if err := validResponse.VisitGetRadarItemsResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}

// AddRadarItem operation middleware
func (sh *strictHandler) AddRadarItem(ctx *gin.Context, radarId string) {
	var request AddRadarItemRequestObject

	request.RadarId = radarId

	var body AddRadarItemJSONRequestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Status(http.StatusBadRequest)
		ctx.Error(err)
		return
	}
	request.Body = &body

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.AddRadarItem(ctx, request.(AddRadarItemRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "AddRadarItem")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(AddRadarItemResponseObject); ok {
		if err := validResponse.VisitAddRadarItemResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}

// UpdateRadarItem operation middleware
func (sh *strictHandler) UpdateRadarItem(ctx *gin.Context, radarId string, radarItemId string) {
	var request UpdateRadarItemRequestObject

	request.RadarId = radarId
	request.RadarItemId = radarItemId

	var body UpdateRadarItemJSONRequestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Status(http.StatusBadRequest)
		ctx.Error(err)
		return
	}
	request.Body = &body

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.UpdateRadarItem(ctx, request.(UpdateRadarItemRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "UpdateRadarItem")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(UpdateRadarItemResponseObject); ok {
		if err := validResponse.VisitUpdateRadarItemResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+RYXW/yNhT+K5G3i03KSGh7MXG1sa4VqraqXSutqlDlxgcwIrZrO+0Yyn9/ZTsB8sFn",
	"A1R6r7Ac+/g5z3N8zjEzFPFYcAZMK9SZIQlvCSjd5YSCnbjHBMvfCbl3H8xUxJkGZodYiAmNsKacBWPF",
	"mZlT0QhibEZCcgFSZ5Y01RMwA/gPx8IM0V9TTxr7yEd6KsyM0pKyIUpT30KhEgjqPGd7+/Nl/HUMkUZp",
	"cZ2WCaS+g9zTEDcCm4CKJBVmrfWiiNNHDMdQ++EtwURipl8oqflectBa8QuHFS3s6vujIFjDd+l+U64f",
	"KWDtjBKcKXfsH1xCF+8VuT9KGKAO+iFYXOrAfVXBn1Jy6QAUNEUPI/Cya++ZkzBlyqPsHU8o8SgTiVZG",
	"TwPrb66veMLIcUApAREdUCCeBMUTGYH3gZXHuPYGFkWG6p8kikCpTyitFhbWaf3K+QQwq4idb6+Xu+jY",
	"rQBpQXkDyqgaAfE+qB55mQ3vJ8a9zIufUR7RNpG5EGmMemu4DmFOZ374NejutHd5OgDzVH4IAMb4tiAO",
	"ysQuQPIMd0och8Sw9vzUz0zZ++oySOVGR5zUl6UYlMJD2FyS8oW+s1W92xkPNWdLwBrIC7ZUDLiMzQgZ",
	"un7R1Ba6CixXJBfJJ0zGSn08jHrXD49P8dnNr1fj29vBv/qpfVO7XUNsz54PttR5bgtLiafLdXtHc3fZ",
	"tjqTOxVSHyU2snbhryQdJSg/1F9Wo2B62dOcwJUiW7IaEXpTP/XJQNjYju0sZjNqrG3w1qi0UpG7JX+K",
	"qhyGwfqGdWVnStmA5zkRRxYlxJhOUCef+u1VJoxPEvp/i4FG1Q4IopG7I16MGR5CDJap7DK57/fZHXoH",
	"qdy2ditshcYaF8CwoKiDzlth6xz5SGA9sgwFWNDgvR1Y63ZGcNdp8rw76RHUQaboZgcs3oXTVRFUeDoG",
	"5Xdjucs9C9urDWXrgkrzk/roIgw3byz10EYRjYfKKJc53TdzRSKCmf19oSR1L58JaKiycmnnc2IEljgG",
	"DYbH5xmiRgNDdB7zHZQbReXO318qieVo61fo2tLrpTJ9EV5st2fe1ld4Qv3UR0OoiY1r0JaC7rRHvhgN",
	"tV3rvpHTHI0iqaEx66KOEEz73N/iSzrdW4xSr3hqLdbd/GDe8qwN+55d9QXjvu6d0mTW9DfUCtspfb1I",
	"Lv0duH81Kj9Gj1WRXFzOJzTEeZ3alFcOq4i/zpRDeTyFm8pXNW/sk+esNP0WAAD//40uAZysFwAA",
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
