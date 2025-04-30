package telegram

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func NewClient(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (client *Client) Updates(offset int, limit int) ([]update, error) {
	query := url.Values{}
	query.Add("offset", strconv.Itoa(offset))
	query.Add("limit", strconv.Itoa(limit))

	data, err := client.getRequest(getUpdates, query)
	if err != nil {
		return nil, err
	}

	var updates receivedUpdates

	if err := json.Unmarshal(data, &updates); err != nil {
		slog.Error("getRequest: error of parse response data:", err.Error())
		return nil, err
	}

	return updates.Updates, nil
}

func (client *Client) getRequest(method string, query url.Values) ([]byte, error) {

	url := url.URL{
		Scheme: "https",
		Host:   client.host,
		Path:   path.Join(client.basePath, method),
	}

	slog.Info("getRequest: done url for GET request with query parameters:", url.String(), query.Encode())

	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		slog.Error("getRequest: error of making GET request:", err.Error())
		return nil, err
	}

	request.URL.RawQuery = query.Encode()

	response, err := client.client.Do(request)
	if err != nil {
		slog.Error("getRequest: error of send GET request:", err.Error())
		return nil, err
	}

	slog.Info("getRequest: response status and header:", response.Status, response.Header)

	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		slog.Error("getRequest: error of read response body: %s", err.Error())
		return nil, err
	}

	slog.Info("getRequest: response body:", body)

	return body, err
}

func (client *Client) postRequest(method string, data []byte) ([]byte, error) {

	slog.Info("postRequest: get body for request:", data)

	url := url.URL{
		Scheme: "https",
		Host:   client.host,
		Path:   path.Join(client.basePath, method),
	}

	slog.Info("postRequest: done url for POST request:", url.String())

	requestBody := bytes.NewBuffer(data)

	request, err := http.NewRequest(http.MethodPost, url.String(), requestBody)
	if err != nil {
		slog.Error("postRequest: error of making POST request:", err.Error())
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := client.client.Do(request)
	if err != nil {
		slog.Error("postRequest: error of send POST request:", err.Error())
		return nil, err
	}

	slog.Info("ostRequest: response status and header:", response.Status, response.Header)

	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		slog.Error("postRequest: error of read response body: %s", err.Error())
		return nil, err
	}

	slog.Info("postRequest: response body:", body)

	return body, err
}
