package lotus

import "fmt"
import "net/http"
import "io/ioutil"
import "encoding/json"
import "errors"

type Client struct {
	Subdomain string
  User string
	Password string
}

type User struct {
	Url string
	Id int
  Name string
	Email string
	Role string
	AuthenticityToken string `json:"authenticity_token"`
	Locale string
	Tags []string
}

type View struct {
	Url string
	Id int
	Title string
	Active bool
}

type Ticket struct {
	Url string
	Id int
	Subject string
	Description string
	Tags []string
}

type ViewResult struct {
  NextPageUrl string `json:"next_page"`
  PreviousPageUrl string `json:"previous_page"`
	Count int
  Tickets []Ticket
}

func (c *Client) ExecuteViewNextPage(r *ViewResult) (result *ViewResult, err error) {
	if r.NextPageUrl != "" {
		result, err = c.getViewResults(r.NextPageUrl)
	} else {
		result, err = nil, errors.New("no next page")
	}

	return
}

func (c *Client) ExecuteView(v View) (*ViewResult, error) {
	var url = fmt.Sprintf("/views/%d/tickets.json", v.Id);
	return c.getViewResults(url);
}

func (c *Client) getViewResults(path string) (data *ViewResult, err error) {
	body, err := c.perform("GET", path)
	json.Unmarshal(body, &data)
	return data, err;
}

func (c *Client) Me() User {
	var data = make(map[string]User)
	var body, _ = c.perform("GET", "/users/me.json")
	json.Unmarshal(body, &data)
  return data["user"]
}

func (c *Client) Views() []View {
	var data = make(map[string][]View)
	var body, _ = c.perform("GET", "/views.json")
	json.Unmarshal(body, &data)
  return data["views"]
}

func (c *Client) perform(method string, path string) (body []byte, err error){
	var httpClient http.Client
	var url = fmt.Sprintf("https://%s.zendesk.com/api/v2%s", c.Subdomain, path)
	req, err := http.NewRequest(method, url, nil)
	req.SetBasicAuth(c.User, c.Password)
	resp, err := httpClient.Do(req)
	if err != nil {
		body = nil
	} else {
		body, err = ioutil.ReadAll(resp.Body)
	}
	return
}
