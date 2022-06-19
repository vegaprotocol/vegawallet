package api

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"code.vegaprotocol.io/vegawallet/libs/jsonrpc"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

var ErrCouldNotReadRequestBody = errors.New("couldn't read the HTTP request body")

type HTTPServer struct {
	server *http.Server

	api *jsonrpc.API
}

func NewHTTPServer(addr string, api *jsonrpc.API) *HTTPServer {
	httpServer := &HTTPServer{
		api: api,
	}

	router := httprouter.New()
	router.Handle(http.MethodPost, "/api/v1/request", httpServer.HandleRequestV1)

	httpServer.server = &http.Server{
		Addr:    addr,
		Handler: cors.AllowAll().Handler(router),
	}

	return httpServer
}

func (s *HTTPServer) Start() error {
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Stop() error {
	return s.server.Shutdown(context.Background())
}

func (s *HTTPServer) HandleRequestV1(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	request, err := unmarshallRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// Failing to unmarshall the request prevent us from retrieving the
		// request ID. So, it's left empty.
		s.writeResponse(w, jsonrpc.NewErrorResponse("", err))
		return
	}

	response := s.api.DispatchRequest(request)

	// If the request doesn't have an ID, it's a notification. Notifications do
	// not send content back, even id an error occurred.
	if request.IsNotification() {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if response.Error != nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	s.writeResponse(w, response)
}

func (s *HTTPServer) writeResponse(w http.ResponseWriter, response *jsonrpc.Response) {
	marshaledResponse, err := json.Marshal(response)
	if err != nil {
		return
	}

	if _, err = w.Write(marshaledResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func unmarshallRequest(r *http.Request) (*jsonrpc.Request, *jsonrpc.ErrorDetails) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, jsonrpc.NewParseError(ErrCouldNotReadRequestBody)
	}

	request := &jsonrpc.Request{}

	if len(body) == 0 {
		return request, nil
	}

	if err := json.Unmarshal(body, request); err != nil {
		var syntaxError *json.SyntaxError
		if errors.As(err, &syntaxError) {
			return nil, jsonrpc.NewParseError(err)
		}
		return nil, jsonrpc.NewInvalidRequest(err)
	}

	return request, nil
}
