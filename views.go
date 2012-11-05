package lotus

import "encoding/json"
import "fmt"
import "errors"

type View struct {
	Url    string
	Id     int
	Title  string
	Active bool
}

type Column struct {
	Id    string
	Title string
}

type ViewResult struct {
	NextPageUrl     string `json:"next_page"`
	PreviousPageUrl string `json:"previous_page"`
	Count           int
	Columns         []Column
	Users           []ViewResultData
	Groups          []ViewResultData
	Organizations   []ViewResultData
	Rows            []*ViewResultRow
}

type ViewResultData struct {
	Id   int
	Name string
	Url  string
}

type ViewResultRow struct {
	Score          int
	Subject        string
	RequesterId    int `json:"requester_id"`
	AssigneeId     int `json:"assignee_id"`
	GroupId        int `json:"group_id"`
	OrganizationId int `json:"organization_id"`

	Created      string
	Type         string
	Priority     string
	Tags         []string
	Assignee     ViewResultData
	Requester    ViewResultData
	Group        ViewResultData
	Organization ViewResultData
	Ticket       ViewTicket
}

type ViewTicket struct {
	Url         string
	Id          int
	Subject     string
	Description string
	Status      string
	Type        string
	Priority    string
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
	var url = fmt.Sprintf("/views/%d/execute.json", v.Id)
	return c.getViewResults(url)
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
	finalizeViewResult(data)
	return data, err
}

func makeCache(lst []ViewResultData) map[int]ViewResultData {
	var cache = make(map[int]ViewResultData)
	if len(lst) != 0 {
		for _, item := range lst {
			cache[item.Id] = item
		}
	}
	return cache
}

func finalizeViewResult(result *ViewResult) {
	var (
		userMap  = makeCache(result.Users)
		groupMap = makeCache(result.Groups)
		orgMap   = makeCache(result.Organizations)
	)

	for _, row := range result.Rows {

		if row.RequesterId != 0 {
			row.Requester = userMap[row.RequesterId]
		}

		if row.AssigneeId != 0 {
			row.Assignee = userMap[row.AssigneeId]
		}

		if row.GroupId != 0 {
			row.Group = groupMap[row.GroupId]
		}

		if row.OrganizationId != 0 {
			row.Organization = orgMap[row.OrganizationId]
		}
	}
}
