package soap

import "github.com/valyala/fasthttp"

type ClientOption func(c *Client)

func WithFastHttpClient(cli *fasthttp.Client) ClientOption {
	return func(c *Client) {
		c.cli = cli
	}
}

func WithHttpHeaders(headers map[string]string) ClientOption {
	return func(c *Client) {
		c.httpHeaders = headers
	}
}
