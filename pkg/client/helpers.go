package client

import (
	"fmt"
	"nutritionapp/pkg/server"
)

// Make a request to the server, and return the response, without type checking
func makeRequest(c *Client, reqType string, data any) (any, error) {
	resp := make(chan server.Response)
	c.requests <- server.Request{
		Type:   reqType,
		Data:   data,
		Return: resp,
	}

	a := <-resp
	if a.Error != nil {
		return server.Response{}, a.Error
	}
	return a.Data, nil
}

// Make a request to the server, and return the response casted to the specified type
func makeRequestTyped[T any](c *Client, reqType string, data any) (*T, error) {
	data, err := makeRequest(c, reqType, data)
	if err != nil {
		return nil, err
	}

	result, ok := data.(T)
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", data)
	}

	return &result, nil
}
