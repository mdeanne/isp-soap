package soap

import (
	"bytes"
	"encoding/xml"
)

type Envelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Header  *Header
	Body    Body
}

type Header struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`

	Items []interface{} `xml:",omitempty"`
}

type Body struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`

	Fault   *SOAPFault `xml:",omitempty"`
	Content []byte     `xml:",innerxml"`
}

// UnmarshalXML unmarshals SOAPBody xml
func (b *Body) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	val := struct {
		Value []byte `xml:",innerxml"`
	}{}
	if err := d.DecodeElement(&val, &start); err != nil {
		return nil
	}
	token, err := xml.NewDecoder(bytes.NewReader(val.Value)).Token()
	if err != nil {
		return err
	}
	switch se := token.(type) {
	case xml.StartElement:
		if se.Name.Space == "http://schemas.xmlsoap.org/soap/envelope/" && se.Name.Local == "Fault" {
			fault := &SOAPFault{}
			if err := xml.Unmarshal(val.Value, fault); err != nil {
				return err
			}
			b.Fault = fault
			return nil
		}
	}
	b.Content = val.Value
	return nil
}

type SOAPFault struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`

	Code   string `xml:"faultcode,omitempty"`
	String string `xml:"faultstring,omitempty"`
	Actor  string `xml:"faultactor,omitempty"`
	Detail string `xml:"detail,omitempty"`
}

func (f *SOAPFault) Error() string {
	return f.String
}

type RequestBody interface {
	Xml() ([]byte, error)
}

type xmlReqBody struct {
	body []byte
}

func (b xmlReqBody) Xml() ([]byte, error) {
	return b.body, nil
}

func Xml(body []byte) RequestBody {
	return xmlReqBody{
		body: body,
	}
}

type anyReqBody struct {
	body interface{}
}

func (s anyReqBody) Xml() ([]byte, error) {
	return xml.Marshal(s.body)
}

func Any(body interface{}) RequestBody {
	return anyReqBody{body: body}
}
