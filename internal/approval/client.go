package approval

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	http    *http.Client
	baseURL string
}

func NewClient(httpClient *http.Client, baseURL string) *Client {
	return &Client{http: httpClient, baseURL: baseURL}
}

func (c *Client) Create(req CreateRequest) (*Response, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Post(c.baseURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("backend %d döndü", resp.StatusCode)
	}
	return decode(resp)
}

func (c *Client) waitOnce(id int64) (*Response, error) {
	url := fmt.Sprintf("%s/%d/wait", c.baseURL, id)
	resp, err := c.http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return decode(resp)
}

func (c *Client) WaitForResolution(id int64, maxRounds int) (*Response, error) {
	var last *Response
	for i := 0; i < maxRounds; i++ {
		resp, err := c.waitOnce(id)
		if err != nil {
			return nil, err
		}
		last = resp
		if resp.Status.IsTerminal() {
			return resp, nil
		}
	}
	return last, nil
}

func decode(resp *http.Response) (*Response, error) {
	var out Response
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
