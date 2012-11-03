package lotus

import "fmt"
import "net/http"
import "io/ioutil"
import "encoding/json"
import "errors"

func Hello(name string) string {
	return fmt.Sprintf("Hello %s", name);
}

func fetchBytes(url string) []byte {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
  return body
}

func Fetch(url string) string {
	return fmt.Sprintf("%s", fetchBytes(url))
}

type Client struct {
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
	var url = fmt.Sprintf("https://shajith.zendesk.com/api/v2/views/%d/tickets.json", v.Id);
	return c.getViewResults(url);
}

func (c *Client) getViewResults(url string) (data *ViewResult, err error) {
	body, err := perform("GET", url, c.User, c.Password)
	json.Unmarshal(body, &data)
	return data, err;
}

func (c *Client) Me() User {
	var data = make(map[string]User)
	var body, _ = perform("GET", "https://shajith.zendesk.com/api/v2/users/me.json", c.User, c.Password)
	json.Unmarshal(body, &data)
  return data["user"]
}

func (c *Client) Views() []View {
	var data = make(map[string][]View)
	var body, _ = perform("GET", "https://shajith.zendesk.com/api/v2/views.json", c.User, c.Password)
	json.Unmarshal(body, &data)
  return data["views"]
}

func perform(method string, url string, user string, password string) (body []byte, err error){
	var httpClient http.Client

	req, err := http.NewRequest(method, url, nil)
	req.SetBasicAuth(user, password)
	resp, err := httpClient.Do(req)
	if (err != nil) {
		body = nil
  } else {
		body, err = ioutil.ReadAll(resp.Body)
	}

	return
}
