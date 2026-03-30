package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	BaseURL     = "https://api.line.me"
	BaseDataURL = "https://api-data.line.me"
	BaseOAURL   = "https://developers-oaplus.line.biz"
)

type Client struct {
	AccessToken string
	HTTPClient  *http.Client
}

func New(token string) *Client {
	return &Client{
		AccessToken: token,
		HTTPClient:  &http.Client{},
	}
}

func (c *Client) do(method, rawURL string, body any, headers map[string]string) ([]byte, int, error) {
	var bodyReader io.Reader

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, rawURL, bodyReader)
	if err != nil {
		return nil, 0, err
	}

	if c.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	if resp.StatusCode >= 400 {
		return nil, resp.StatusCode, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(data))
	}

	return data, resp.StatusCode, nil
}

func (c *Client) Get(path string, params map[string]string) ([]byte, error) {
	u := BaseURL + path
	if len(params) > 0 {
		q := url.Values{}
		for k, v := range params {
			q.Set(k, v)
		}
		u += "?" + q.Encode()
	}
	data, _, err := c.do(http.MethodGet, u, nil, nil)
	return data, err
}

func (c *Client) Post(path string, body any) ([]byte, error) {
	data, _, err := c.do(http.MethodPost, BaseURL+path, body, nil)
	return data, err
}

func (c *Client) Put(path string, body any) ([]byte, error) {
	data, _, err := c.do(http.MethodPut, BaseURL+path, body, nil)
	return data, err
}

func (c *Client) Delete(path string) ([]byte, error) {
	data, _, err := c.do(http.MethodDelete, BaseURL+path, nil, nil)
	return data, err
}

func (c *Client) PostForm(rawURL string, fields map[string]string) ([]byte, error) {
	form := url.Values{}
	for k, v := range fields {
		form.Set(k, v)
	}
	req, err := http.NewRequest(http.MethodPost, rawURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(data))
	}
	return data, nil
}

func PrintJSON(data []byte) {
	var buf bytes.Buffer
	if err := json.Indent(&buf, data, "", "  "); err != nil {
		fmt.Println(string(data))
		return
	}
	fmt.Println(buf.String())
}
