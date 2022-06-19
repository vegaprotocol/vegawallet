package jsonrpc

// Result is just a nicer way to describe what's expected to be returned by the
// handlers.
type Result interface{}

type Response struct {
	// Version specifies the version of the JSON-RPC protocol.
	// MUST be exactly "2.0".
	Version string `json:"jsonrpc"`

	// Result is REQUIRED on success. This member MUST NOT exist if there was an
	// error invoking the method.
	Result Result `json:"result,omitempty"`

	// Error is REQUIRED on error. This member MUST NOT exist if there was no
	// error triggered during invocation.
	Error *ErrorDetails `json:"error,omitempty"`

	// ID is an identifier established by the Client that MUST contain a String.
	// This member is REQUIRED. It MUST be the same as the value of the id member
	// in the Request Object.
	// If there was an error in detecting the id in the Request object (e.g.
	// Parse error/Invalid Request), it MUST be empty.
	ID string `json:"id,omitempty"`
}

type ErrorCode int16

const (
	// ErrorCodeParseError Invalid JSON was received by the server. An error
	// occurred on the server while parsing the JSON text.
	ErrorCodeParseError ErrorCode = -32700
	// ErrorCodeInvalidRequest The JSON sent is not a valid Request object.
	ErrorCodeInvalidRequest ErrorCode = -3260
	// ErrorCodeMethodNotFound The method does not exist / is not available.
	ErrorCodeMethodNotFound ErrorCode = -32601
	// ErrorCodeInvalidParams Invalid method parameter(s).
	ErrorCodeInvalidParams ErrorCode = -32602
	// ErrorCodeInternalError Internal JSON-RPC error.
	ErrorCodeInternalError ErrorCode = -32603
)

// 	ErrorDetails is returned when an RPC call encounters an error.
type ErrorDetails struct {
	// Code indicates the error type that occurred.
	Code ErrorCode `json:"code"`

	// Message provides a short description of the error.
	// The message SHOULD be limited to a concise single sentence.
	Message string `json:"message"`

	// Data is a primitive or a structured value that contains additional
	// information about the error. This may be omitted.
	// The value of this member is defined by the Server (e.g. detailed error
	// information, nested errors etc.).
	Data interface{} `json:"data,omitempty"`
}

func NewParseError(data interface{}) *ErrorDetails {
	return &ErrorDetails{
		Code:    ErrorCodeParseError,
		Message: "Parse error",
		Data:    data,
	}
}

func NewInvalidRequest(data interface{}) *ErrorDetails {
	return &ErrorDetails{
		Code:    ErrorCodeInvalidRequest,
		Message: "Invalid Request",
		Data:    data,
	}
}

func NewMethodNotFound(data interface{}) *ErrorDetails {
	return &ErrorDetails{
		Code:    ErrorCodeMethodNotFound,
		Message: "Method not found",
		Data:    data,
	}
}

func NewInvalidParams(data interface{}) *ErrorDetails {
	return &ErrorDetails{
		Code:    ErrorCodeInvalidParams,
		Message: "Invalid params",
		Data:    data,
	}
}

func NewInternalError(data interface{}) *ErrorDetails {
	return &ErrorDetails{
		Code:    ErrorCodeInternalError,
		Message: "Internal error",
		Data:    data,
	}
}

func NewServerError(code ErrorCode, data interface{}) *ErrorDetails {
	if code > -32000 || code < -32099 {
		panic("server error code should be between [-32000, -32099]")
	}
	return &ErrorDetails{
		Code:    code,
		Message: "Server error",
		Data:    data,
	}
}

func NewErrorResponse(id string, details *ErrorDetails) *Response {
	return &Response{
		Version: JSONRPC2,
		Error:   details,
		ID:      id,
	}
}

func NewSuccessfulResponse(id string, result Result) *Response {
	return &Response{
		Version: JSONRPC2,
		Result:  result,
		ID:      id,
	}
}
