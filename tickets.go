package lotus

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
			user, err = c.User(t.AssigneeId)
			if err == nil {
				t.assignee = user
			}
		} else {
			user, err = nil, nil
		}
	}

	return
}
