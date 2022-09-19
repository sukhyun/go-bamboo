package bamboo

import (
	"fmt"
	"net/http"
	"strconv"
)

// PlanService handles communication with the plan related methods
type PlanService service

// PlanCreateBranchOptions specifies the optional parameters
// for the CreatePlanBranch method
type PlanCreateBranchOptions struct {
	VCSBranch string
}

// PlanResponse encapsultes a response from the plan service
type PlanResponse struct {
	*ResourceMetadata
	Plans *Plans `json:"plans"`
}

// Plans is a collection of Plan objects
type Plans struct {
	*CollectionMetadata
	PlanList []*Plan `json:"plan"`
}

// Plan is the definition of a single plan
type Plan struct {
	ShortName string   `json:"shortName,omitempty"`
	ShortKey  string   `json:"shortKey,omitempty"`
	Type      string   `json:"type,omitempty"`
	Enabled   bool     `json:"enabled,omitempty"`
	Link      *Link    `json:"link,omitempty"`
	Key       string   `json:"key,omitempty"`
	Name      string   `json:"name,omitempty"`
	PlanKey   *PlanKey `json:"planKey,omitempty"`
}

// PlanKey holds the plan-key for a plan
type PlanKey struct {
	Key string `json:"key,omitempty"`
}

// SpecResponse is the information of specification
type SpecResponse struct {
	Spec *SpecDetail `json:"spec"`
}

type SpecDetail struct {
	ProjectKey string `json:"projectKey,omitempty"`
	BuildKey string `json:"buildKey,omitempty"`
	Code string `json:"code,omitempty`
}

// CreatePlanBranch will create a plan branch with the given branch name for the specified build
func (p *PlanService) CreatePlanBranch(planKey, branchName string, options *PlanCreateBranchOptions) (bool, *http.Response, error) {
	var u string
	if !emptyStrings(planKey, branchName) {
		u = fmt.Sprintf("plan/%s/branch/%s.json", planKey, branchName)
	} else {
		return false, nil, &simpleError{"Project key and/or branch name cannot be empty"}
	}

	request, err := p.client.NewRequest(http.MethodPut, u, nil)
	if err != nil {
		return false, nil, err
	}

	if options != nil && options.VCSBranch != "" {
		values := request.URL.Query()
		values.Add("vcsBranch", options.VCSBranch)
		request.URL.RawQuery = values.Encode()
	}

	response, err := p.client.Do(request, nil)
	if err != nil {
		return false, response, err
	}

	if !(response.StatusCode == 200) {
		return false, response, &simpleError{fmt.Sprintf("Create returned %d", response.StatusCode)}
	}

	return true, response, nil
}

// NumberOfPlans returns the number of plans on the Bamboo server
func (p *PlanService) NumberOfPlans() (int, *http.Response, error) {
	request, err := p.client.NewRequest(http.MethodGet, "plan.json", nil)
	if err != nil {
		return 0, nil, err
	}

	// Restrict results to one for speed
	values := request.URL.Query()
	values.Add("max-results", "1")
	request.URL.RawQuery = values.Encode()

	planResp := PlanResponse{}
	response, err := p.client.Do(request, &planResp)
	if err != nil {
		return 0, response, err
	}

	if response.StatusCode != 200 {
		return 0, response, &simpleError{fmt.Sprintf("Getting the number of plans returned %s", response.Status)}
	}

	return planResp.Plans.Size, response, nil
}

// ListPlans gets information on all plans
func (p *PlanService) ListPlans() ([]*Plan, *http.Response, error) {
	// Get number of plans to set max-results
	numPlans, resp, err := p.NumberOfPlans()
	if err != nil {
		return nil, resp, err
	}

	request, err := p.client.NewRequest(http.MethodGet, "plan.json", nil)
	if err != nil {
		return nil, nil, err
	}

	q := request.URL.Query()
	q.Add("max-results", strconv.Itoa(numPlans))
	request.URL.RawQuery = q.Encode()

	planResp := PlanResponse{}
	response, err := p.client.Do(request, &planResp)
	if err != nil {
		return nil, response, err
	}

	if response.StatusCode != 200 {
		return nil, response, &simpleError{fmt.Sprintf("Getting plan information returned %s", response.Status)}
	}

	return planResp.Plans.PlanList, response, nil
}

// ListPlanKeys get all the plan keys for all build plans on Bamboo
func (p *PlanService) ListPlanKeys() ([]string, *http.Response, error) {
	plans, response, err := p.ListPlans()
	if err != nil {
		return nil, response, err
	}
	keys := make([]string, len(plans))

	for i, p := range plans {
		keys[i] = p.Key
	}
	return keys, response, nil
}

// ListPlanNames returns a list of ShortNames of all plans
func (p *PlanService) ListPlanNames() ([]string, *http.Response, error) {
	plans, response, err := p.ListPlans()
	if err != nil {
		return nil, response, err
	}
	names := make([]string, len(plans))

	for i, p := range plans {
		names[i] = p.ShortName
	}
	return names, response, nil
}

// PlanNameMap returns a map[string]string where the PlanKey is the key and the ShortName is the value
func (p *PlanService) PlanNameMap() (map[string]string, *http.Response, error) {
	plans, response, err := p.ListPlans()
	if err != nil {
		return nil, response, err
	}

	planMap := make(map[string]string, len(plans))

	for _, p := range plans {
		planMap[p.Key] = p.ShortName
	}
	return planMap, response, nil
}

// DisablePlan will disable a plan or plan branch
func (p *PlanService) DisablePlan(planKey string) (*http.Response, error) {
	u := fmt.Sprintf("plan/%s/enable", planKey)
	request, err := p.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}

	response, err := p.client.Do(request, nil)
	if err != nil {
		return response, err
	}
	return response, nil
}

// GetSpecs gets informations on plan's spec
func (p *PlanService) GetSpecs(key string) (string, *http.Response, error) {
	requestValue := fmt.Sprintf("plan/%s/specs?format=YAML", key)
	request, err := p.client.NewRequest(http.MethodGet, requestValue, nil)
	if err != nil {
		return "", nil, err
	}
  
	// Restrict results to one for speed
	values := request.URL.Query()
	values.Add("max-results", "1")
	request.URL.RawQuery = values.Encode()
  
	specResp := SpecResponse{}
	response, err := p.client.Do(request, &specResp)
	if err != nil {
		return "", response, err
	}
  
	if response.StatusCode != 200 {
		return "", response, &simpleError{fmt.Sprintf("Getting the spec of plans returned %s", response.Status)}
	}
  
	return specResp.Spec.Code, response, nil
}
