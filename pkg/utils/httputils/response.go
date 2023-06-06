// Package httputils provides utilities for HTTP related operations
package httputils

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

// the constants below describes the available response codes and error codes in this system.
// the first digit (x0000) defines the category of the error
// 1 = for rendering related errors (REST response)
// 2 = validation related errors
// 3 = HTTP status related errors
// 4 = CRUD related errors
// the last three digit (00xxx) defines the incremental value of the same category and type of code
const (
	// RenderFailed is an application error code for a failed rendering.
	// a failed rendering may happen when the fetched response body is invalid
	RenderFailed = 1001

	// InvalidRequestJSON is an application error code where system is unable to extract JSON on the request body
	InvalidRequestJSON = 1002

	// RequestJSONExtractionFailed is an application error code where system is unable to unmarshal captured JSON data
	RequestJSONExtractionFailed = 1003

	// InputValidationError is an application error code where the captured input data got validation error
	InputValidationError = 2001

	// UnauthorizedAccess is an application error code where the identity has no authorize to access
	UnauthorizedAccess = 2002

	// BadRequest is an application error to represent bad request
	BadRequest = 3001

	// CreateDataFailed is an application error to represent that create process failed
	CreateDataFailed = 4001

	// UpdateDataFailed is an application error to represent that update process failed
	UpdateDataFailed = 4002

	// DeleteDataFailed is an application error to represent that delete process failed
	DeleteDataFailed = 4003
)

// responseText is list of error test for each application-level error code
var responseText = map[int]string{
	RenderFailed:                "failed to render a valid response body",
	InvalidRequestJSON:          "failed to extract request body",
	RequestJSONExtractionFailed: "failed to read JSON body from the request",

	InputValidationError: "got input validation error",
	UnauthorizedAccess:   "identity is unauthorized to access this API",

	BadRequest: "bad request",

	CreateDataFailed: "insert process failed",
	UpdateDataFailed: "update process failed",
	DeleteDataFailed: "delete process failed",
}

// ResponseText returns a text for the HTTP status code in the application level.
// It returns the empty string if the code is unknown.
func ResponseText(identifier string, code int) string {
	if identifier != "" {
		return fmt.Sprintf("[%s] %s", identifier, responseText[code])
	} else {
		return responseText[code]
	}
}

// Response renderer type for handling all sorts of http response
type Response struct {
	HTTPStatusCode int         `json:"-"`                                                                            // response response status code
	Data           interface{} `json:"data,omitempty"`                                                               // always set as empty
	MessageText    string      `json:"message,omitempty" example:"Resource not found."`                              // user-level status message
	AppErrCode     int64       `json:"code,omitempty" example:"404"`                                                 // application-specific error code
	ErrorText      string      `json:"error,omitempty" example:"The requested resource was not found on the server"` // application-level error message, for debugging
} // @name  Response

// Render implements the github.com/go-chi/render.Renderer interface for Response
func (e *Response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// httpErrPayload returns a structured error response
func httpErrPayload(errText string, appRespCode int64, httpStatusCode int, err error) render.Renderer {
	errorTxt := errText
	if err != nil {
		errorTxt = err.Error()
	}
	return &Response{
		HTTPStatusCode: httpStatusCode,
		AppErrCode:     appRespCode,
		ErrorText:      errorTxt,
		MessageText:    errText,
	}
}

// httpOKPayload returns a structured OK response
func httpOKPayload(respBody Response) render.Renderer {
	return &Response{
		HTTPStatusCode: http.StatusOK,
		Data:           respBody.Data,
		MessageText:    respBody.MessageText,
	}
}

// RenderErrResponse renders the error http response
func RenderErrResponse(w http.ResponseWriter, r *http.Request, errText string, appErrCode int64,
	httpStatusCode int, err error) {
	_ = render.Render(w, r, httpErrPayload(errText, appErrCode, httpStatusCode, err))
}

// RenderOKResponse returns a rendered http response
// rendering may fails and returns an error, otherwise it returns nil value
func RenderOKResponse(w http.ResponseWriter, r *http.Request, respBody Response) error {
	return render.Render(w, r, httpOKPayload(respBody))
}
