package todoist

import (
	"fmt"
	"net/http"
)

var SECTIONS_ENDPOINT = "sections"

type Section struct {
	Id        string `json:"id"`
	ProjectId string `json:"project_id"`
	Order     int    `json:"order"`
	Name      string `json:"name"`
}

func (client Client) GetAllSections(projectId string) (*[]Section, error) {
	query := fmt.Sprintf("project_id=%s", projectId)

	var sections *[]Section
	err := client.sendRequest(http.MethodGet, SECTIONS_ENDPOINT, query, &sections)
	if err != nil {
		return nil, err
	}

	return sections, nil
}