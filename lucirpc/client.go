package lucirpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	humanReadableGetSection = "get section"
	humanReadableLogin      = "login"

	methodGetAll = "get_all"
	methodLogin  = "login"

	pathAuth = "/cgi-bin/luci/rpc/auth"
	pathUCI  = "/cgi-bin/luci/rpc/uci"

	queryKeyAuth = "auth"
)

type Client struct {
	jsonRPCClientUCI jsonRPCClient
}

func (c *Client) GetSection(
	ctx context.Context,
	config string,
	section string,
) (map[string]string, error) {
	requestBody := jsonRPCRequestBody{
		Method: methodGetAll,
		Params: []string{config, section},
	}
	responseBody, err := c.jsonRPCClientUCI.Invoke(
		ctx,
		humanReadableGetSection,
		requestBody,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to %s: %w", humanReadableGetSection, err)
	}

	var result map[string]string
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %s response: %w", humanReadableGetSection, err)
	}

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
	requestBody := jsonRPCRequestBody{
		Method: methodLogin,
		Params: []string{username, password},
	}
	jsonRPCClient := jsonRPCNewClient(
		*httpClient,
		address,
	)
	responseBody, err := jsonRPCClient.Invoke(
		ctx,
		humanReadableLogin,
		requestBody,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to %s: %w", humanReadableLogin, err)
	}

	var authToken string
	err = json.Unmarshal(responseBody, &authToken)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %s response: %w", humanReadableLogin, err)
	}

	query := url.Values{}
	query.Add(queryKeyAuth, authToken)
	addressUCI := url.URL{
		Host:     host,
		Path:     pathUCI,
		RawQuery: query.Encode(),
		Scheme:   scheme,
	}
	jsonRPCClientUCI := jsonRPCNewClient(
		*httpClient,
		addressUCI,
	)
	client := &Client{
		jsonRPCClientUCI: jsonRPCClientUCI,
	}
	return client, nil
}

type jsonRPCClient struct {
	address url.URL
	client  http.Client
}

func (c jsonRPCClient) Invoke(
	ctx context.Context,
	humanReadableMethod string,
	requestBody jsonRPCRequestBody,
) (json.RawMessage, error) {
	buffer := bytes.Buffer{}
	encoder := json.NewEncoder(&buffer)
	err := encoder.Encode(requestBody)
	if err != nil {
		return nil, fmt.Errorf("problem encoding %s request: %w", humanReadableMethod, err)
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.address.String(),
		&buffer,
	)
	if err != nil {
		return nil, fmt.Errorf("problem creating %s request: %w", humanReadableMethod, err)
	}

	response, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("problem sending request to %s: %w", humanReadableMethod, err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("expected %s to respond with a 200: got %s", humanReadableMethod, response.Status)
	}

	var responseBody jsonRPCResponseBody
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&responseBody)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %s response: %w", humanReadableMethod, err)
	}

	if responseBody.Error == nil && responseBody.Result == nil {
		return nil, fmt.Errorf("invalid %s response: expected either an error or a result, got neither", humanReadableMethod)
	}

	if responseBody.Error != nil {
		return nil, fmt.Errorf("%s error: %s", humanReadableMethod, *responseBody.Error)
	}

	return *responseBody.Result, nil
}

func jsonRPCNewClient(
	httpClient http.Client,
	address url.URL,
) jsonRPCClient {
	return jsonRPCClient{
		address: address,
		client:  httpClient,
	}
}

type jsonRPCRequestBody struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
}

type jsonRPCResponseBody struct {
	Error  *string          `json:"error"`
	Result *json.RawMessage `json:"result"`
}
