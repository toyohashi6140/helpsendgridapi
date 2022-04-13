package request

import (
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
)

const (
	ENDPOINT string = "/v3/mail/send"
	HOST     string = "https://api.sendgrid.com"
)

type request rest.Request

func New(apiKey string) *request {
	r := request(sendgrid.GetRequest(apiKey, ENDPOINT, HOST))
	return &r
}

// Post request methodにpostを指定
func (r *request) Post() *request {
	r.Method = rest.Post
	return r
}

// Get request methodにgetを指定
func (r *request) Get() *request {
	r.Method = rest.Get
	return r
}

// Do jsonを含めたあとで実際のリクエストを行う
func (r *request) Do() (*rest.Response, error) {
	res, err := sendgrid.API(rest.Request(*r))
	if err != nil {
		return nil, err
	}
	return res, nil
}
