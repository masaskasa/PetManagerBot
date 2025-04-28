package telegram

import (
	"encoding/json"
	"io"
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

func New(host string, token string) *Client {
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

	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	request.URL.RawQuery = query.Encode()

	response, err := client.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}
