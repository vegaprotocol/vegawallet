package jsonrpc

import (
	"fmt"
	"strings"
)

const JSONRPC2 string = "2.0"

type API struct {
	// commands maps a method to a command.
	commands map[string]Command
}

func New() *API {
	return &API{
		commands: map[string]Command{},
	}
}

func (a *API) DispatchRequest(request *Request) *Response {
	if err := request.Check(); err != nil {
		return invalidRequestResponse(request, err)
	}

	commands, ok := a.commands[request.Method]
	if !ok {
		return unsupportedMethodResponse(request)
	}

	result, errorDetails := commands.Handle(request.Params)
	if errorDetails != nil {
		return NewErrorResponse(request.ID, errorDetails)
	}
	return NewSuccessfulResponse(request.ID, result)
}

func (a *API) RegisterCommand(method string, handler Command) {
	if len(strings.Trim(method, " \t\r\n")) == 0 {
		panic("method cannot be empty")
	}

	if handler == nil {
		panic("handler cannot be nil")
	}

	if _, ok := a.commands[method]; ok {
		panic(fmt.Sprintf("method \"%s\" is already registered", method))
	}

	a.commands[method] = handler
}

func invalidRequestResponse(request *Request, err error) *Response {
	return NewErrorResponse(request.ID, NewInvalidRequest(err))
}

func unsupportedMethodResponse(request *Request) *Response {
	return NewErrorResponse(request.ID, NewMethodNotFound(fmt.Sprintf("method \"%s\" is not supported", request.Method)))
}
