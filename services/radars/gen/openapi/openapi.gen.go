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

// Error defines model for Error.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Radar defines model for Radar.
type Radar struct {
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	UniqId    string    `json:"uniq_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AddRadarOut defines model for AddRadarOut.
type AddRadarOut = Radar

// BadRequest defines model for BadRequest.
type BadRequest = Error

// GetRadarByIdOut defines model for GetRadarByIdOut.
type GetRadarByIdOut = Radar

// NotFound defines model for NotFound.
type NotFound = Error

// Operation defines model for Operation.
type Operation struct {
	Success bool `json:"success"`
}

// UpdateRadarOut defines model for UpdateRadarOut.
type UpdateRadarOut = Radar

// AddRadar defines model for AddRadar.
type AddRadar struct {
	Title string `json:"title"`
}

// UpdateRadar defines model for UpdateRadar.
type UpdateRadar struct {
	Title string `json:"title"`
}

// AddRadarJSONBody defines parameters for AddRadar.
type AddRadarJSONBody struct {
	Title string `json:"title"`
}

// UpdateRadarJSONBody defines parameters for UpdateRadar.
type UpdateRadarJSONBody struct {
	Title string `json:"title"`
}

// AddRadarJSONRequestBody defines body for AddRadar for application/json ContentType.
type AddRadarJSONRequestBody AddRadarJSONBody

// UpdateRadarJSONRequestBody defines body for UpdateRadar for application/json ContentType.
type UpdateRadarJSONRequestBody UpdateRadarJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Add a new radar
	// (POST /api/v1/radars)
	AddRadar(c *gin.Context)
	// Deletes a radar
	// (DELETE /api/v1/radars/{radar_id})
	DeleteRadar(c *gin.Context, radarId string)
	// Find radar by ID
	// (GET /api/v1/radars/{radar_id})
	GetRadarById(c *gin.Context, radarId string)
	// Update existing radar
	// (PUT /api/v1/radars/{radar_id})
	UpdateRadar(c *gin.Context, radarId string)
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
}

type AddRadarOutJSONResponse Radar

type BadRequestJSONResponse Error

type GetRadarByIdOutJSONResponse Radar

type NotFoundJSONResponse Error

type OperationJSONResponse struct {
	Success bool `json:"success"`
}

type UpdateRadarOutJSONResponse Radar

type AddRadarRequestObject struct {
	Body *AddRadarJSONRequestBody
}

type AddRadarResponseObject interface {
	VisitAddRadarResponse(w http.ResponseWriter) error
}

type AddRadar201JSONResponse struct{ AddRadarOutJSONResponse }

func (response AddRadar201JSONResponse) VisitAddRadarResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	return json.NewEncoder(w).Encode(response)
}

type AddRadar400JSONResponse struct{ BadRequestJSONResponse }

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

type DeleteRadar200JSONResponse struct{ OperationJSONResponse }

func (response DeleteRadar200JSONResponse) VisitDeleteRadarResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type DeleteRadar404JSONResponse struct{ NotFoundJSONResponse }

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

type GetRadarById200JSONResponse struct{ GetRadarByIdOutJSONResponse }

func (response GetRadarById200JSONResponse) VisitGetRadarByIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetRadarById400JSONResponse struct{ BadRequestJSONResponse }

func (response GetRadarById400JSONResponse) VisitGetRadarByIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type GetRadarById404JSONResponse struct{ NotFoundJSONResponse }

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

type UpdateRadar200JSONResponse struct{ UpdateRadarOutJSONResponse }

func (response UpdateRadar200JSONResponse) VisitUpdateRadarResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type UpdateRadar400JSONResponse struct{ BadRequestJSONResponse }

func (response UpdateRadar400JSONResponse) VisitUpdateRadarResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type UpdateRadar404JSONResponse struct{ NotFoundJSONResponse }

func (response UpdateRadar404JSONResponse) VisitUpdateRadarResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Add a new radar
	// (POST /api/v1/radars)
	AddRadar(ctx context.Context, request AddRadarRequestObject) (AddRadarResponseObject, error)
	// Deletes a radar
	// (DELETE /api/v1/radars/{radar_id})
	DeleteRadar(ctx context.Context, request DeleteRadarRequestObject) (DeleteRadarResponseObject, error)
	// Find radar by ID
	// (GET /api/v1/radars/{radar_id})
	GetRadarById(ctx context.Context, request GetRadarByIdRequestObject) (GetRadarByIdResponseObject, error)
	// Update existing radar
	// (PUT /api/v1/radars/{radar_id})
	UpdateRadar(ctx context.Context, request UpdateRadarRequestObject) (UpdateRadarResponseObject, error)
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

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9SW34rjNhTGX0Wc9qIFNU5296LkqjtMZwlLG0hnocsSBsU6ThRsSSvJk00Hv3uRZGds",
	"xyHNdEjplYX+HH36fUc6foJUFVpJlM7C9AkMfi3RuhvFBYaO95wvGGfGt1MlHUrnm0zrXKTMCSWTrVXS",
	"99l0gwXzLW2URuPqEE64HH0Dv7FC+yb8ticmRKXg9tr3WGeEXENV0aBBGOQw/VKvXR6mqdUWUwdVd54z",
	"JVYUPmnOHP5P5IYeq5W0XdDz0l0k/nuDGUzhu+TZyCSO2iSiCLtztKkR2geBKfxRpilaCxWFG8YX0fRX",
	"2/ZXY9TgtvcbJHWGEb8TE9ISIR9ZLjgRUpcuSPqALii/2c/4tXH8rtydKiW/DgyrMRWZQE4MWlWaFMmO",
	"WSKVI1lQUVGYazQsrnlxTtv6gGeyeqVUjkwepXWzfDixu8c6yCWZkMJukJOdcBtSxyA/SEXqU/wI3Ut7",
	"Xa8rWocKWKJNR+BSxcNb0Lv2FAq0lq2Hxnrsmok0xjpGSOHwYvX2Nsgc8gcWUGTKFL4FHtdPThR4/BrR",
	"yx4vCqUUXx8E7y4Yl1trd/eb2Yf7T5+LNx9/vtvO59mf7vPk42CMYOAlMnuEGhGNeto+eSf+cP4Jmakm",
	"b1gaVGDBRA7TpuuXlSmlykvx10hiQN67iphuIiRSMMnWWPj8O9CM44sa4iMaG5dNRuPR2EdTGiXTAqbw",
	"djQevQUKmrlNMDFhWiSPkyREDz1axae2K+E954QRibuDWaq5SDMexxsBzzV6f+oWdMp4cljbrzlvxpPT",
	"Eep5SbswVRTejcfn17SKSrhnZVEwsx88pmNr67OgBrT087vQkqfwfRC8itxydHhM8Db0E3YCYBxuGGpm",
	"WIEOvSVf+pHCJCI4cYrUu/ksg2mwFShIVvi0aGRBv7jT1iPVz/3lkQf/gOdzDQgOvDu/4lDHuvwjBdui",
	"1OEPy4rCGgfyc4GuNNIvtEKuczxBuV26z2Ge3RKV1dfOKWLCDv816f6/x0sy/t9ZdCckr6ms9mR2O+iR",
	"Lgc8iqWU4DdhhRNyfcKj9m/yGYvmpYt/JUHGqxpz4RvWFl29xNjef8b1fW250zGnb25VVX8HAAD//0gA",
	"JBaWDQAA",
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