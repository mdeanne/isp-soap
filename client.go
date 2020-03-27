package soap

import (
	"encoding/xml"
	"fmt"
	"github.com/valyala/fasthttp"
)

const (
	contentType      = `text/xml; charset=utf-8`
	soapActionHeader = "SOAPAction"
)

type Client struct {
	cli         *fasthttp.Client
	url         string
	httpHeaders map[string]string
}

func (c *Client) Call(action string, body RequestBody, opts ...CallOption) (*CallResponse, error) {
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(request)
		fasthttp.ReleaseResponse(response)
	}()

	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	for k, v := range c.httpHeaders {
		request.Header.Set(k, v)
	}
	for k, v := range options.httpHeaders {
		request.Header.Set(k, v)
	}
	request.SetRequestURI(c.url)
	request.Header.SetMethod("POST")
	request.Header.SetContentType(contentType)
	request.Header.Set(soapActionHeader, action)

	reqEnv := Envelope{
		Body: Body{},
	}
	if body != nil {
		bytes, err := body.Xml()
		if err != nil {
			return nil, fmt.Errorf("prepare request body: %v", err)
		}
		reqEnv.Body.Content = bytes
	}
	bytes, err := xml.Marshal(reqEnv)
	if err != nil {
		return nil, fmt.Errorf("marhal envelope: %v", err)
	}

	request.Header.SetContentLength(len(bytes))
	request.SetBody(bytes)

	if options.timeout > 0 {
		err = c.cli.DoTimeout(request, response, options.timeout)
	} else {
		err = c.cli.Do(request, response)
	}
	if err != nil {
		return nil, err //internal or network error
	}

	resEnv := Envelope{}
	if err := xml.Unmarshal(response.Body(), &resEnv); err != nil {
		return nil, fmt.Errorf("unmarshal enveloper: %v", err)
	}

	return &CallResponse{
		HttpStatusCode: response.StatusCode(),
		Response:       resEnv,
	}, nil
}

func NewClient(url string, opts ...ClientOption) *Client {
	c := &Client{
		cli: &fasthttp.Client{},
		url: url,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
