package microsoft

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// Client Outlook API Client
type Client struct {
	baseURL string
	version string
	ctx     context.Context
	client  *http.Client
}

// NewClient creates a new Microsoft Graph API client.
// The client uses version 1.0 by default.
func NewClient(ctx context.Context, client *http.Client) *Client {
	c := &Client{
		baseURL: "https://graph.microsoft.com",
		version: "v1.0",
		ctx:     ctx,
		client:  client,
	}
	return c
}

// Get perform http GET
func (c *Client) Get(path string, target interface{}) (resp *http.Response, err error) {
	url := c.buildURL(path)
	resp, err = sendRequest(c.ctx, c.client, "GET", url, nil)
	if err != nil {
		return resp, err
	}
	err = applyResponseBody(resp, target)
	return resp, err
}

// Post perform http POST
func (c *Client) Post(path string, body interface{}) (resp *http.Response, err error) {
	url := c.buildURL(path)
	return sendRequest(c.ctx, c.client, "POST", url, body)
}

// Patch perform http Patch
func (c *Client) Patch(path string, body interface{}) (resp *http.Response, err error) {
	url := c.buildURL(path)
	return sendRequest(c.ctx, c.client, "PATCH", url, body)
}

// Delete perform http DELETE
func (c *Client) Delete(path string) (resp *http.Response, err error) {
	url := c.buildURL(path)
	return sendRequest(c.ctx, c.client, "DELETE", url, nil)
}

// Prepends the path with the base URL unless a fully-qualified URL is given.
func (c *Client) buildURL(path string) string {
	url, _ := url.Parse(path)
	if url.Scheme != "" {
		return path
	}

	return fmt.Sprintf("%s/%s/%s", c.baseURL, c.version, strings.TrimPrefix(path, "/"))
}

func sendRequest(ctx context.Context, client *http.Client, method, path string, body interface{}) (resp *http.Response, err error) {
	// Verify and encode the URL.
	u, err := url.Parse(path)
	if err != nil {
		return nil, errors.Wrap(err, "invalid url")
	}
	path = u.String()

	// Declare `b` so that its nil value can be passed to http.NewRequest()
	// when `body` is nil.
	var b []byte

	if body != nil {
		b, err = json.Marshal(body)
		if err != nil {
			return nil, errors.Wrap(err, "invalid body")
		}
	}

	req, err := http.NewRequest(method, path, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return client.Do(req.WithContext(ctx))
}

func applyResponseBody(resp *http.Response, target interface{}) error {
	// Don't bother capturing the response body if there isn't a target.
	if target == nil {
		return nil
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	switch resp.ContentLength {
	case 0:
		return nil
	case -1:
		// Can't safely use the JSON decoder because the body might be empty.
		body, err := ReadHTTPResponse(resp)
		if err != nil {
			return err
		}

		// An error isn't returned if there is nothing to unmarshal.
		// It may be that a header or parameter was passed which caused the
		// server to respond without a body. The caller can check the response
		// to discover why the target wasn't populated if a body was expected.
		if len(body) > 0 {
			return json.Unmarshal(body, target)
		}
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

// ReadHTTPResponse returns the response body.
func ReadHTTPResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
