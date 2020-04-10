package soap

import (
	"encoding/xml"
	"github.com/valyala/fasthttp"
)

type SOAPPart interface {
	Fault() *SOAPFault
	Body() []byte
	Header() *Header
	UnmarshalBody(ptr interface{}) error
}

type HttpPart interface {
	StatusCode() int
	Body() []byte
}

type CallResponse struct {
	http     httpPart
	envelope Envelope
}

func (r CallResponse) Fault() *SOAPFault {
	return r.envelope.Body.Fault
}

func (r CallResponse) Body() []byte {
	return r.envelope.Body.Content
}

func (r CallResponse) Header() *Header {
	return r.envelope.Header
}

func (r CallResponse) UnmarshalBody(ptr interface{}) error {
	return xml.Unmarshal(r.Body(), ptr)
}

type httpPart struct {
	httpStatusCode int
	body           []byte
}

func (p httpPart) StatusCode() int {
	return p.httpStatusCode
}

func (p httpPart) Body() []byte {
	return p.body
}

func (r CallResponse) IsSuccess() bool {
	return r.http.httpStatusCode == fasthttp.StatusOK && r.envelope.Body.Fault == nil
}

func (r CallResponse) HTTP() HttpPart {
	return r.http
}
