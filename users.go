package lotus

import "fmt"
import "encoding/json"

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

func (c *Client) Me() (*User, error) {
	return c.getUser("me")
}

func (c *Client) User(id int) (*User, error) {
	return c.getUser(fmt.Sprintf("%s", id))
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
