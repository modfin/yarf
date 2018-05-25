package thttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const (
	//StdPort defines the standard port for thttp yarf
	StdPort string = "23456"

	//StdProtocol defines the standard protocol for thttp yarf
	StdProtocol string = "http"
)

// HTTPTransporter implements the yarf.Transport for using http as a transport protocol
type HTTPTransporter struct {
	options   Options
	functions map[string]func(requestData []byte) (responseData []byte)

	initServer sync.Once
	server     *http.Server
	mux        *http.ServeMux
}

// Options defines the options used by the http yarf transport
type Options struct {
	Server    Server
	Client    Client
	Discovery Discovery
}

// Server defines the server config used
type Server struct {
	Addr      string
	TLSConfig *tls.Config
	timeout   time.Duration
}

// Client defines the server config used
type Client struct {
	Timeout time.Duration
}

// NewHTTPTransporter a constructor for the HTTPTransporter
func NewHTTPTransporter(options Options) (*HTTPTransporter, error) {
	t := HTTPTransporter{
		options: options,
		mux:     http.NewServeMux(),
	}

	return &t, nil
}

// Call implements client side call of transporter
func (h *HTTPTransporter) Call(ctx context.Context, function string, requestData []byte) (response []byte, err error) {

	url, err := h.options.Discovery.URL()

	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(requestData)

	req, err := http.NewRequest("POST", url+"/"+function, r)
	req = req.WithContext(ctx)
	req.Header.Set("WithContent-Type", "application/octet-stream")

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer func() { resp.Body.Close() }()

	return ioutil.ReadAll(resp.Body)
}

// Start initiates the http server to receive requests
func (h *HTTPTransporter) Start() error {

	h.server = &http.Server{
		Addr:           stringOr(h.options.Server.Addr, ":"+StdPort),
		TLSConfig:      h.options.Server.TLSConfig,
		Handler:        h.mux,
		ReadTimeout:    durationOr(h.options.Server.timeout, 10*time.Second),
		WriteTimeout:   durationOr(h.options.Server.timeout, 10*time.Second),
		MaxHeaderBytes: 1 << 20,
	}

	return h.server.ListenAndServe()
}

// Close halts the http server to receive requests
func (h *HTTPTransporter) Close() error {
	return h.server.Shutdown(context.Background())
}

// Listen defines the function that will handle yarf requests
func (h *HTTPTransporter) Listen(function string, toExec func(ctx context.Context, requestData []byte) (responseData []byte)) (err error) {

	h.mux.HandleFunc("/"+function, func(res http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		if req.Method != "POST" {
			//FIX returning error
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		reqData, err := ioutil.ReadAll(req.Body)

		if err != nil {
			//FIX returning error
			res.WriteHeader(http.StatusInsufficientStorage)
			return
		}

		if err != nil {
			//FIX returning error
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		//fmt.Println("CTX", req.Context())
		respData := toExec(req.Context(), reqData)

		res.WriteHeader(http.StatusOK)
		res.Header().Set("content-type", "application/octet-stream")
		_, err = res.Write(respData)

		if err != nil {
			//FIX returning error
			return
		}

	})

	return
}
