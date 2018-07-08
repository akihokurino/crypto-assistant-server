package api

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

type HttpClient interface {
	Get(ctx context.Context, urlString string, params map[string]string) (map[string]interface{}, error)
}

type httpClient struct {
	client *http.Client
}

func NewHttpClient() HttpClient {
	return &httpClient{}
}

func (c *httpClient) Get(ctx context.Context, urlString string, params map[string]string) (map[string]interface{}, error) {
	client := urlfetch.Client(ctx)

	req, err := http.NewRequest("GET", urlString, nil)

	if err != nil {
		return nil, err
	}

	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}
	req.URL.RawQuery = values.Encode()

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var f interface{}
	json.Unmarshal(body, &f)
	m := f.(map[string]interface{})

	return m, nil
}
