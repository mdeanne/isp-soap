package soap

import (
	"encoding/xml"
	"github.com/valyala/fasthttp"
)

type CallResponse struct {
	HttpStatusCode int
	Response       Envelope
}

func (r CallResponse) IsSuccess() bool {
	return r.HttpStatusCode == fasthttp.StatusOK && r.Response.Body.Fault == nil
}

func (r CallResponse) SoapFault() *SOAPFault {
	return r.Response.Body.Fault
}

func (r CallResponse) SoapContent() []byte {
	return r.Response.Body.Content
}

func (r CallResponse) SoapHeader() *Header {
	return r.Response.Header
}

func (r CallResponse) UnmarshalContent(ptr interface{}) error {
	return xml.Unmarshal(r.SoapContent(), ptr)
}
