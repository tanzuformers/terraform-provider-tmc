package tanzuclient

import (
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
