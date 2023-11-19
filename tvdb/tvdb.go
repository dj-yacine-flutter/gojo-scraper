package tvdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	Client   *http.Client
	Language string
	Token    string
}

func NewClient(client *http.Client) (tvdb *Client) {
	tvdb = &Client{
		Language: "en",
		Token:    "",
		Client:   client,
	}
	return
}

type RequestArgs struct {
	Body   interface{}
	Method string
	Params url.Values
	Path   string
}

type Errors struct {
	InvalidFilters     []string `json:"invalidFilters"`
	InvalidLanguage    string   `json:"invalidLanguage"`
	InvalidQueryParams []string `json:"invalidQueryParams"`
}

type Links struct {
	First    int `json:"first"`
	Last     int `json:"last"`
	Next     int `json:"next"`
	Previous int `json:"previous"`
}

const (
	ApplicationJson = "application/json"
	BaseURL         = "https://api4.thetvdb.com/v4"
	DefaultLanguage = "en"
)

func (c *Client) BuildUrlPath(path string) (result string) {
	result = fmt.Sprintf("%s%s", BaseURL, path)
	return
}

func (c *Client) ParseResponse(body io.ReadCloser, data interface{}) error {
	return json.NewDecoder(body).Decode(data)
}

func (c *Client) DoRequest(args RequestArgs) (resp *http.Response, err error) {
	var body io.Reader
	if args.Body != nil {
		marshal, _ := json.Marshal(args.Body)
		body = bytes.NewBuffer(marshal)
	}
	req, err := http.NewRequest(args.Method, c.BuildUrlPath(args.Path), body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", ApplicationJson)
	req.Header.Set("Accept", ApplicationJson)
	req.Header.Set("Accept-Language", c.Language)
	if c.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}
	if args.Params != nil {
		req.URL.RawQuery = args.Params.Encode()
	}

	resp, err = c.Client.Do(req)
	if err == nil && resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("received error status code (%d): %s", resp.StatusCode, resp.Status)
		return
	}
	return
}

type AuthenticationResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
	Status string `json:"status"`
}

type AuthenticationRequest struct {
	ApiKey string `json:"apikey"`
	Pin    string `json:"pin"`
}

const AuthenticationLogin string = "/login"

func (c *Client) Login(requestParams *AuthenticationRequest) (err error) {
	resp, err := c.DoRequest(RequestArgs{
		Body:   requestParams,
		Path:   AuthenticationLogin,
		Method: http.MethodPost,
	})
	if err != nil {
		return
	}
	_, err = c.saveToken(resp)
	defer resp.Body.Close()
	return
}

func (c *Client) saveToken(resp *http.Response) (data *AuthenticationResponse, err error) {
	data = new(AuthenticationResponse)
	err = c.ParseResponse(resp.Body, data)
	if err != nil {
		return
	}
	c.Token = data.Data.Token
	return
}
