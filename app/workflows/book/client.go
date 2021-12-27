package book

import (
	"isbnbook/app/log"

	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type client struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
}

func NewClient(baseUrl string) (*client, error) {
	baseURL, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	return &client{
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
	}, nil
}

var logger = log.GetLogger()

func (c *client) Get(url string, params map[string]string) ([]byte, error) {
	copiedURL := *c.BaseURL
	copiedURL.Path = path.Join(copiedURL.Path, url)
	q := copiedURL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	copiedURL.RawQuery = q.Encode()
	logger.Info(fmt.Sprintf("call url:%s", copiedURL.String()))
	resp, err := c.HTTPClient.Get(copiedURL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code is not OK(%d)", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}
