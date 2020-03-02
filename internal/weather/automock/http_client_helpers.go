package automock

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/mock"
	"io"
	"io/ioutil"
	"net/http"
)

func NewClientThatChecksIfRequestHasDefinedDeadline(returnStatusCode int) *HTTPClient {
	m := HTTPClient{}
	m.On("Do", mock.MatchedBy(func(actualRequest *http.Request) bool {
		_, hasDeadline := actualRequest.Context().Deadline()
		return hasDeadline
	})).Return(&http.Response{StatusCode: returnStatusCode, Body: ioutil.NopCloser(bytes.NewBufferString("{}"))}, nil).Once()
	return &m
}

func NewClientThatReturnsWrongStatusCode() *HTTPClient {
	m := HTTPClient{}
	m.On("Do", mock.Anything).Return(&http.Response{StatusCode: http.StatusInternalServerError, Body: ioutil.NopCloser(&bytes.Buffer{})}, nil).Once()
	return &m
}

func NewClientThatReturnsError() *HTTPClient {
	m := HTTPClient{}
	m.On("Do", mock.Anything).Return(nil, fixError())
	return &m
}

func NewClientThatReturnsNotFoundStatusCode() *HTTPClient {
	m := HTTPClient{}
	m.On("Do", mock.Anything).Return(&http.Response{StatusCode: http.StatusNotFound, Body: ioutil.NopCloser(&bytes.Buffer{})}, nil).Once()
	return &m
}

func NewClientThatReturnsErrorOnClosingBody(statusCode int) *HTTPClient {
	m := HTTPClient{}
	m.On("Do", mock.Anything).Return(&http.Response{StatusCode: statusCode, Body: &faultyCloser{Reader: bytes.NewBufferString("{}")}}, nil).Once()
	return &m
}

type faultyCloser struct {
	io.Reader
}

func (c *faultyCloser) Close() error {
	return errors.New("close error")
}

func fixError() error {
	return errors.New("some error")
}
