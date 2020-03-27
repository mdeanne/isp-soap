package soap

import (
	"encoding/xml"
	"fmt"
	"net"
	"testing"
	"time"
)

type GetIpLocationResponse struct {
	GetIpLocationResult string
}

type GetIpLocation struct {
	XMLName xml.Name `xml:"http://lavasoft.com/ GetIpLocation"`
	SIp     string   `xml:"sIp"`
}

func TestClient_Call(t *testing.T) {
	client := NewClient("http://wsgeoip.lavasoft.com/ipservice.asmx")
	ip, err := getOutboundIp()
	if err != nil {
		panic(err)
	}
	req := GetIpLocation{
		SIp: ip,
	}
	response, err := client.Call("http://lavasoft.com/GetIpLocation", Any(req), WithTimeout(1*time.Second))
	if err != nil {
		panic(err)
	}
	if !response.IsSuccess() {
		panic(response.SoapFault())
	}
	resp := GetIpLocationResponse{}
	err = response.UnmarshalContent(&resp)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

func getOutboundIp() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP.To4().String(), nil
}
