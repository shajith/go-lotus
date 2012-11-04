package lotus

import "fmt"
import "net/http"
import "io/ioutil"
import "encoding/json"
import "errors"

type Client struct {
	Subdomain string
	User      string
	Password  string
}

type User struct {
	Url               string
	Id                int
	Name              string
	Email             string
	Role              string
	AuthenticityToken string `json:"authenticity_token"`
	Locale            string
	Tags              []string
}

type View struct {
	Url    string
	Id     int
	Title  string
	Active bool
}

type Ticket struct {
	Url         string
	Id          int
	Subject     string
	Description string
	Tags        []string
	AssigneeId  int `json:"assignee_id"`
	RequesterId int `json:"requester_id"`
	requester   *User
	assignee    *User
}

func (c *Client) Assignee(t *Ticket) (user *User, err error) {
	if t.assignee != nil {
		user = t.assignee
	} else {
		if t.AssigneeId != 0 {
			user, err = c.getUser(fmt.Sprintf("%d", t.AssigneeId))
			if err == nil {
				t.assignee = user
			}
		} else {
			user, err = nil, nil
		}
	}

	return
}

type ViewResult struct {
	NextPageUrl     string `json:"next_page"`
	PreviousPageUrl string `json:"previous_page"`
	Count           int
	Tickets         []Ticket
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
	var url = fmt.Sprintf("/views/%d/tickets.json", v.Id)
	return c.getViewResults(url)
}

func (c *Client) Me() (*User, error) {
	return c.getUser("me")
}

func (c *Client) Views() []View {
	var data = make(map[string][]View)
	var body, _ = c.perform("GET", "/views.json")
	json.Unmarshal(body, &data)
	return data["views"]
}

func (c *Client) getViewResults(url string) (data *ViewResult, err error) {
	body, err := c.perform("GET", url)
	json.Unmarshal(body, &data)
	return data, err
}

func (c *Client) getUser(id string) (user *User, err error) {
	var data = make(map[string]User)
	body, err := c.perform("GET", fmt.Sprintf("/users/%s.json", id))

	if err == nil {
		json.Unmarshal(body, &data)
		u := data["user"]
		user = &u
	}

	return
}

func (c *Client) perform(method string, path string) (body []byte, err error) {

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
