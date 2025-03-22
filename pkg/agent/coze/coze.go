package coze

import (
	"github.com/coze-dev/coze-go"
)

type CozeClient struct {
	Client  coze.CozeAPI
	BaseURL string
	auth    coze.Auth
}

type Options func(*CozeClient)

// 初始化API
func NewCozeClient(token string, opts ...Options) *CozeClient {
	auth := coze.NewTokenAuth(token)
	client := &CozeClient{
		Client:  coze.NewCozeAPI(auth, coze.WithBaseURL(coze.CnBaseURL)),
		BaseURL: coze.ComBaseURL,
		auth:    coze.NewTokenAuth(token),
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// 更换BaseURL，默认是api.coze.cn
func WithBaseURL(baseURL string) Options {
	return func(c *CozeClient) {
		c.Client = coze.NewCozeAPI(c.auth, coze.WithBaseURL(baseURL))
		c.BaseURL = baseURL
	}
}
