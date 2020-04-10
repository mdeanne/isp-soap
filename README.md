# isp-soap
Simple fasthttp client wrapping for SOAP.
## Example
```go
func main() {
	client := soap.NewClient("http://wsgeoip.lavasoft.com/ipservice.asmx")
	ip, err := getOutboundIp()
	if err != nil {
		panic(err)
	}
	req := GetIpLocation{
		SIp: ip,
	}
	//soap.Xml() - if you have xml bytes already
	response, err := client.Call("http://lavasoft.com/GetIpLocation", soap.Any(req), soap.WithTimeout(1*time.Second))
	if err != nil {
		panic(err) //network error or xml.UnmarshalError maybe occurred
	}
	if !response.IsSuccess() { //check http status and SOAP fault
		panic(string(response.HTTP().Body())) //access to original http response body
	}
	resp := GetIpLocationResponse{}
	err = response.UnmarshalBody(&resp) //unmarshal SOAP body to expected struct
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
```
