package tanzuclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ClusterGroup struct {
	// Name of the Cluster Group
	FullName *FullName `json:"fullName"`
	// Metadata about the Cluster Group
	Meta *MetaData `json:"meta"`
}

type ClusterGroupResponse struct {
	ClusterGroup ClusterGroup `json:"clusterGroup"`
}

// Fetch Details about an existing Cluster Group using its name
func (c *Client) GetClusterGroup(name string) (*ClusterGroup, error) {
	requestURL := fmt.Sprintf("%s/v1alpha1/clustergroups/%s", c.baseURL, name)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	res := ClusterGroupResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.ClusterGroup, nil
}

// Create a new Cluster Group with a given name.
// Also accepts a description for the Cluster Group and
// a set of labels to be added
func (c *Client) CreateClusterGroup(name string, description string, labels map[string]interface{}) (*ClusterGroup, error) {

	requestURL := c.baseURL + "/v1alpha1/clustergroups"

	newClusterGroup := &ClusterGroup{
		FullName: &FullName{
			Name: name,
		},
		Meta: &MetaData{
			Description: description,
			Labels:      labels,
		},
	}

	newCgObject := &ClusterGroupResponse{
		ClusterGroup: *newClusterGroup,
	}

	// Create JSON object for the request Body
	json_data, err := json.Marshal(newCgObject) // returns []byte
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}

	res := ClusterGroupResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.ClusterGroup, nil
}

func (c *Client) DeleteClusterGroup(name string) error {
	requestURL := c.baseURL + "/v1alpha1/clustergroups/" + name

	req, err := http.NewRequest("DELETE", requestURL, nil)
	if err != nil {
		return err
	}

	res := ClusterGroupResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateClusterGroup(name string, description string, labels map[string]interface{}) (*ClusterGroup, error) {

	requestURL := c.baseURL + "/v1alpha1/clustergroups/" + name

	newClusterGroup := &ClusterGroup{
		FullName: &FullName{
			Name: name,
		},
		Meta: &MetaData{
			Description: description,
			Labels:      labels,
		},
	}

	newCgObject := &ClusterGroupResponse{
		ClusterGroup: *newClusterGroup,
	}

	// Create JSON object for the request Body
	json_data, err := json.Marshal(newCgObject) // returns []byte
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", requestURL, bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}

	res := ClusterGroupResponse{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.ClusterGroup, nil
}
