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
) (map[string]json.RawMessage, error) {
	marshalledConfig, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize config %q for %s: %w", config, humanReadableGetSection, err)
	}

	marshalledSection, err := json.Marshal(section)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize section %q for %s: %w", section, humanReadableGetSection, err)
	}

	requestBody := jsonRPCRequestBody{
		Method: methodGetAll,
		Params: []json.RawMessage{
			marshalledConfig,
			marshalledSection,
		},
	}
	responseBody, err := c.jsonRPCClientUCI.InvokeNotNull(
		ctx,
		humanReadableGetSection,
		requestBody,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to %s: %w", humanReadableGetSection, err)
	}

	// Depending on the `config` and `section`,
	// this method can return a response that is an array instead of an object.
	// We have to handle that case as well.
	var unknownResult any
	err = json.Unmarshal(responseBody, &unknownResult)
	if err != nil {
		return nil, fmt.Errorf("unable to determine type of %s response: %w", humanReadableGetSection, err)
	}

	_, ok := unknownResult.([]any)
	if ok {
		return nil, fmt.Errorf("incorrect config (%q) and/or section (%q): result from LuCI: %s", config, section, responseBody)
	}

	var result map[string]json.RawMessage
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
	marshalledUsername, err := json.Marshal(username)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize username for %s: %w", humanReadableLogin, err)
	}

	marshalledPassword, err := json.Marshal(password)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize password for %s: %w", humanReadableLogin, err)
	}

	requestBody := jsonRPCRequestBody{
		Method: methodLogin,
		Params: []json.RawMessage{
			marshalledUsername,
			marshalledPassword,
		},
	}
	jsonRPCClient := jsonRPCNewClient(
		*httpClient,
		address,
	)
	responseBody, err := jsonRPCClient.InvokeNotNull(
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

func (c jsonRPCClient) InvokeNotNull(
	ctx context.Context,
	humanReadableMethod string,
	requestBody jsonRPCRequestBody,
) (json.RawMessage, error) {
	result, err := c.Invoke(
		ctx,
		humanReadableMethod,
		requestBody,
	)
	if err != nil {
		return json.RawMessage{}, err
	}

	if result == nil {
		return nil, fmt.Errorf("invalid %s response: expected either an error or a result, got neither", humanReadableMethod)
	}

	return *result, nil
}

func (c jsonRPCClient) Invoke(
	ctx context.Context,
	humanReadableMethod string,
	requestBody jsonRPCRequestBody,
) (*json.RawMessage, error) {
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

	if responseBody.Error != nil {
		return nil, fmt.Errorf("%s error: %s", humanReadableMethod, *responseBody.Error)
	}

	return responseBody.Result, nil
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
	Method string            `json:"method"`
	Params []json.RawMessage `json:"params"`
}

type jsonRPCResponseBody struct {
	Error  *string          `json:"error"`
	Result *json.RawMessage `json:"result"`
}
