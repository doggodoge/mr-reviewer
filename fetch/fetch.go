package fetch

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"mr-reviewer/config"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

//go:embed  gql/fetchMergeRequests.txt
var fetchMergeRequests string

type MRsResponse struct {
	Data Data `json:"data"`
}

type Data struct {
	Project Project `json:"project"`
}

type Project struct {
	MergeRequests MergeRequests `json:"mergeRequests"`
}

type MergeRequests struct {
	Nodes []MR `json:"nodes"`
}

type MR struct {
	Title  string `json:"title"`
	Author struct {
		Name string `json:"name"`
	} `json:"author"`
	ApprovalsRequired int `json:"approvalsRequired"`
	HeadPipeline      struct {
		Status string `json:"status"`
	} `json:"headPipeline"`
	URL string `json:"webUrl"`
}

func (mrs *MRsResponse) ToListItems(showDraft bool) []list.Item {
	var l []list.Item

	if !showDraft {
		mrs = FilterDraft(mrs)
	}

	for _, mr := range mrs.Data.Project.MergeRequests.Nodes {
		var approved string
		if mr.ApprovalsRequired != 0 {
			approved = "No"
		} else {
			approved = "Yes"
		}
		desc := fmt.Sprintf(
			"%s | Approved: %s | Pipeline Status: %s",
			mr.Author.Name,
			approved,
			mr.HeadPipeline.Status,
		)

		l = append(l, config.Repository{
			Name:  mr.Title,
			Desc:  desc,
			Route: mr.URL,
		})
	}

	return l
}

func FetchMRsFromRepo(c *config.Config, project string) (*MRsResponse, error) {
	query := fmt.Sprintf(fetchMergeRequests, project)
	token := fmt.Sprintf("Bearer %s", c.Token)

	body := map[string]string{
		"query": query,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewBuffer(b)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/graphql", c.BasePath), reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var mrs MRsResponse
	err = json.Unmarshal(respBody, &mrs)
	if err != nil {
		return nil, err
	}

	return &mrs, nil
}

func FilterDraft(mrs *MRsResponse) *MRsResponse {
	result := []MR{}
	for _, mr := range mrs.Data.Project.MergeRequests.Nodes {
		if !strings.HasPrefix(strings.ToLower(mr.Title), "draft") {
			result = append(result, mr)
		}
	}

	return &MRsResponse{
		Data{
			Project{
				MergeRequests{
					Nodes: result,
				},
			},
		},
	}
}
