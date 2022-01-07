package repos

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"isbnbook/app/log"
)

type Client interface {
	Get(url string, params map[string]string) ([]byte, error)
}

type httpClient struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
	logger     log.AppLogger
}

func NewClient(baseUrl string) (Client, error) {
	baseURL, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	log := log.GetLogger()
	return &httpClient{
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
		logger:     log,
	}, nil
}

func (c *httpClient) Get(url string, params map[string]string) ([]byte, error) {
	copiedURL := *c.BaseURL
	copiedURL.Path = path.Join(copiedURL.Path, url)
	q := copiedURL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	copiedURL.RawQuery = q.Encode()
	c.logger.Info(fmt.Sprintf("call url:%s", copiedURL.String()))
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
