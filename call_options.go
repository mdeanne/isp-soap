package soap

import "time"

type CallOption func(opts *callOptions)

type callOptions struct {
	timeout     time.Duration
	httpHeaders map[string]string
}

func WithTimeout(timeout time.Duration) CallOption {
	return func(opts *callOptions) {
		opts.timeout = timeout
	}
}

func WithCallHttpHeaders(headers map[string]string) CallOption {
	return func(opts *callOptions) {
		opts.httpHeaders = headers
	}
}

func defaultOptions() *callOptions {
	return &callOptions{}
}
