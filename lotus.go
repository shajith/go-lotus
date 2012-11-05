package lotus

import "fmt"
import "net/http"
import "io/ioutil"

type Client struct {
	Subdomain string
	Username  string
	Password  string
	debug     bool
}

func (c *Client) perform(method string, path string) (body []byte, err error) {
	if c.debug {
		fmt.Println("Performing: ", path)
	}
	var httpClient http.Client
	var url = fmt.Sprintf("https://%s.zendesk.com/api/v2%s", c.Subdomain, path)
	req, err := http.NewRequest(method, url, nil)
	req.SetBasicAuth(c.Username, c.Password)
	resp, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	return
}

func (c *Client) Debug(on bool) {
	c.debug = on
}

func New(subdomain string, user string, password string) *Client {
	return &Client{Subdomain: subdomain, Username: user, Password: password}
}
