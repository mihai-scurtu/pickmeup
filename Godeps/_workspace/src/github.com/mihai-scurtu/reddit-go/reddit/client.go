package reddit

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"net/url"
	"time"
	"fmt"
	"errors"
)

const (
	URL = "https://api.reddit.com"
	OAUTH_URL = "https://oauth.reddit.com"
)

var httpClient = &http.Client{}

type Client struct {
	UserAgent string

	Username string
	Password string
	ClientId string
	ClientSecret string

	Token string
	TokenExpires int64
}

type tokenResponse struct {
	Token string `json:"access_token"`
	Expires int64 `json:"expires_in"`
}

func NewClient(userAgent string) *Client {

	return &Client{UserAgent: userAgent}
}

func (c *Client) createRequest(method string, uri string, body io.Reader) *http.Request {
	url := c.url(uri)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", c.UserAgent)

	if c.Token != "" {
		req.Header.Add("Authorization", "Bearer " + c.Token)
	}

	return req
}

func (c *Client) Get(uri string) []byte {
	req := c.createRequest("GET", uri, nil)

	resp, err := c.run(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body
}

func (c *Client) GetFrontPage() *PostListing {
	return c.GetPostListing("/")
}

func (this *Client) GetPostListing(uri string) *PostListing {
	var l *PostListing

	resp := this.Get(uri)

	//log.Println(string(resp))

	if err := json.Unmarshal(resp, &l); err !=nil {
		log.Fatal(err)
	}

	return l
}

func (c *Client) GetNewPosts() *PostListing {
	return c.GetPostListing("/new")
}

func (c *Client) GetComments(subreddits ...string) *CommentListing {
	var l *CommentListing

	uri := "/r/" + strings.Join(subreddits, "+") + "/comments"


	resp := c.Get(uri)
	fmt.Println(string(resp))

	if err := json.Unmarshal(resp, &l); err != nil {
		log.Fatal(err)
	}

	return l
}

func (this *Client) url(uri string) string {
	if len(uri) == 0 || uri[0] != '/' {
		uri = "/" + uri
	}

	base := URL
	if this.Token != "" {
		base = OAUTH_URL
	}
 	return base + uri
}

func (this *Client) run(req *http.Request) (*http.Response, error) {
	if this.Token != "" && this.TokenExpires < time.Now().Unix() {
		this.GetToken()
	}

	return httpClient.Do(req)
}

func (this *Client) GetToken() error {
	form := url.Values{}

	form.Add("grant_type", "password")
	form.Add("username", this.Username)
	form.Add("password", this.Password)

	req, err := http.NewRequest("POST", "https://www.reddit.com/api/v1/access_token", strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("User-Agent", this.UserAgent)
	req.SetBasicAuth(this.ClientId, this.ClientSecret)

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	tr := tokenResponse{}
	err = json.Unmarshal(body, &tr)
	if err != nil {
		log.Fatal(err)
	}

	if tr.Token == "" {
		return errors.New(string(body))
	}

	this.Token = tr.Token
	this.TokenExpires = time.Now().Unix() + tr.Expires - 600

	return nil
}
