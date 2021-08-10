package tanzuclient

import (
	"fmt"
	"net/http"
)

type Workspace struct {
	// The name of the workspace.
	FullName FullName `json:"fullName"`
	// The metadata of the workspace.
	Meta MetaData `json:"meta"`
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

func (c *Client) GetAllWorkspaces() ([]Workspace, error) {
	tmcURL := fmt.Sprintf("%s/v1alpha1/workspaces", c.baseURL)

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
