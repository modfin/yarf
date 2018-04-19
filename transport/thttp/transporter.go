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
	STD_PORT     string = "23456"
	STD_PROTOCOL string = "http"
)

type HttpTransporter struct {
	options   Options
	functions map[string]func(requestData []byte) (responseData []byte)

	initServer sync.Once
	server     *http.Server
	mux        *http.ServeMux
}

type Options struct {
	Server    Server
	Client    Client
	Discovery Discovery
}
type Server struct {
	Addr      string
	TLSConfig *tls.Config
	timeout   time.Duration
}
type Client struct {
	Timeout time.Duration
}

func NewHttpTransporter(options Options) (*HttpTransporter, error) {
	t := HttpTransporter{
		options: options,
		mux:     http.NewServeMux(),
	}

	return &t, nil
}

func (h *HttpTransporter) Call(ctx context.Context, function string, requestData []byte) (response []byte, err error) {
	//ctx, _ = context.WithTimeout(ctx, durationOr(h.options.Client.Timeout, 10 * time.Second))

	url, err := h.options.Discovery.Url()

	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(requestData)
	resp, err := http.Post(url+"/"+function, "application/octet-stream", r)

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

func (h *HttpTransporter) Start() error {

	h.server = &http.Server{
		Addr:           stringOr(h.options.Server.Addr, ":"+STD_PORT),
		TLSConfig:      h.options.Server.TLSConfig,
		Handler:        h.mux,
		ReadTimeout:    durationOr(h.options.Server.timeout, 10*time.Second),
		WriteTimeout:   durationOr(h.options.Server.timeout, 10*time.Second),
		MaxHeaderBytes: 1 << 20,
	}

	return h.server.ListenAndServe()
}

func (h *HttpTransporter) Stop(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}

func (h *HttpTransporter) Listen(function string, toExec func(requestData []byte) (responseData []byte)) (err error) {

	h.mux.HandleFunc("/"+function, func(res http.ResponseWriter, req *http.Request) {
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

		respData := toExec(reqData)
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
