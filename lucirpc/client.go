package lucirpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	methodGetAll = "get_all"
	methodLogin  = "login"

	pathAuth = "/cgi-bin/luci/rpc/auth"
	pathUCI  = "/cgi-bin/luci/rpc/uci"

	queryKeyAuth = "auth"
)

type Client struct {
	addressUCI url.URL
	client     *http.Client
}

func (c *Client) GetSection(
	ctx context.Context,
	config string,
	section string,
) (map[string]string, error) {
	requestBody := getSectionRequestBody{
		Method: methodGetAll,
		Params: [2]string{config, section},
	}
	buffer := bytes.Buffer{}
	encoder := json.NewEncoder(&buffer)
	err := encoder.Encode(requestBody)
	if err != nil {
		return nil, fmt.Errorf("problem encoding get section request: %w", err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.addressUCI.String(),
		&buffer,
	)
	if err != nil {
		return nil, fmt.Errorf("problem creating get section request: %w", err)
	}

	response, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("problem sending request to get section: %w", err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("expected 200 response: got %s", response.Status)
	}

	var responseBody getSectionResponseBody
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&responseBody)
	if err != nil {
		return nil, fmt.Errorf("unable to process get section response: %w", err)
	}

	if responseBody.Error != nil {
		return nil, fmt.Errorf("unable to get section: %s", *responseBody.Error)
	}

	if responseBody.Result == nil {
		return nil, errors.New("invalid get section response: expected either an error or a result, got neither")
	}

	result := *responseBody.Result
	return result, nil
}

func NewClient(
	ctx context.Context,
	scheme string,
	hostname string,
	port uint16,
	username string,
	password string,
) (*Client, error) {
	host := hostname
	if port != 0 {
		host = fmt.Sprintf("%s:%d", host, port)
	}

	address := url.URL{
		Host:   host,
		Path:   pathAuth,
		Scheme: scheme,
	}
	httpClient := &http.Client{}
	requestBody := authRequestBody{
		Method: methodLogin,
		Params: [2]string{username, password},
	}
	buffer := bytes.Buffer{}
	encoder := json.NewEncoder(&buffer)
	err := encoder.Encode(requestBody)
	if err != nil {
		return nil, fmt.Errorf("problem encoding login request: %w", err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		address.String(),
		&buffer,
	)
	if err != nil {
		return nil, fmt.Errorf("problem creating login request: %w", err)
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("problem sending request to login: %w", err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("expected 200 response: got %s", response.Status)
	}

	var responseBody authResponseBody
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&responseBody)
	if err != nil {
		return nil, fmt.Errorf("unable to process login response: %w", err)
	}

	if responseBody.Error != nil {
		return nil, fmt.Errorf("unable to login: %s", *responseBody.Error)
	}

	if responseBody.Result == nil {
		return nil, errors.New("invalid login response: expected either an error or a result, got neither")
	}

	authToken := *responseBody.Result
	query := url.Values{}
	query.Add(queryKeyAuth, authToken)
	addressUCI := url.URL{
		Host:     host,
		Path:     pathUCI,
		RawQuery: query.Encode(),
		Scheme:   scheme,
	}
	client := &Client{
		addressUCI: addressUCI,
		client:     httpClient,
	}
	return client, nil
}

type authRequestBody struct {
	Method string    `json:"method"`
	Params [2]string `json:"params"`
}

type authResponseBody struct {
	Error  *string `json:"error"`
	Result *string `json:"result"`
}

type getSectionRequestBody struct {
	Method string    `json:"method"`
	Params [2]string `json:"params"`
}

type getSectionResponseBody struct {
	Error  *string            `json:"error"`
	Result *map[string]string `json:"result"`
}
