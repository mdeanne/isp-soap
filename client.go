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

func (c *Client) Call(action string, body RequestBody, opts ...CallOption) (CallResponse, error) {
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
	resp := CallResponse{}
	if body != nil {
		bytes, err := body.Xml()
		if err != nil {
			return resp, fmt.Errorf("prepare request body: %v", err)
		}
		reqEnv.Body.Content = bytes
	}
	bytes, err := xml.Marshal(reqEnv)
	if err != nil {
		return resp, fmt.Errorf("marhal envelope: %v", err)
	}

	request.Header.SetContentLength(len(bytes))
	request.SetBody(bytes)

	if options.timeout > 0 {
		err = c.cli.DoTimeout(request, response, options.timeout)
	} else {
		err = c.cli.Do(request, response)
	}
	if err != nil {
		return resp, err //internal or network error
	}

	resEnv := Envelope{}
	copied := make([]byte, len(response.Body()))
	copy(copied, response.Body())
	resp.http.httpStatusCode = response.StatusCode()
	resp.http.body = copied
	if err := xml.Unmarshal(copied, &resEnv); err != nil {
		return resp, xml.UnmarshalError(fmt.Sprintf("unmarshal enveloper: %v", err))
	}
	resp.envelope = resEnv
	return resp, nil
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
