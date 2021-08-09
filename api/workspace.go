package api

import "fmt"

type Workspace struct {
	// The name of the workspace.
	Name string `json:"name"`
	// The description of the workspace.
	Description string `json:"description"`
}

func (c *Client) GetWorkspace(name string) (*Workspace, error) {
	tmcURL := fmt.Sprintf("%workspaces/%s", c.baseURL, name)

	var a Workspace

	err := c.get(tmcURL, &a)
	if err != nil {
		return nil, err
	}

	return &a, nil
}
