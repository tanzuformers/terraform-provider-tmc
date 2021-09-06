package tanzuclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Workspace struct {
	// The name of the workspace.
	FullName *SimpleFullName `json:"fullName"`
	// The metadata of the workspace.
	Meta *MetaData `json:"meta"`
}

type WorkspaceResponse struct {
	Workspace Workspace `json:"workspace"`
}

type AllWorkspaces struct {
	Workspaces []Workspace `json:"workspaces"`
}

func (c *Client) GetWorkspace(name string) (*Workspace, error) {
	tmcURL := fmt.Sprintf("%s/v1alpha1/workspaces/%s", c.baseURL, name)

	req, err := http.NewRequest("GET", tmcURL, nil)
	if err != nil {
		return nil, err
	}

	res := WorkspaceResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Workspace, nil
}

func (c *Client) GetAllWorkspaces(labels map[string]interface{}) ([]Workspace, error) {
	queryString := buildLabelQuery(labels)

	tmcURL := fmt.Sprintf("%s/v1alpha1/workspaces?query=%s", c.baseURL, queryString)

	req, err := http.NewRequest("GET", tmcURL, nil)
	if err != nil {
		return nil, err
	}

	res := AllWorkspaces{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return res.Workspaces, nil
}

func (c *Client) CreateWorkspace(name string, description string, labels map[string]interface{}) (*Workspace, error) {
	tmcURL := fmt.Sprintf("%s/v1alpha1/workspaces", c.baseURL)

	workspace := &Workspace{
		FullName: &SimpleFullName{
			Name: name,
		},
		Meta: &MetaData{
			Description: description,
			SimpleMetaData: SimpleMetaData{
				Labels: labels,
			},
		},
	}

	workspaceResponse := &WorkspaceResponse{
		Workspace: *workspace,
	}

	// Create JSON object for the request Body
	json_data, err := json.Marshal(workspaceResponse) // returns []byte
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", tmcURL, bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}

	res := WorkspaceResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Workspace, nil
}

func (c *Client) DeleteWorkspace(name string) error {
	tmcURL := fmt.Sprintf("%s/v1alpha1/workspaces/%s", c.baseURL, name)

	req, err := http.NewRequest("DELETE", tmcURL, nil)
	if err != nil {
		return err
	}

	res := WorkspaceResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateWorkspace(name string, description string, labels map[string]interface{}) (*Workspace, error) {
	tmcURL := fmt.Sprintf("%s/v1alpha1/workspaces/%s", c.baseURL, name)

	workspace := &Workspace{
		FullName: &SimpleFullName{
			Name: name,
		},
		Meta: &MetaData{
			SimpleMetaData: SimpleMetaData{
				Labels: labels,
			},
			Description: description,
		},
	}

	workspaceResponse := &WorkspaceResponse{
		Workspace: *workspace,
	}

	// Create JSON object for the request Body
	json_data, err := json.Marshal(workspaceResponse) // returns []byte
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", tmcURL, bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}

	res := WorkspaceResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Workspace, nil
}
