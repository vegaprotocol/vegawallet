package jsonrpc

type Command interface {
	Handle(params Params) (Result, *ErrorDetails)
}
