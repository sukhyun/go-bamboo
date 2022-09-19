package bamboo

import (
	"errors"
	"fmt"
	"net/http"
)

// ProjectService handles communication with the project related methods
type ProjectService service

// ProjectResponse the REST response from the server
type ProjectResponse struct {
	*ResourceMetadata
	Projects *Projects `json:"projects"`
}

// Projects is a collection of project elements
type Projects struct {
	*CollectionMetadata
	ProjectList []*Project `json:"project"`
}

// Project is a single project definition
type Project struct {
	Key         string `json:"key,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Link        *Link  `json:"link,omitempty"`
}

// ProjectInformation is the information for a single project
type ProjectInformation struct {
	Key         string                   `json:"key,omitempty"`
	Name        string                   `json:"name,omitempty"`
	Description string                   `json:"description,omitempty"`
	NumPlans    *ProjectPlansInformation `json:"plans"`
}

// ProjectPlansInformation holds the number of plans in a project
type ProjectPlansInformation struct {
	Size int `json:"size,omitempty"`
}

// ProjectRepos is the information for project's repositories
type ProjectRepos struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Url string `json:"url"`
	Location string `json:"location"`
	Icon string `json:"string"`
	Type string `json:"type"`
	IsAdmin bool `json:"isAdmin"`
}

type ProjectRepositoryResult struct {
	*Index
	Repositories []*ProjectRepos `json:"results"`
}

// ProjectInfo get the information on the specific project
func (p *ProjectService) ProjectInfo(projectKey string) (*ProjectInformation, *http.Response, error) {
	var u string
	if !emptyStrings(projectKey) {
		u = fmt.Sprintf("project/%s.json", projectKey)
	} else {
		return nil, nil, &simpleError{fmt.Sprintf("Project key cannot be an empty string")}
	}

	request, err := p.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	projectInfo := ProjectInformation{}
	response, err := p.client.Do(request, &projectInfo)
	if err != nil {
		return nil, nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, response,
			errors.New("Getting Project Information returned: " + response.Status)
	}

	return &projectInfo, response, nil
}

// ProjectPlans returns a list of plans for a given project
func (p *ProjectService) ProjectPlans(projectKey string) ([]*Plan, *http.Response, error) {
	var u string
	if !emptyStrings(projectKey) {
		u = fmt.Sprintf("project/%s.json", projectKey)
	} else {
		return nil, nil, &simpleError{fmt.Sprintf("Project key cannot be an empty string")}
	}

	request, err := p.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	values := request.URL.Query()
	values.Set("expand", "plans")
	values.Set("max-result", "1000")
	request.URL.RawQuery = values.Encode()

	projectResponse := PlanResponse{}
	response, err := p.client.Do(request, &projectResponse)
	if err != nil {
		return nil, nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, response,
			errors.New("Getting Project Plans returned: " + response.Status)
	}

	return projectResponse.Plans.PlanList, response, nil
}

// ListProjects lists all projects
func (p *ProjectService) ListProjects() ([]*Project, *http.Response, error) {
	u := "project.json"

	request, err := p.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	projectResp := ProjectResponse{}
	response, err := p.client.Do(request, &projectResp)
	if err != nil {
		return nil, nil, err
	}

	if !(response.StatusCode == 200) {
		return nil, response, &simpleError{fmt.Sprintf("List projects returned %s", response.Status)}
	}

	return projectResp.Projects.ProjectList, response, nil
}

func (p *ProjectService) ProjectRepositories(projectKey string) ([]*ProjectRepos, *http.Response, error) {
	u := fmt.Sprintf("project/%s/repositories", projectKey)
  
	request, err := p.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
  
	projectInfo := ProjectRepositoryResult{}
	response, err := p.client.Do(request, &projectInfo)
	if err != nil {
		return nil, nil, err
	}
  
	if response.StatusCode != http.StatusOK {
		return nil, response,
		&simpleError{fmt.Sprintf("Getting Project Information returned: %s", response.Status)}
	}
  
	return projectInfo.Repositories, response, nil
 }
