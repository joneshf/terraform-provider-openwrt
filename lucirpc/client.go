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
	humanReadableCommitChanges = "commit changes"
	humanReadableCreateSection = "create section"
	humanReadableDeleteSection = "delete section"
	humanReadableGetSection    = "get section"
	humanReadableLogin         = "login"
	humanReadableShowChanges   = "show changes"
	humanReadableUpdateSection = "update section"

	methodChanges = "changes"
	methodCommit  = "commit"
	methodDelete  = "delete"
	methodGetAll  = "get_all"
	methodLogin   = "login"
	methodSection = "section"
	methodTSet    = "tset"

	pathAuth = "/cgi-bin/luci/rpc/auth"
	pathUCI  = "/cgi-bin/luci/rpc/uci"

	queryKeyAuth = "auth"
)

type Client struct {
	jsonRPCClientUCI jsonRPCClient
}

func (c *Client) CommitChanges(
	ctx context.Context,
	config string,
) (bool, error) {
	marshalledConfig, err := json.Marshal(config)
	if err != nil {
		return false, fmt.Errorf("unable to serialize config %q for %s: %w", config, humanReadableCommitChanges, err)
	}

	requestBody := jsonRPCRequestBody{
		Method: methodCommit,
		Params: []json.RawMessage{
			marshalledConfig,
		},
	}
	responseBody, err := c.jsonRPCClientUCI.Invoke(
		ctx,
		humanReadableCommitChanges,
		requestBody,
	)
	if err != nil {
		return false, fmt.Errorf("unable to %s: %w", humanReadableCommitChanges, err)
	}

	// The result can be `true` to indicate success,
	// or `null` to indicate failure.
	var result bool
	if responseBody == nil {
		return false, nil
	}

	err = json.Unmarshal(*responseBody, &result)
	if err != nil {
		return false, fmt.Errorf("unable to parse %s response: %w", humanReadableCommitChanges, err)
	}

	return result, nil
}

func (c *Client) CreateSection(
	ctx context.Context,
	config string,
	sectionType string,
	section string,
	options Options,
) (bool, error) {
	marshalledConfig, err := json.Marshal(config)
	if err != nil {
		return false, fmt.Errorf("unable to serialize config %q for %s: %w", config, humanReadableCreateSection, err)
	}

	marshalledSectionType, err := json.Marshal(sectionType)
	if err != nil {
		return false, fmt.Errorf("unable to serialize sectionType %q for %s: %w", sectionType, humanReadableCreateSection, err)
	}

	marshalledSection, err := json.Marshal(section)
	if err != nil {
		return false, fmt.Errorf("unable to serialize section %q for %s: %w", section, humanReadableCreateSection, err)
	}

	marshalledOptions, err := json.Marshal(options)
	if err != nil {
		return false, fmt.Errorf("unable to serialize options %q for %s: %w", options, humanReadableCreateSection, err)
	}

	requestBody := jsonRPCRequestBody{
		Method: methodSection,
		Params: []json.RawMessage{
			marshalledConfig,
			marshalledSectionType,
			marshalledSection,
			marshalledOptions,
		},
	}
	responseBody, err := c.jsonRPCClientUCI.Invoke(
		ctx,
		humanReadableCreateSection,
		requestBody,
	)
	if err != nil {
		return false, fmt.Errorf("unable to %s: %w", humanReadableCreateSection, err)
	}

	// The result can be `true` to indicate success,
	// or `null` to indicate failure.
	var result bool
	if responseBody == nil {
		return false, nil
	}

	err = json.Unmarshal(*responseBody, &result)
	if err != nil {
		return false, fmt.Errorf("unable to parse %s response: %w", humanReadableCreateSection, err)
	}

	if !result {
		return false, fmt.Errorf("unable to %s: it is not clear why this happened", humanReadableCreateSection)
	}

	result, err = c.CommitChanges(
		ctx,
		config,
	)
	if err != nil {
		return false, fmt.Errorf("was able to %s, but could not %s: %w", humanReadableCreateSection, humanReadableCommitChanges, err)
	}

	return result, nil
}

func (c *Client) DeleteSection(
	ctx context.Context,
	config string,
	section string,
) (bool, error) {
	marshalledConfig, err := json.Marshal(config)
	if err != nil {
		return false, fmt.Errorf("unable to serialize config %q for %s: %w", config, humanReadableDeleteSection, err)
	}

	marshalledSection, err := json.Marshal(section)
	if err != nil {
		return false, fmt.Errorf("unable to serialize section %q for %s: %w", section, humanReadableDeleteSection, err)
	}

	requestBody := jsonRPCRequestBody{
		Method: methodDelete,
		Params: []json.RawMessage{
			marshalledConfig,
			marshalledSection,
		},
	}
	responseBody, err := c.jsonRPCClientUCI.Invoke(
		ctx,
		humanReadableDeleteSection,
		requestBody,
	)
	if err != nil {
		return false, fmt.Errorf("unable to %s: %w", humanReadableDeleteSection, err)
	}

	// The result can be `true` to indicate success,
	// or `null` to indicate failure.
	var result bool
	if responseBody == nil {
		return false, nil
	}

	err = json.Unmarshal(*responseBody, &result)
	if err != nil {
		return false, fmt.Errorf("unable to parse %s response: %w", humanReadableDeleteSection, err)
	}

	if !result {
		return false, fmt.Errorf("unable to %s: it is not clear why this happened", humanReadableDeleteSection)
	}

	result, err = c.CommitChanges(
		ctx,
		config,
	)
	if err != nil {
		return false, fmt.Errorf("was able to %s, but could not %s: %w", humanReadableDeleteSection, humanReadableCommitChanges, err)
	}

	return result, nil
}

func (c *Client) GetSection(
	ctx context.Context,
	config string,
	section string,
) (Options, error) {
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
	responseBody, err := c.jsonRPCClientUCI.Invoke(
		ctx,
		humanReadableGetSection,
		requestBody,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to %s: %w", humanReadableGetSection, err)
	}

	if responseBody == nil {
		return nil, fmt.Errorf("could not find section %s.%s", config, section)
	}

	// Depending on the `config` and `section`,
	// this method can return a response that is an array instead of an object.
	// We have to handle that case as well.
	var unknownResult any
	err = json.Unmarshal(*responseBody, &unknownResult)
	if err != nil {
		return nil, fmt.Errorf("unable to determine type of %s response: %w", humanReadableGetSection, err)
	}

	_, ok := unknownResult.([]any)
	if ok {
		return nil, fmt.Errorf("incorrect config (%q) and/or section (%q): result from LuCI: %s", config, section, *responseBody)
	}

	var result Options
	err = json.Unmarshal(*responseBody, &result)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %s response: %w", humanReadableGetSection, err)
	}

	return result, nil
}

func (c *Client) ShowChanges(
	ctx context.Context,
	config string,
) ([][]string, error) {
	marshalledConfig, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize config %q for %s: %w", config, humanReadableShowChanges, err)
	}

	requestBody := jsonRPCRequestBody{
		Method: methodChanges,
		Params: []json.RawMessage{
			marshalledConfig,
		},
	}
	responseBody, err := c.jsonRPCClientUCI.Invoke(
		ctx,
		humanReadableShowChanges,
		requestBody,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to %s: %w", humanReadableShowChanges, err)
	}

	result := [][]string{}
	if responseBody == nil {
		return result, nil
	}

	err = json.Unmarshal(*responseBody, &result)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %s response: %w", humanReadableShowChanges, err)
	}

	return result, nil
}

func (c *Client) UpdateSection(
	ctx context.Context,
	config string,
	section string,
	options Options,
) (bool, error) {
	marshalledConfig, err := json.Marshal(config)
	if err != nil {
		return false, fmt.Errorf("unable to serialize config %q for %s: %w", config, humanReadableUpdateSection, err)
	}

	marshalledSection, err := json.Marshal(section)
	if err != nil {
		return false, fmt.Errorf("unable to serialize section %q for %s: %w", section, humanReadableUpdateSection, err)
	}

	marshalledOptions, err := json.Marshal(options)
	if err != nil {
		return false, fmt.Errorf("unable to serialize options %q for %s: %w", options, humanReadableCreateSection, err)
	}

	requestBody := jsonRPCRequestBody{
		Method: methodTSet,
		Params: []json.RawMessage{
			marshalledConfig,
			marshalledSection,
			marshalledOptions,
		},
	}
	responseBody, err := c.jsonRPCClientUCI.Invoke(
		ctx,
		humanReadableUpdateSection,
		requestBody,
	)
	if err != nil {
		return false, fmt.Errorf("unable to %s: %w", humanReadableUpdateSection, err)
	}

	// The result can be `true` to indicate success,
	// or `null` to indicate failure.
	var result bool
	if responseBody == nil {
		return false, nil
	}

	err = json.Unmarshal(*responseBody, &result)
	if err != nil {
		return false, fmt.Errorf("unable to parse %s response: %w", humanReadableUpdateSection, err)
	}

	if !result {
		return false, fmt.Errorf("unable to %s: it is not clear why this happened", humanReadableUpdateSection)
	}

	result, err = c.CommitChanges(
		ctx,
		config,
	)
	if err != nil {
		return false, fmt.Errorf("was able to %s, but could not %s: %w", humanReadableUpdateSection, humanReadableCommitChanges, err)
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

	defer response.Body.Close()
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
